// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/keeper.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
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

// AddSecretStor mocks base method.
func (m *MockKeeperer) AddSecretStor(ctx context.Context, entity entities.NewSecretEncoded, stype entities.SecretType) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSecretStor", ctx, entity, stype)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSecretStor indicates an expected call of AddSecretStor.
func (mr *MockKeepererMockRecorder) AddSecretStor(ctx, entity, stype interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSecretStor", reflect.TypeOf((*MockKeeperer)(nil).AddSecretStor), ctx, entity, stype)
}

// AddUser mocks base method.
func (m *MockKeeperer) AddUser(ctx context.Context, login, hash, email, otpKey string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, login, hash, email, otpKey)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockKeepererMockRecorder) AddUser(ctx, login, hash, email, otpKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockKeeperer)(nil).AddUser), ctx, login, hash, email, otpKey)
}

// DeleteSecretStor mocks base method.
func (m *MockKeeperer) DeleteSecretStor(ctx context.Context, secretID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecretStor", ctx, secretID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecretStor indicates an expected call of DeleteSecretStor.
func (mr *MockKeepererMockRecorder) DeleteSecretStor(ctx, secretID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecretStor", reflect.TypeOf((*MockKeeperer)(nil).DeleteSecretStor), ctx, secretID)
}

// DropKeys mocks base method.
func (m *MockKeeperer) DropKeys(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DropKeys", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DropKeys indicates an expected call of DropKeys.
func (mr *MockKeepererMockRecorder) DropKeys(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DropKeys", reflect.TypeOf((*MockKeeperer)(nil).DropKeys), ctx)
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

// GetSecretStor mocks base method.
func (m *MockKeeperer) GetSecretStor(ctx context.Context, secretID string) (*entities.SecretEncoded, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretStor", ctx, secretID)
	ret0, _ := ret[0].(*entities.SecretEncoded)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretStor indicates an expected call of GetSecretStor.
func (mr *MockKeepererMockRecorder) GetSecretStor(ctx, secretID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretStor", reflect.TypeOf((*MockKeeperer)(nil).GetSecretStor), ctx, secretID)
}

// GetSecretsStor mocks base method.
func (m *MockKeeperer) GetSecretsStor(ctx context.Context, userID string, stype entities.SecretType) ([]*entities.SecretEncoded, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretsStor", ctx, userID, stype)
	ret0, _ := ret[0].([]*entities.SecretEncoded)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretsStor indicates an expected call of GetSecretsStor.
func (mr *MockKeepererMockRecorder) GetSecretsStor(ctx, userID, stype interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretsStor", reflect.TypeOf((*MockKeeperer)(nil).GetSecretsStor), ctx, userID, stype)
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

// UpdateSecretStor mocks base method.
func (m *MockKeeperer) UpdateSecretStor(ctx context.Context, entity entities.NewSecretEncoded, secretID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecretStor", ctx, entity, secretID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecretStor indicates an expected call of UpdateSecretStor.
func (mr *MockKeepererMockRecorder) UpdateSecretStor(ctx, entity, secretID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecretStor", reflect.TypeOf((*MockKeeperer)(nil).UpdateSecretStor), ctx, entity, secretID)
}

// MockFileKeeper is a mock of FileKeeper interface.
type MockFileKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockFileKeeperMockRecorder
}

// MockFileKeeperMockRecorder is the mock recorder for MockFileKeeper.
type MockFileKeeperMockRecorder struct {
	mock *MockFileKeeper
}

// NewMockFileKeeper creates a new mock instance.
func NewMockFileKeeper(ctrl *gomock.Controller) *MockFileKeeper {
	mock := &MockFileKeeper{ctrl: ctrl}
	mock.recorder = &MockFileKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileKeeper) EXPECT() *MockFileKeeperMockRecorder {
	return m.recorder
}

// DeleteFile mocks base method.
func (m *MockFileKeeper) DeleteFile(ctx context.Context, backet, fileID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, backet, fileID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockFileKeeperMockRecorder) DeleteFile(ctx, backet, fileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockFileKeeper)(nil).DeleteFile), ctx, backet, fileID)
}

// DownloadFile mocks base method.
func (m *MockFileKeeper) DownloadFile(ctx context.Context, backet, fileID string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", ctx, backet, fileID)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockFileKeeperMockRecorder) DownloadFile(ctx, backet, fileID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockFileKeeper)(nil).DownloadFile), ctx, backet, fileID)
}

// UploadFile mocks base method.
func (m *MockFileKeeper) UploadFile(ctx context.Context, backet, fileID string, fr io.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, backet, fileID, fr)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockFileKeeperMockRecorder) UploadFile(ctx, backet, fileID, fr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockFileKeeper)(nil).UploadFile), ctx, backet, fileID, fr)
}
