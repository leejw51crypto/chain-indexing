package ddd_projection_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex/entity/event"
)

type MockProjection struct {
	mock.Mock
}

func NewMockProjection() *MockProjection {
	return &MockProjection{}
}

func (projection *MockProjection) Id() string {
	mockArgs := projection.Called()

	return mockArgs.String(0)
}

func (projection *MockProjection) GetEventsToListen() []string {
	mockArgs := projection.Called()

	return mockArgs.Get(0).([]string)
}

func (projection *MockProjection) GetLastHandledEventHeight() *int64 {
	mockArgs := projection.Called()

	return mockArgs.Get(0).(*int64)
}

func (projection *MockProjection) OnInit() error {
	mockArgs := projection.Called()

	return mockArgs.Error(0)
}

func (projection *MockProjection) HandleEvents(evt []event.Event) error {
	mockArgs := projection.Called(evt)

	return mockArgs.Error(0)
}