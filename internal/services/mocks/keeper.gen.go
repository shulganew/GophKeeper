// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/keeper.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
	entities "github.com/shulganew/GophKeeper/internal/entities"
)

// MockKeeperer is a mock of Keeperer interface.
type MockKeeperer struct {
	ctrl     *gomock.Controller
	recorder *MockKeepererMockRecorder
}

// MockKeepererMockRecorder is the mock recorder for MockKeeperer.
type MockKeepererMockRecorder struct {
	mock *MockKeeperer
}

// NewMockKeeperer creates a new mock instance.
func NewMockKeeperer(ctrl *gomock.Controller) *MockKeeperer {
	mock := &MockKeeperer{ctrl: ctrl}
	mock.recorder = &MockKeepererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeeperer) EXPECT() *MockKeepererMockRecorder {
	return m.recorder
}

// AddSite mocks base method.
func (m *MockKeeperer) AddSite(ctx context.Context, site entities.NewSecretEncoded) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSite", ctx, site)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSite indicates an expected call of AddSite.
func (mr *MockKeepererMockRecorder) AddSite(ctx, site interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSite", reflect.TypeOf((*MockKeeperer)(nil).AddSite), ctx, site)
}

// AddUser mocks base method.
func (m *MockKeeperer) AddUser(ctx context.Context, login, hash, email string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, login, hash, email)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockKeepererMockRecorder) AddUser(ctx, login, hash, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockKeeperer)(nil).AddUser), ctx, login, hash, email)
}

// GetByLogin mocks base method.
func (m *MockKeeperer) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLogin", ctx, login)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLogin indicates an expected call of GetByLogin.
func (mr *MockKeepererMockRecorder) GetByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLogin", reflect.TypeOf((*MockKeeperer)(nil).GetByLogin), ctx, login)
}

// GetSites mocks base method.
func (m *MockKeeperer) GetSites(ctx context.Context, userID string, stype entities.SecretType) ([]entities.SecretEncoded, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSites", ctx, userID, stype)
	ret0, _ := ret[0].([]entities.SecretEncoded)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSites indicates an expected call of GetSites.
func (mr *MockKeepererMockRecorder) GetSites(ctx, userID, stype interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSites", reflect.TypeOf((*MockKeeperer)(nil).GetSites), ctx, userID, stype)
}

// LoadEKeysc mocks base method.
func (m *MockKeeperer) LoadEKeysc(ctx context.Context) ([]entities.EKeyDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadEKeysc", ctx)
	ret0, _ := ret[0].([]entities.EKeyDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadEKeysc indicates an expected call of LoadEKeysc.
func (mr *MockKeepererMockRecorder) LoadEKeysc(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadEKeysc", reflect.TypeOf((*MockKeeperer)(nil).LoadEKeysc), ctx)
}

// SaveEKeyc mocks base method.
func (m *MockKeeperer) SaveEKeyc(ctx context.Context, eKeyc entities.EKeyDB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEKeyc", ctx, eKeyc)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEKeyc indicates an expected call of SaveEKeyc.
func (mr *MockKeepererMockRecorder) SaveEKeyc(ctx, eKeyc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEKeyc", reflect.TypeOf((*MockKeeperer)(nil).SaveEKeyc), ctx, eKeyc)
}

// SaveEKeysc mocks base method.
func (m *MockKeeperer) SaveEKeysc(ctx context.Context, eKeysc []entities.EKeyDB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveEKeysc", ctx, eKeysc)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveEKeysc indicates an expected call of SaveEKeysc.
func (mr *MockKeepererMockRecorder) SaveEKeysc(ctx, eKeysc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveEKeysc", reflect.TypeOf((*MockKeeperer)(nil).SaveEKeysc), ctx, eKeysc)
}
