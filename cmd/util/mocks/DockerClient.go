// Code generated by MockGen. DO NOT EDIT.
// Source: ./docker_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"
	time "time"

	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	network "github.com/docker/docker/api/types/network"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// MockDockerClient is a mock of DockerClient interface.
type MockDockerClient struct {
	ctrl     *gomock.Controller
	recorder *MockDockerClientMockRecorder
}

// MockDockerClientMockRecorder is the mock recorder for MockDockerClient.
type MockDockerClientMockRecorder struct {
	mock *MockDockerClient
}

// NewMockDockerClient creates a new mock instance.
func NewMockDockerClient(ctrl *gomock.Controller) *MockDockerClient {
	mock := &MockDockerClient{ctrl: ctrl}
	mock.recorder = &MockDockerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDockerClient) EXPECT() *MockDockerClientMockRecorder {
	return m.recorder
}

// ContainerCreate mocks base method.
func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerCreate", ctx, config, hostConfig, networkingConfig, platform, containerName)
	ret0, _ := ret[0].(container.ContainerCreateCreatedBody)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerCreate indicates an expected call of ContainerCreate.
func (mr *MockDockerClientMockRecorder) ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerCreate", reflect.TypeOf((*MockDockerClient)(nil).ContainerCreate), ctx, config, hostConfig, networkingConfig, platform, containerName)
}

// ContainerList mocks base method.
func (m *MockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerList", ctx, options)
	ret0, _ := ret[0].([]types.Container)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerList indicates an expected call of ContainerList.
func (mr *MockDockerClientMockRecorder) ContainerList(ctx, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerList", reflect.TypeOf((*MockDockerClient)(nil).ContainerList), ctx, options)
}

// ContainerRemove mocks base method.
func (m *MockDockerClient) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerRemove", ctx, containerID, options)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerRemove indicates an expected call of ContainerRemove.
func (mr *MockDockerClientMockRecorder) ContainerRemove(ctx, containerID, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerRemove", reflect.TypeOf((*MockDockerClient)(nil).ContainerRemove), ctx, containerID, options)
}

// ContainerStart mocks base method.
func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerStart", ctx, containerID, options)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStart indicates an expected call of ContainerStart.
func (mr *MockDockerClientMockRecorder) ContainerStart(ctx, containerID, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStart", reflect.TypeOf((*MockDockerClient)(nil).ContainerStart), ctx, containerID, options)
}

// ContainerStop mocks base method.
func (m *MockDockerClient) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContainerStop", ctx, containerID, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStop indicates an expected call of ContainerStop.
func (mr *MockDockerClientMockRecorder) ContainerStop(ctx, containerID, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStop", reflect.TypeOf((*MockDockerClient)(nil).ContainerStop), ctx, containerID, timeout)
}

// ImageBuild mocks base method.
func (m *MockDockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImageBuild", ctx, buildContext, options)
	ret0, _ := ret[0].(types.ImageBuildResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageBuild indicates an expected call of ImageBuild.
func (mr *MockDockerClientMockRecorder) ImageBuild(ctx, buildContext, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageBuild", reflect.TypeOf((*MockDockerClient)(nil).ImageBuild), ctx, buildContext, options)
}

// ImagePull mocks base method.
func (m *MockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImagePull", ctx, refStr, options)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImagePull indicates an expected call of ImagePull.
func (mr *MockDockerClientMockRecorder) ImagePull(ctx, refStr, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImagePull", reflect.TypeOf((*MockDockerClient)(nil).ImagePull), ctx, refStr, options)
}
