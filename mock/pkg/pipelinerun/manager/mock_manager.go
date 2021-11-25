// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/pipelinerun/manager/manager.go

// Package mock_manager is a generated GoMock package.
package mock_manager

import (
	context "context"
	q "g.hz.netease.com/horizon/lib/q"
	models "g.hz.netease.com/horizon/pkg/pipelinerun/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockManager) Create(ctx context.Context, pipelinerun *models.Pipelinerun) (*models.Pipelinerun, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, pipelinerun)
	ret0, _ := ret[0].(*models.Pipelinerun)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockManagerMockRecorder) Create(ctx, pipelinerun interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManager)(nil).Create), ctx, pipelinerun)
}

// GetByID mocks base method
func (m *MockManager) GetByID(ctx context.Context, pipelinerunID uint) (*models.Pipelinerun, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, pipelinerunID)
	ret0, _ := ret[0].(*models.Pipelinerun)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockManagerMockRecorder) GetByID(ctx, pipelinerunID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockManager)(nil).GetByID), ctx, pipelinerunID)
}

// GetByClusterID mocks base method
func (m *MockManager) GetByClusterID(ctx context.Context, clusterID uint, canRollback bool, query q.Query) (int, []*models.Pipelinerun, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByClusterID", ctx, clusterID, canRollback, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]*models.Pipelinerun)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByClusterID indicates an expected call of GetByClusterID
func (mr *MockManagerMockRecorder) GetByClusterID(ctx, clusterID, canRollback, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByClusterID", reflect.TypeOf((*MockManager)(nil).GetByClusterID), ctx, clusterID, canRollback, query)
}

// DeleteByID mocks base method
func (m *MockManager) DeleteByID(ctx context.Context, pipelinerunID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, pipelinerunID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID
func (mr *MockManagerMockRecorder) DeleteByID(ctx, pipelinerunID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockManager)(nil).DeleteByID), ctx, pipelinerunID)
}

// UpdateConfigCommitByID mocks base method
func (m *MockManager) UpdateConfigCommitByID(ctx context.Context, pipelinerunID uint, commit string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfigCommitByID", ctx, pipelinerunID, commit)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfigCommitByID indicates an expected call of UpdateConfigCommitByID
func (mr *MockManagerMockRecorder) UpdateConfigCommitByID(ctx, pipelinerunID, commit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfigCommitByID", reflect.TypeOf((*MockManager)(nil).UpdateConfigCommitByID), ctx, pipelinerunID, commit)
}

// GetLatestByClusterIDAndAction mocks base method
func (m *MockManager) GetLatestByClusterIDAndAction(ctx context.Context, clusterID uint, action string) (*models.Pipelinerun, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestByClusterIDAndAction", ctx, clusterID, action)
	ret0, _ := ret[0].(*models.Pipelinerun)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestByClusterIDAndAction indicates an expected call of GetLatestByClusterIDAndAction
func (mr *MockManagerMockRecorder) GetLatestByClusterIDAndAction(ctx, clusterID, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestByClusterIDAndAction", reflect.TypeOf((*MockManager)(nil).GetLatestByClusterIDAndAction), ctx, clusterID, action)
}

// GetLatestSuccessByClusterID mocks base method
func (m *MockManager) GetLatestSuccessByClusterID(ctx context.Context, clusterID uint) (*models.Pipelinerun, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestSuccessByClusterID", ctx, clusterID)
	ret0, _ := ret[0].(*models.Pipelinerun)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestSuccessByClusterID indicates an expected call of GetLatestSuccessByClusterID
func (mr *MockManagerMockRecorder) GetLatestSuccessByClusterID(ctx, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestSuccessByClusterID", reflect.TypeOf((*MockManager)(nil).GetLatestSuccessByClusterID), ctx, clusterID)
}

// UpdateResultByID mocks base method
func (m *MockManager) UpdateResultByID(ctx context.Context, pipelinerunID uint, result *models.Result) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateResultByID", ctx, pipelinerunID, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResultByID indicates an expected call of UpdateResultByID
func (mr *MockManagerMockRecorder) UpdateResultByID(ctx, pipelinerunID, result interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateResultByID", reflect.TypeOf((*MockManager)(nil).UpdateResultByID), ctx, pipelinerunID, result)
}