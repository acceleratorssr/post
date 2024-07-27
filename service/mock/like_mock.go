// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/like.go
//
// Generated by this command:
//
//	mockgen -source=./service/like.go -destination=./service/mock/like_mock.go --package=svcmocks
//

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	domain "post/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLikeService is a mock of LikeService interface.
type MockLikeService struct {
	ctrl     *gomock.Controller
	recorder *MockLikeServiceMockRecorder
}

// MockLikeServiceMockRecorder is the mock recorder for MockLikeService.
type MockLikeServiceMockRecorder struct {
	mock *MockLikeService
}

// NewMockLikeService creates a new mock instance.
func NewMockLikeService(ctrl *gomock.Controller) *MockLikeService {
	mock := &MockLikeService{ctrl: ctrl}
	mock.recorder = &MockLikeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLikeService) EXPECT() *MockLikeServiceMockRecorder {
	return m.recorder
}

// Collect mocks base method.
func (m *MockLikeService) Collect(ctx context.Context, ObjType string, ObjID, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Collect", ctx, ObjType, ObjID, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Collect indicates an expected call of Collect.
func (mr *MockLikeServiceMockRecorder) Collect(ctx, ObjType, ObjID, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockLikeService)(nil).Collect), ctx, ObjType, ObjID, uid)
}

// GetListAllOfLikes mocks base method.
func (m *MockLikeService) GetListBatchOfLikes(ctx context.Context, ObjType string, offset, limit int) ([]domain.Like, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListBatchOfLikes", ctx, ObjType, offset, limit)
	ret0, _ := ret[0].([]domain.Like)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListAllOfLikes indicates an expected call of GetListAllOfLikes.
func (mr *MockLikeServiceMockRecorder) GetListAllOfLikes(ctx, ObjType, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListBatchOfLikes", reflect.TypeOf((*MockLikeService)(nil).GetListBatchOfLikes), ctx, ObjType, offset, limit)
}

// IncrReadCount mocks base method.
func (m *MockLikeService) IncrReadCount(ctx context.Context, ObjType string, ObjID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrReadCount", ctx, ObjType, ObjID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrReadCount indicates an expected call of IncrReadCount.
func (mr *MockLikeServiceMockRecorder) IncrReadCount(ctx, ObjType, ObjID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrReadCount", reflect.TypeOf((*MockLikeService)(nil).IncrReadCount), ctx, ObjType, ObjID)
}

// Like mocks base method.
func (m *MockLikeService) Like(ctx context.Context, ObjType string, ObjID, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Like", ctx, ObjType, ObjID, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Like indicates an expected call of Like.
func (mr *MockLikeServiceMockRecorder) Like(ctx, ObjType, ObjID, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Like", reflect.TypeOf((*MockLikeService)(nil).Like), ctx, ObjType, ObjID, uid)
}

// UnLike mocks base method.
func (m *MockLikeService) UnLike(ctx context.Context, ObjType string, ObjID, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLike", ctx, ObjType, ObjID, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnLike indicates an expected call of UnLike.
func (mr *MockLikeServiceMockRecorder) UnLike(ctx, ObjType, ObjID, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLike", reflect.TypeOf((*MockLikeService)(nil).UnLike), ctx, ObjType, ObjID, uid)
}
