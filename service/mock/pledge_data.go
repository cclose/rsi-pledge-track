package mock

import (
	"github.com/stretchr/testify/mock"
	"github.org/cclose/rsi-pledge-track/model"
	"time"
)

type PledgeDataServiceMock struct {
	mock.Mock
}

func (m *PledgeDataServiceMock) Insert(pd *model.PledgeData) error {
	args := m.Mock.Called(pd)
	return args.Error(0)
}

func (m *PledgeDataServiceMock) Get(id int, offset int) (*model.PledgeData, error) {
	args := m.Mock.Called(id, offset)
	return args.Get(0).(*model.PledgeData), args.Error(1)
}

func (m *PledgeDataServiceMock) GetAll(offset int, limit int) ([]*model.PledgeData, error) {
	args := m.Mock.Called(limit, offset)
	return args.Get(0).([]*model.PledgeData), args.Error(1)
}

func (m *PledgeDataServiceMock) GetByTimestamp(time time.Time, offset int) (*model.PledgeData, error) {
	args := m.Mock.Called(time, offset)
	return args.Get(0).(*model.PledgeData), args.Error(1)
}

func (m *PledgeDataServiceMock) GetAfterTimestamp(time time.Time, offset int, limit int) ([]*model.PledgeData, error) {
	args := m.Mock.Called(time, limit, offset)
	return args.Get(0).([]*model.PledgeData), args.Error(1)
}
