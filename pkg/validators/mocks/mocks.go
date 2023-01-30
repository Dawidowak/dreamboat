// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blocknative/dreamboat/pkg/validators (interfaces: RegistrationManager,Verifier,State)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	structs "github.com/blocknative/dreamboat/pkg/structs"
	validators "github.com/blocknative/dreamboat/pkg/validators"
	verify "github.com/blocknative/dreamboat/pkg/verify"
	types "github.com/flashbots/go-boost-utils/types"
	gomock "github.com/golang/mock/gomock"
)

// MockRegistrationManager is a mock of RegistrationManager interface.
type MockRegistrationManager struct {
	ctrl     *gomock.Controller
	recorder *MockRegistrationManagerMockRecorder
}

// MockRegistrationManagerMockRecorder is the mock recorder for MockRegistrationManager.
type MockRegistrationManagerMockRecorder struct {
	mock *MockRegistrationManager
}

// NewMockRegistrationManager creates a new mock instance.
func NewMockRegistrationManager(ctrl *gomock.Controller) *MockRegistrationManager {
	mock := &MockRegistrationManager{ctrl: ctrl}
	mock.recorder = &MockRegistrationManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegistrationManager) EXPECT() *MockRegistrationManagerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockRegistrationManager) Check(arg0 *types.RegisterValidatorRequestMessage) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockRegistrationManagerMockRecorder) Check(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockRegistrationManager)(nil).Check), arg0)
}

// GetRegistration mocks base method.
func (m *MockRegistrationManager) GetRegistration(arg0 context.Context, arg1 types.PublicKey) (types.SignedValidatorRegistration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegistration", arg0, arg1)
	ret0, _ := ret[0].(types.SignedValidatorRegistration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRegistration indicates an expected call of GetRegistration.
func (mr *MockRegistrationManagerMockRecorder) GetRegistration(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegistration", reflect.TypeOf((*MockRegistrationManager)(nil).GetRegistration), arg0, arg1)
}

// SendStore mocks base method.
func (m *MockRegistrationManager) SendStore(arg0 validators.StoreReq) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendStore", arg0)
}

// SendStore indicates an expected call of SendStore.
func (mr *MockRegistrationManagerMockRecorder) SendStore(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendStore", reflect.TypeOf((*MockRegistrationManager)(nil).SendStore), arg0)
}

// MockVerifier is a mock of Verifier interface.
type MockVerifier struct {
	ctrl     *gomock.Controller
	recorder *MockVerifierMockRecorder
}

// MockVerifierMockRecorder is the mock recorder for MockVerifier.
type MockVerifierMockRecorder struct {
	mock *MockVerifier
}

// NewMockVerifier creates a new mock instance.
func NewMockVerifier(ctrl *gomock.Controller) *MockVerifier {
	mock := &MockVerifier{ctrl: ctrl}
	mock.recorder = &MockVerifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVerifier) EXPECT() *MockVerifierMockRecorder {
	return m.recorder
}

// GetVerifyChan mocks base method.
func (m *MockVerifier) GetVerifyChan(arg0 uint) chan verify.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifyChan", arg0)
	ret0, _ := ret[0].(chan verify.Request)
	return ret0
}

// GetVerifyChan indicates an expected call of GetVerifyChan.
func (mr *MockVerifierMockRecorder) GetVerifyChan(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifyChan", reflect.TypeOf((*MockVerifier)(nil).GetVerifyChan), arg0)
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

// Beacon mocks base method.
func (m *MockState) Beacon() *structs.BeaconState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Beacon")
	ret0, _ := ret[0].(*structs.BeaconState)
	return ret0
}

// Beacon indicates an expected call of Beacon.
func (mr *MockStateMockRecorder) Beacon() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Beacon", reflect.TypeOf((*MockState)(nil).Beacon))
}
