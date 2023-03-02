// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blocknative/dreamboat/pkg/beacon (interfaces: Datastore,ValidatorCache,BeaconClient,State)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	client "github.com/blocknative/dreamboat/pkg/beacon/client"
	structs "github.com/blocknative/dreamboat/pkg/structs"
	types "github.com/flashbots/go-boost-utils/types"
	gomock "github.com/golang/mock/gomock"
)

// MockDatastore is a mock of Datastore interface.
type MockDatastore struct {
	ctrl     *gomock.Controller
	recorder *MockDatastoreMockRecorder
}

// MockDatastoreMockRecorder is the mock recorder for MockDatastore.
type MockDatastoreMockRecorder struct {
	mock *MockDatastore
}

// NewMockDatastore creates a new mock instance.
func NewMockDatastore(ctrl *gomock.Controller) *MockDatastore {
	mock := &MockDatastore{ctrl: ctrl}
	mock.recorder = &MockDatastoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatastore) EXPECT() *MockDatastoreMockRecorder {
	return m.recorder
}

// GetRegistration mocks base method.
func (m *MockDatastore) GetRegistration(arg0 context.Context, arg1 types.PublicKey) (types.SignedValidatorRegistration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegistration", arg0, arg1)
	ret0, _ := ret[0].(types.SignedValidatorRegistration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegistration indicates an expected call of GetRegistration.
func (mr *MockDatastoreMockRecorder) GetRegistration(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegistration", reflect.TypeOf((*MockDatastore)(nil).GetRegistration), arg0, arg1)
}

// MockValidatorCache is a mock of ValidatorCache interface.
type MockValidatorCache struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorCacheMockRecorder
}

// MockValidatorCacheMockRecorder is the mock recorder for MockValidatorCache.
type MockValidatorCacheMockRecorder struct {
	mock *MockValidatorCache
}

// NewMockValidatorCache creates a new mock instance.
func NewMockValidatorCache(ctrl *gomock.Controller) *MockValidatorCache {
	mock := &MockValidatorCache{ctrl: ctrl}
	mock.recorder = &MockValidatorCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidatorCache) EXPECT() *MockValidatorCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockValidatorCache) Get(arg0 types.PublicKey) (structs.ValidatorCacheEntry, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(structs.ValidatorCacheEntry)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockValidatorCacheMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockValidatorCache)(nil).Get), arg0)
}

// MockBeaconClient is a mock of BeaconClient interface.
type MockBeaconClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconClientMockRecorder
}

// MockBeaconClientMockRecorder is the mock recorder for MockBeaconClient.
type MockBeaconClientMockRecorder struct {
	mock *MockBeaconClient
}

// NewMockBeaconClient creates a new mock instance.
func NewMockBeaconClient(ctrl *gomock.Controller) *MockBeaconClient {
	mock := &MockBeaconClient{ctrl: ctrl}
	mock.recorder = &MockBeaconClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBeaconClient) EXPECT() *MockBeaconClientMockRecorder {
	return m.recorder
}

// Genesis mocks base method.
func (m *MockBeaconClient) Genesis() (structs.GenesisInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Genesis")
	ret0, _ := ret[0].(structs.GenesisInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Genesis indicates an expected call of Genesis.
func (mr *MockBeaconClientMockRecorder) Genesis() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Genesis", reflect.TypeOf((*MockBeaconClient)(nil).Genesis))
}

// GetProposerDuties mocks base method.
func (m *MockBeaconClient) GetProposerDuties(arg0 structs.Epoch) (*client.RegisteredProposersResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposerDuties", arg0)
	ret0, _ := ret[0].(*client.RegisteredProposersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProposerDuties indicates an expected call of GetProposerDuties.
func (mr *MockBeaconClientMockRecorder) GetProposerDuties(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposerDuties", reflect.TypeOf((*MockBeaconClient)(nil).GetProposerDuties), arg0)
}

// GetWithdrawals mocks base method.
func (m *MockBeaconClient) GetWithdrawals(arg0 structs.Slot) (*client.GetWithdrawalsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawals", arg0)
	ret0, _ := ret[0].(*client.GetWithdrawalsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawals indicates an expected call of GetWithdrawals.
func (mr *MockBeaconClientMockRecorder) GetWithdrawals(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawals", reflect.TypeOf((*MockBeaconClient)(nil).GetWithdrawals), arg0)
}

// KnownValidators mocks base method.
func (m *MockBeaconClient) KnownValidators(arg0 structs.Slot) (client.AllValidatorsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KnownValidators", arg0)
	ret0, _ := ret[0].(client.AllValidatorsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KnownValidators indicates an expected call of KnownValidators.
func (mr *MockBeaconClientMockRecorder) KnownValidators(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KnownValidators", reflect.TypeOf((*MockBeaconClient)(nil).KnownValidators), arg0)
}

// PublishBlock mocks base method.
func (m *MockBeaconClient) PublishBlock(arg0 *types.SignedBeaconBlock) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishBlock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishBlock indicates an expected call of PublishBlock.
func (mr *MockBeaconClientMockRecorder) PublishBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishBlock", reflect.TypeOf((*MockBeaconClient)(nil).PublishBlock), arg0)
}

// SubscribeToHeadEvents mocks base method.
func (m *MockBeaconClient) SubscribeToHeadEvents(arg0 context.Context, arg1 chan client.HeadEvent) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SubscribeToHeadEvents", arg0, arg1)
}

// SubscribeToHeadEvents indicates an expected call of SubscribeToHeadEvents.
func (mr *MockBeaconClientMockRecorder) SubscribeToHeadEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToHeadEvents", reflect.TypeOf((*MockBeaconClient)(nil).SubscribeToHeadEvents), arg0, arg1)
}

// SyncStatus mocks base method.
func (m *MockBeaconClient) SyncStatus() (*client.SyncStatusPayloadData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncStatus")
	ret0, _ := ret[0].(*client.SyncStatusPayloadData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SyncStatus indicates an expected call of SyncStatus.
func (mr *MockBeaconClientMockRecorder) SyncStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncStatus", reflect.TypeOf((*MockBeaconClient)(nil).SyncStatus))
}

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// Duties mocks base method.
func (m *MockState) Duties() structs.DutiesState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Duties")
	ret0, _ := ret[0].(structs.DutiesState)
	return ret0
}

// Duties indicates an expected call of Duties.
func (mr *MockStateMockRecorder) Duties() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Duties", reflect.TypeOf((*MockState)(nil).Duties))
}

// HeadSlot mocks base method.
func (m *MockState) HeadSlot() structs.Slot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeadSlot")
	ret0, _ := ret[0].(structs.Slot)
	return ret0
}

// HeadSlot indicates an expected call of HeadSlot.
func (mr *MockStateMockRecorder) HeadSlot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeadSlot", reflect.TypeOf((*MockState)(nil).HeadSlot))
}

// KnownValidators mocks base method.
func (m *MockState) KnownValidators() structs.ValidatorsState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KnownValidators")
	ret0, _ := ret[0].(structs.ValidatorsState)
	return ret0
}

// KnownValidators indicates an expected call of KnownValidators.
func (mr *MockStateMockRecorder) KnownValidators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KnownValidators", reflect.TypeOf((*MockState)(nil).KnownValidators))
}

// KnownValidatorsUpdateTime mocks base method.
func (m *MockState) KnownValidatorsUpdateTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KnownValidatorsUpdateTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// KnownValidatorsUpdateTime indicates an expected call of KnownValidatorsUpdateTime.
func (mr *MockStateMockRecorder) KnownValidatorsUpdateTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KnownValidatorsUpdateTime", reflect.TypeOf((*MockState)(nil).KnownValidatorsUpdateTime))
}

// SetDuties mocks base method.
func (m *MockState) SetDuties(arg0 structs.DutiesState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDuties", arg0)
}

// SetDuties indicates an expected call of SetDuties.
func (mr *MockStateMockRecorder) SetDuties(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDuties", reflect.TypeOf((*MockState)(nil).SetDuties), arg0)
}

// SetGenesis mocks base method.
func (m *MockState) SetGenesis(arg0 structs.GenesisInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGenesis", arg0)
}

// SetGenesis indicates an expected call of SetGenesis.
func (mr *MockStateMockRecorder) SetGenesis(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGenesis", reflect.TypeOf((*MockState)(nil).SetGenesis), arg0)
}

// SetHeadSlot mocks base method.
func (m *MockState) SetHeadSlot(arg0 structs.Slot) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHeadSlot", arg0)
}

// SetHeadSlot indicates an expected call of SetHeadSlot.
func (mr *MockStateMockRecorder) SetHeadSlot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeadSlot", reflect.TypeOf((*MockState)(nil).SetHeadSlot), arg0)
}

// SetKnownValidators mocks base method.
func (m *MockState) SetKnownValidators(arg0 structs.ValidatorsState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetKnownValidators", arg0)
}

// SetKnownValidators indicates an expected call of SetKnownValidators.
func (mr *MockStateMockRecorder) SetKnownValidators(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKnownValidators", reflect.TypeOf((*MockState)(nil).SetKnownValidators), arg0)
}

// SetReady mocks base method.
func (m *MockState) SetReady() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetReady")
}

// SetReady indicates an expected call of SetReady.
func (mr *MockStateMockRecorder) SetReady() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReady", reflect.TypeOf((*MockState)(nil).SetReady))
}

// SetWithdrawals mocks base method.
func (m *MockState) SetWithdrawals(arg0 structs.WithdrawalsState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetWithdrawals", arg0)
}

// SetWithdrawals indicates an expected call of SetWithdrawals.
func (mr *MockStateMockRecorder) SetWithdrawals(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWithdrawals", reflect.TypeOf((*MockState)(nil).SetWithdrawals), arg0)
}

// Withdrawals mocks base method.
func (m *MockState) Withdrawals() structs.WithdrawalsState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdrawals")
	ret0, _ := ret[0].(structs.WithdrawalsState)
	return ret0
}

// Withdrawals indicates an expected call of Withdrawals.
func (mr *MockStateMockRecorder) Withdrawals() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdrawals", reflect.TypeOf((*MockState)(nil).Withdrawals))
}
