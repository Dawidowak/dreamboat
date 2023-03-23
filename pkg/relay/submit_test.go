package relay

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/blocknative/dreamboat/blstools"
	rpctypes "github.com/blocknative/dreamboat/pkg/client/sim/types"
	"github.com/blocknative/dreamboat/pkg/relay/mocks"
	"github.com/blocknative/dreamboat/pkg/structs"
	"github.com/blocknative/dreamboat/pkg/structs/forks/bellatrix"
	"github.com/blocknative/dreamboat/pkg/structs/forks/capella"
	ccommon "github.com/blocknative/dreamboat/test/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/flashbots/go-boost-utils/bls"
	"github.com/flashbots/go-boost-utils/types"
	"github.com/golang/mock/gomock"
	"github.com/lthibault/log"
	"github.com/stretchr/testify/require"
)

type fields struct {
	d           Datastore
	a           Auctioneer
	ver         Verifier
	config      RelayConfig
	cache       ValidatorCache
	vstore      ValidatorStore
	bvc         BlockValidationClient
	beacon      Beacon
	beaconState State
}

func simpletest(t require.TestingT, ctrl *gomock.Controller, fork structs.ForkVersion, submitRequest structs.SubmitBlockRequest, sk *bls.SecretKey, pubKey types.PublicKey, relaySigningDomain types.Domain, genesisTime uint64) fields {

	conf := RelayConfig{
		BuilderSigningDomain: relaySigningDomain,
		SecretKey:            sk,
		PubKey:               pubKey,
	}

	ds := mocks.NewMockDatastore(ctrl)
	state := mocks.NewMockState(ctrl)
	cache := mocks.NewMockValidatorCache(ctrl)
	vstore := mocks.NewMockValidatorStore(ctrl)
	verify := mocks.NewMockVerifier(ctrl)
	bvc := mocks.NewMockBlockValidationClient(ctrl)
	a := mocks.NewMockAuctioneer(ctrl)

	// Submit Block

	state.EXPECT().Genesis().MaxTimes(1).Return(
		structs.GenesisInfo{GenesisTime: genesisTime},
	)
	state.EXPECT().HeadSlot().MaxTimes(1).Return(
		structs.Slot(submitRequest.Slot()),
	)

	state.EXPECT().Randao().MaxTimes(1).Return(submitRequest.Random().String())

	ds.EXPECT().CheckSlotDelivered(context.Background(), submitRequest.Slot()).MaxTimes(1).Return(
		false, nil,
	)
	cache.EXPECT().Get(submitRequest.ProposerPubkey()).Return(
		structs.ValidatorCacheEntry{}, false,
	)

	state.EXPECT().ForkVersion(structs.Slot(submitRequest.Slot())).Times(2).Return(fork)

	vstore.EXPECT().GetRegistration(context.Background(), submitRequest.ProposerPubkey()).MaxTimes(1).Return(
		types.SignedValidatorRegistration{
			Message: &types.RegisterValidatorRequestMessage{
				FeeRecipient: submitRequest.ProposerFeeRecipient(),
			},
		}, nil,
	)

	cache.EXPECT().Add(submitRequest.ProposerPubkey(), gomock.Any()).Return(false) // todo check ValidatorCacheEntry disregarding time.Now()
	msg, err := submitRequest.ComputeSigningRoot(relaySigningDomain)
	require.NoError(t, err)
	verify.EXPECT().Enqueue(context.Background(), submitRequest.Signature(), submitRequest.BuilderPubkey(), msg).Times(1)

	bvc.EXPECT().IsSet().Times(1).Return(true)
	switch fork {
	case structs.ForkBellatrix:
		bvc.EXPECT().ValidateBlock(context.Background(), &rpctypes.BuilderBlockValidationRequest{
			SubmitBlockRequest: submitRequest,
		}).Return(nil)
	case structs.ForkCapella:
		bvc.EXPECT().ValidateBlockV2(context.Background(), &rpctypes.BuilderBlockValidationRequest{
			SubmitBlockRequest: submitRequest,
		}).Return(nil)
	}

	contents, err := submitRequest.PreparePayloadContents(sk, &pubKey, relaySigningDomain)
	require.NoError(t, err)
	log.Debug(contents)

	ds.EXPECT().PutPayload(context.Background(), submitRequest.ToPayloadKey(), contents.Payload, conf.TTL).Return(nil)

	var bl *structs.CompleteBlockstruct
	a.EXPECT().AddBlock(gomock.Any()).Times(1).DoAndReturn(func(block *structs.CompleteBlockstruct) bool {
		bl = block
		return true
	})

	ds.EXPECT().PutHeader(gomock.Any(), gomock.Any(), conf.TTL)

	////   GetHeader

	a.EXPECT().MaxProfitBlock(structs.Slot(submitRequest.Slot())).Times(1).
		DoAndReturn(func(slot structs.Slot) (block *structs.CompleteBlockstruct, a bool) {
			return bl, true
		})

	ds.EXPECT().CacheBlock(gomock.Any(), structs.PayloadKey{
		BlockHash: submitRequest.BlockHash(),
		Slot:      structs.Slot(submitRequest.Slot()),
		Proposer:  submitRequest.ProposerPubkey()},
		gomock.Any()).Times(1).DoAndReturn(func(c context.Context, s structs.PayloadKey, block *structs.CompleteBlockstruct) error {
		if bl != block {
			t.Errorf("cache block is different structure")
		}
		return nil
	})
	switch fork {
	case structs.ForkCapella:
		state.EXPECT().Withdrawals().Times(1).Return(structs.WithdrawalsState{
			Slot: structs.Slot(submitRequest.Slot()),
			Root: types.Root{},
		})
	}

	//state.EXPECT().ForkVersion(structs.Slot(submitRequest.Slot())).Times(1).Return(fork)
	/*structs.HeaderData{
		Slot:           structs.Slot(submitRequest.Slot()),
		Marshaled:      m,
		HeaderAndTrace: contents.Header,
	}*/
	//	sbr, ok := submitRequest.(*bellatrix.SubmitBlockRequest)
	//if ok {

	//}
	return fields{
		config:      conf,
		d:           ds,
		a:           a,
		ver:         verify,
		cache:       cache,
		vstore:      vstore,
		bvc:         bvc,
		beacon:      mocks.NewMockBeacon(ctrl),
		beaconState: state,
	}
}

func TestRelay_SubmitBlock(t *testing.T) {

	type args struct {
		ctx context.Context
		m   *structs.MetricGroup
		sbr structs.SubmitBlockRequest
	}

	relaySigningDomain, err := ccommon.ComputeDomain(
		types.DomainTypeAppBuilder,
		types.Root{}.String())
	require.NoError(t, err)
	genesisTime := uint64(time.Now().Unix())
	sk, pubKey, err := blstools.GenerateNewKeypair()
	require.NoError(t, err)

	tests := []struct {
		name string
		fork structs.ForkVersion

		args    args
		wantErr bool
	}{
		{
			name: "bellatrix simple",
			fork: structs.ForkBellatrix,
			args: args{
				ctx: context.Background(),
				m:   &structs.MetricGroup{},
				sbr: func() structs.SubmitBlockRequest {
					return validSubmitBlockRequestBellatrix(t, sk, pubKey, relaySigningDomain, genesisTime)
				}(),
			},
		},
		{
			name: "capella simple",
			fork: structs.ForkCapella,
			args: args{
				ctx: context.Background(),
				m:   &structs.MetricGroup{},
				sbr: func() structs.SubmitBlockRequest {
					return validSubmitBlockRequestCapella(t, sk, pubKey, relaySigningDomain, genesisTime)
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := log.New()
			controller := gomock.NewController(t)

			f := simpletest(t, controller, tt.fork, tt.args.sbr, sk, pubKey, relaySigningDomain, genesisTime)
			rs := NewRelay(l,
				f.config,
				f.beacon,
				f.cache,
				f.vstore,
				f.ver,
				f.beaconState,
				f.d,
				f.a,
				f.bvc)

			if err := rs.SubmitBlock(tt.args.ctx, tt.args.m, tt.args.sbr); (err != nil) != tt.wantErr {
				t.Errorf("Relay.SubmitBlock() error = %v, wantErr %v", err, tt.wantErr)
			}

			var pHash types.Hash
			switch s := tt.args.sbr.(type) {
			case *bellatrix.SubmitBlockRequest:
				pHash = s.BellatrixExecutionPayload.EpParentHash
			case *capella.SubmitBlockRequest:
				pHash = s.CapellaExecutionPayload.EpParentHash
			}

			req := structs.HeaderRequest{
				"slot":        strconv.FormatUint(tt.args.sbr.Slot(), 10),
				"parent_hash": pHash.String(),
				"pubkey":      tt.args.sbr.ProposerPubkey().String(),
			}
			if _, err := rs.GetHeader(tt.args.ctx, tt.args.m, req); (err != nil) != tt.wantErr {
				t.Errorf("Relay.GetHeader() error = %v, wantErr %v", err, tt.wantErr)
			}

			//rs.GetPayload(tt.args.ctx, tt.args.m,)
		})
	}
}

/*
func (hr HeaderRequest) Slot() (Slot, error) {
	slot, err := strconv.Atoi(hr["slot"])
	return Slot(slot), err
}

func (hr HeaderRequest) ParentHash() (types.Hash, error) {
	var parentHash types.Hash
	err := parentHash.UnmarshalText([]byte(strings.ToLower(hr["parent_hash"])))
	return parentHash, err
}

func (hr HeaderRequest) Pubkey() (PubKey, error) {
	var pk PubKey
	if err := pk.UnmarshalText([]byte(strings.ToLower(hr["pubkey"]))
}*/

func validSubmitBlockRequestBellatrix(t require.TestingT, sk *bls.SecretKey, pubKey types.PublicKey, domain types.Domain, genesisTime uint64) *bellatrix.SubmitBlockRequest {

	slot := rand.Uint64()

	payload := randomPayload()
	payload.EpTimestamp = genesisTime + (slot * 12)

	msg := types.BidTrace{
		Slot:                 slot,
		ParentHash:           payload.EpParentHash,
		BlockHash:            payload.EpBlockHash,
		BuilderPubkey:        pubKey,
		ProposerPubkey:       types.PublicKey(random48Bytes()),
		ProposerFeeRecipient: types.Address(random20Bytes()),
		Value:                types.IntToU256(rand.Uint64()),
	}

	signature, err := types.SignMessage(&msg, domain, sk)
	require.NoError(t, err)

	return &bellatrix.SubmitBlockRequest{
		BellatrixSignature:        signature,
		BellatrixMessage:          msg,
		BellatrixExecutionPayload: *payload,
	}
}

func validSubmitBlockRequestCapella(t require.TestingT, sk *bls.SecretKey, pubKey types.PublicKey, domain types.Domain, genesisTime uint64) *capella.SubmitBlockRequest {

	slot := rand.Uint64()
	random := randomPayload()

	payload := capella.ExecutionPayload{
		ExecutionPayload: *random,
		EpWithdrawals: []*structs.Withdrawal{{
			Index:          1,
			ValidatorIndex: 1,
			Address:        types.Address(random20Bytes()),
			Amount:         1234,
		}},
	}

	payload.EpTimestamp = genesisTime + (slot * 12)

	msg := types.BidTrace{
		Slot:                 slot,
		ParentHash:           payload.EpParentHash,
		BlockHash:            payload.EpBlockHash,
		BuilderPubkey:        pubKey,
		ProposerPubkey:       types.PublicKey(random48Bytes()),
		ProposerFeeRecipient: types.Address(random20Bytes()),
		Value:                types.IntToU256(rand.Uint64()),
	}

	signature, err := types.SignMessage(&msg, domain, sk)
	require.NoError(t, err)

	return &capella.SubmitBlockRequest{
		CapellaSignature:        signature,
		CapellaMessage:          msg,
		CapellaExecutionPayload: payload,
	}
}

func random48Bytes() (b [48]byte) {
	rand.Read(b[:])
	return b
}

func random20Bytes() (b [20]byte) {
	rand.Read(b[:])
	return b
}

func random32Bytes() (b [32]byte) {
	rand.Read(b[:])
	return b
}

func random256Bytes() (b [256]byte) {
	rand.Read(b[:])
	return b
}

func randomPayload() *bellatrix.ExecutionPayload {
	return &bellatrix.ExecutionPayload{
		EpParentHash:    types.Hash(random32Bytes()),
		EpFeeRecipient:  types.Address(random20Bytes()),
		EpStateRoot:     types.Hash(random32Bytes()),
		EpReceiptsRoot:  types.Hash(random32Bytes()),
		EpLogsBloom:     types.Bloom(random256Bytes()),
		EpRandom:        random32Bytes(),
		EpBlockNumber:   rand.Uint64(),
		EpGasLimit:      rand.Uint64(),
		EpGasUsed:       rand.Uint64(),
		EpTimestamp:     rand.Uint64(),
		EpExtraData:     types.ExtraData{},
		EpBaseFeePerGas: types.IntToU256(rand.Uint64()),
		EpBlockHash:     types.Hash(random32Bytes()),
		EpTransactions:  randomTransactions(2),
	}
}

func randomTransactions(size int) []hexutil.Bytes {
	txs := make([]hexutil.Bytes, 0, size)
	for i := 0; i < size; i++ {
		tx := make([]byte, rand.Intn(32))
		rand.Read(tx)
		txs = append(txs, tx)
	}
	return txs
}
