package validators_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/blocknative/dreamboat/blstools"
	pkg "github.com/blocknative/dreamboat/pkg"
	"github.com/blocknative/dreamboat/pkg/datastore"
	relay "github.com/blocknative/dreamboat/pkg/relay"
	mock_relay "github.com/blocknative/dreamboat/pkg/relay/mocks"
	"github.com/blocknative/dreamboat/pkg/structs"
	"github.com/blocknative/dreamboat/pkg/validators"
	"github.com/blocknative/dreamboat/pkg/verify"
	"github.com/flashbots/go-boost-utils/bls"
	"github.com/flashbots/go-boost-utils/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	badger "github.com/ipfs/go-ds-badger2"
	"github.com/lthibault/log"
	"github.com/stretchr/testify/require"
)

func TestGetValidators(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	bs := mock_relay.NewMockState(ctrl)

	l := log.New()
	ver := verify.NewVerificationManager(l, 20)
	ver.RunVerify(300)

	relaySigningDomain, err := pkg.ComputeDomain(
		types.DomainTypeAppBuilder,
		pkg.GenesisForkVersionRopsten,
		types.Root{}.String())
	require.NoError(t, err)

	r := validators.NewRegister(l, relaySigningDomain, bs, ver, nil, nil)

	fbn := &structs.BeaconState{
		DutiesState: structs.DutiesState{
			ProposerDutiesResponse: structs.BuilderGetValidatorsResponseEntrySlice{{
				Slot:  0,
				Entry: &types.SignedValidatorRegistration{},
			}},
		},
	}
	bs.EXPECT().Beacon().Return(fbn).Times(1)

	validators := r.GetValidators(structs.NewMetricGroup(4))
	require.NotNil(t, validators)
}

func TestRegisterValidator(t *testing.T) {
	t.Parallel()

	const N = 100

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)

	var datadir = "/tmp/" + t.Name() + uuid.New().String()
	store, _ := badger.NewDatastore(datadir, &badger.DefaultOptions)
	ds := &datastore.Datastore{TTLStorage: &datastore.TTLDatastoreBatcher{TTLDatastore: store}}
	bs := mock_relay.NewMockState(ctrl)

	relaySigningDomain, err := pkg.ComputeDomain(
		types.DomainTypeAppBuilder,
		pkg.GenesisForkVersionRopsten,
		types.Root{}.String())
	require.NoError(t, err)

	l := log.New()

	storeMgr, err := validators.NewStoreManager(l, 20, 20000, 20)
	if err != nil {
		t.Error()
	}
	storeMgr.RunStore(ds, time.Minute, 300)

	ver := verify.NewVerificationManager(l, 20000)
	ver.RunVerify(300)

	fbn := &structs.BeaconState{
		ValidatorsState: structs.ValidatorsState{
			KnownValidators: make(map[types.PubkeyHex]struct{}),
		},
	}
	registrations := make([]types.SignedValidatorRegistration, 0, N)
	for i := 0; i < N; i++ {
		registration, _ := validValidatorRegistration(t, relaySigningDomain)
		registrations = append(registrations, *registration)

		fbn.ValidatorsState.KnownValidators[registration.Message.Pubkey.PubkeyHex()] = struct{}{}
	}
	bs.EXPECT().Beacon().Return(fbn)

	vr := validators.NewRegister(l, relaySigningDomain, bs, ver, storeMgr, ds)
	err = vr.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	for _, registration := range registrations {
		key := structs.PubKey{PublicKey: registration.Message.Pubkey}
		gotRegistration, err := ds.GetRegistration(ctx, key)
		require.NoError(t, err)
		require.EqualValues(t, registration, gotRegistration)
	}
}

func TestBrokenSignatureRegisterValidator(t *testing.T) {
	t.Parallel()

	const N = 10000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)

	var datadir = "/tmp/" + t.Name() + uuid.New().String()
	store, _ := badger.NewDatastore(datadir, &badger.DefaultOptions)
	ds := &datastore.Datastore{TTLStorage: &datastore.TTLDatastoreBatcher{TTLDatastore: store}}
	bs := mock_relay.NewMockState(ctrl)

	relaySigningDomain, err := pkg.ComputeDomain(
		types.DomainTypeAppBuilder,
		pkg.GenesisForkVersionRopsten,
		types.Root{}.String())
	require.NoError(t, err)

	config := relay.RelayConfig{
		TTL: time.Minute,
	}

	l := log.New()

	storeMgr, err := validators.NewStoreManager(l, 20, 20000, 20)
	if err != nil {
		t.Error()
	}
	storeMgr.RunStore(ds, config.TTL, 300)

	ver := verify.NewVerificationManager(l, 20000)
	ver.RunVerify(300)

	r := validators.NewRegister(l, relaySigningDomain, bs, ver, storeMgr, ds)
	fbn := &structs.BeaconState{
		ValidatorsState: structs.ValidatorsState{
			KnownValidators: make(map[types.PubkeyHex]struct{}),
		},
	}

	registrations := make([]types.SignedValidatorRegistration, 0, N)
	for i := 0; i < N; i++ {
		registration, _ := validValidatorRegistration(t, relaySigningDomain)
		registrations = append(registrations, *registration)

		fbn.ValidatorsState.KnownValidators[registration.Message.Pubkey.PubkeyHex()] = struct{}{}
	}

	registrations[N/2].Signature = types.Signature{}
	bs.EXPECT().Beacon().Return(fbn)

	err = r.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
	require.Error(t, err)
	//t.Logf("returned %s", err.Error())
	time.Sleep(3 * time.Second)

	var errored bool
	for i, registration := range registrations {
		key := structs.PubKey{PublicKey: registration.Message.Pubkey}
		gotRegistration, err := ds.GetRegistration(ctx, key)
		if !errored {
			if i != N/2 {
				if err == nil || err.Error() != "datastore: key not found" {
					require.NoError(t, err)
					require.EqualValues(t, registration, gotRegistration)
				}
			} else {
				errored = true
				require.Error(t, err)
			}
		}
	}
}

func TestNotKnownRegisterValidator(t *testing.T) {
	t.Parallel()

	const N = 10000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)

	var datadir = "/tmp/" + t.Name() + uuid.New().String()
	store, _ := badger.NewDatastore(datadir, &badger.DefaultOptions)
	ds := &datastore.Datastore{TTLStorage: &datastore.TTLDatastoreBatcher{TTLDatastore: store}}
	bs := mock_relay.NewMockState(ctrl)

	relaySigningDomain, err := pkg.ComputeDomain(
		types.DomainTypeAppBuilder,
		pkg.GenesisForkVersionRopsten,
		types.Root{}.String())
	require.NoError(t, err)

	config := relay.RelayConfig{
		TTL: time.Minute,
	}

	l := log.New()

	storeMgr, err := validators.NewStoreManager(l, 20, 20000, 20)
	if err != nil {
		t.Error()
	}
	storeMgr.RunStore(ds, config.TTL, 300)

	ver := verify.NewVerificationManager(l, 20000)
	ver.RunVerify(300)

	r := validators.NewRegister(l, relaySigningDomain, bs, ver, storeMgr, ds)
	fbn := &structs.BeaconState{
		ValidatorsState: structs.ValidatorsState{
			KnownValidators: make(map[types.PubkeyHex]struct{}),
		},
	}

	registrations := make([]types.SignedValidatorRegistration, 0, N)
	for i := 0; i < N; i++ {
		registration, _ := validValidatorRegistration(t, relaySigningDomain)
		registrations = append(registrations, *registration)
		if i != N/2 {
			fbn.ValidatorsState.KnownValidators[registration.Message.Pubkey.PubkeyHex()] = struct{}{}
		}

	}

	bs.EXPECT().Beacon().Return(fbn)
	err = r.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
	require.Error(t, err)
	//t.Logf("returned %s", err.Error())
}

func BenchmarkRegisterValidator(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const N = 10_000

	ctrl := gomock.NewController(b)

	var datadir = "/tmp/" + b.Name() + uuid.New().String()
	store, _ := badger.NewDatastore(datadir, &badger.DefaultOptions)
	ds := &datastore.Datastore{TTLStorage: &datastore.TTLDatastoreBatcher{TTLDatastore: store}}
	bs := mock_relay.NewMockState(ctrl)

	relaySigningDomain, _ := pkg.ComputeDomain(
		types.DomainTypeAppBuilder,
		pkg.GenesisForkVersionRopsten,
		types.Root{}.String())

	config := relay.RelayConfig{
		TTL: 5 * time.Minute,
	}

	l := log.New()

	storeMgr, err := validators.NewStoreManager(l, 20, 20000, 20)
	if err != nil {
		b.Error()
	}
	storeMgr.RunStore(ds, config.TTL, 300)

	ver := verify.NewVerificationManager(l, 20000)
	ver.RunVerify(300)

	r := validators.NewRegister(l, relaySigningDomain, bs, ver, storeMgr, ds)

	fbn := &structs.BeaconState{
		ValidatorsState: structs.ValidatorsState{
			KnownValidators: make(map[types.PubkeyHex]struct{}),
		},
	}

	registrations := make([]types.SignedValidatorRegistration, 0, N)
	for i := 0; i < N; i++ {
		registration, _ := validValidatorRegistration(b, relaySigningDomain)
		registrations = append(registrations, *registration)
		fbn.ValidatorsState.KnownValidators[registration.Message.Pubkey.PubkeyHex()] = struct{}{}
	}
	bs.EXPECT().Beacon().Return(fbn).AnyTimes()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := r.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkRegisterValidatorParallel(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var datadir = "/tmp/" + b.Name() + uuid.New().String()
	store, _ := badger.NewDatastore(datadir, &badger.DefaultOptions)
	ds := &datastore.Datastore{TTLStorage: &datastore.TTLDatastoreBatcher{TTLDatastore: store}}

	relaySigningDomain, _ := pkg.ComputeDomain(
		types.DomainTypeAppBuilder,
		pkg.GenesisForkVersionRopsten,
		types.Root{}.String())

	l := log.New()

	storeMgr, err := validators.NewStoreManager(l, 20, 20000, 20)
	if err != nil {
		b.Error()
	}
	storeMgr.RunStore(ds, 5*time.Minute, 300)

	ver := verify.NewVerificationManager(l, 20000)
	ver.RunVerify(300)

	ctrl := gomock.NewController(b)
	bs := mock_relay.NewMockState(ctrl)

	const N = 10_000

	r := validators.NewRegister(l, relaySigningDomain, bs, ver, storeMgr, ds)
	fbn := &structs.BeaconState{
		ValidatorsState: structs.ValidatorsState{
			KnownValidators: make(map[types.PubkeyHex]struct{}),
		},
	}

	registrations := make([]types.SignedValidatorRegistration, 0, N)
	for i := 0; i < N; i++ {
		registration, _ := validValidatorRegistration(b, relaySigningDomain)
		registrations = append(registrations, *registration)
		fbn.ValidatorsState.KnownValidators[registration.Message.Pubkey.PubkeyHex()] = struct{}{}
	}
	bs.EXPECT().Beacon().Return(fbn).AnyTimes()

	var wg sync.WaitGroup
	wg.Add(b.N)

	var wg2 sync.WaitGroup
	defer wg2.Wait()
	wg2.Add(b.N)

	var wg3 sync.WaitGroup
	defer wg3.Wait()
	wg3.Add(b.N)

	b.ResetTimer()
	b.ReportAllocs()
	b.Logf(" b.N %d", b.N)

	for i := 0; i < b.N; i++ {
		go func() {
			t := time.Now()
			err := r.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
			b.Logf(" RegisterValidator %s", time.Since(t).String())
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	for i := 0; i < b.N; i++ {
		go func() {
			t := time.Now()
			err := r.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
			b.Logf(" RegisterValidator %s", time.Since(t).String())
			if err != nil {
				panic(err)
			}
			wg2.Done()
		}()
	}

	wg2.Wait()
	for i := 0; i < b.N; i++ {
		go func() {
			t := time.Now()
			err := r.RegisterValidator(ctx, structs.NewMetricGroup(4), registrations)
			b.Logf(" RegisterValidator %s", time.Since(t).String())
			if err != nil {
				panic(err)
			}
			wg3.Done()
		}()
	}

}

func validValidatorRegistration(t require.TestingT, domain types.Domain) (*types.SignedValidatorRegistration, *bls.SecretKey) {
	sk, pubKey, err := blstools.GenerateNewKeypair()
	require.NoError(t, err)

	msg := &types.RegisterValidatorRequestMessage{
		FeeRecipient: types.Address{0x42},
		GasLimit:     15_000_000,
		Timestamp:    1652369368,
		Pubkey:       pubKey,
	}

	signature, err := types.SignMessage(msg, domain, sk)
	require.NoError(t, err)
	return &types.SignedValidatorRegistration{
		Message:   msg,
		Signature: signature,
	}, sk
}
