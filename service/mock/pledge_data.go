package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.org/cclose/rsi-pledge-track/model"
	"time"
)

type PledgeDataServiceMock struct {
	mock.Mock
}

func (m *PledgeDataServiceMock) Insert(ctx context.Context, pd *model.PledgeData) error {
	args := m.Mock.Called(ctx, pd)
	return args.Error(0)
}

func (m *PledgeDataServiceMock) Get(id int) (*model.PledgeData, error) {
	args := m.Mock.Called(id)
	return args.Get(0).(*model.PledgeData), args.Error(1)
}

func (m *PledgeDataServiceMock) GetAll(offset int, limit int) ([]*model.PledgeData, error) {
	args := m.Mock.Called(offset, limit)
	return args.Get(0).([]*model.PledgeData), args.Error(1)
}

func (m *PledgeDataServiceMock) GetByTimestamp(time time.Time, offset int) (*model.PledgeData, error) {
	args := m.Mock.Called(time, offset)
	return args.Get(0).(*model.PledgeData), args.Error(1)
}

func (m *PledgeDataServiceMock) GetAfterTimestamp(time time.Time, offset int, limit int) ([]*model.PledgeData, error) {
	args := m.Mock.Called(time, offset, limit)
	return args.Get(0).([]*model.PledgeData), args.Error(1)
}
