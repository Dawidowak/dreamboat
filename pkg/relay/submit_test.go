package relay

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/blocknative/dreamboat/blstools"
	"github.com/blocknative/dreamboat/pkg/relay/mocks"
	"github.com/blocknative/dreamboat/pkg/structs"
	"github.com/blocknative/dreamboat/pkg/structs/forks/bellatrix"
	"github.com/blocknative/dreamboat/test/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/flashbots/go-boost-utils/types"
	"github.com/golang/mock/gomock"
	"github.com/lthibault/log"
	"github.com/stretchr/testify/require"
)

func TestRelay_SubmitBlock(t *testing.T) {
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
	type args struct {
		ctx context.Context
		m   *structs.MetricGroup
		sbr structs.SubmitBlockRequest
	}

	relaySigningDomain, err := common.ComputeDomain(
		types.DomainTypeAppBuilder,
		types.Root{}.String())
	require.NoError(t, err)
	genesisTime := uint64(time.Now().Unix())

	tests := []struct {
		name          string
		GenerateMocks func(*gomock.Controller, structs.SubmitBlockRequest) fields
		args          args
		wantErr       bool
	}{
		{
			name: "bellatrix simple",
			GenerateMocks: func(ctrl *gomock.Controller, submitRequest structs.SubmitBlockRequest) fields {
				ds := mocks.NewMockDatastore(ctrl)
				state := mocks.NewMockState(ctrl)
				cache := mocks.NewMockValidatorCache(ctrl)
				vstore := mocks.NewMockValidatorStore(ctrl)

				state.EXPECT().Genesis().MaxTimes(1).Return(
					structs.GenesisInfo{GenesisTime: genesisTime},
				)
				state.EXPECT().HeadSlot().MaxTimes(1).Return(
					structs.Slot(submitRequest.Slot()),
				)
				ds.EXPECT().CheckSlotDelivered(context.Background(), submitRequest.Slot()).MaxTimes(1).Return(
					false, nil,
				)
				cache.EXPECT().Get(submitRequest.ProposerPubkey()).Return(
					structs.ValidatorCacheEntry{}, false,
				)

				vstore.EXPECT().GetRegistration(context.Background(), submitRequest.ProposerPubkey()).MaxTimes(1).Return(
					types.SignedValidatorRegistration{}, nil,
				)

				return fields{
					config:      RelayConfig{},
					d:           ds,
					a:           mocks.NewMockAuctioneer(ctrl),
					ver:         mocks.NewMockVerifier(ctrl),
					cache:       cache,
					vstore:      vstore,
					bvc:         mocks.NewMockBlockValidationClient(ctrl),
					beacon:      mocks.NewMockBeacon(ctrl),
					beaconState: state,
				}
			},
			args: args{
				ctx: context.Background(),
				m:   &structs.MetricGroup{},
				sbr: func() structs.SubmitBlockRequest {
					return validSubmitBlockRequest(t, relaySigningDomain, genesisTime)
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := log.New()
			controller := gomock.NewController(t)

			f := tt.GenerateMocks(controller, tt.args.sbr)
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
		})
	}
}

func validSubmitBlockRequest(t require.TestingT, domain types.Domain, genesisTime uint64) *bellatrix.SubmitBlockRequest {
	sk, pubKey, err := blstools.GenerateNewKeypair()
	require.NoError(t, err)

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
