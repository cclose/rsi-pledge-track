package service

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"github.org/cclose/rsi-pledge-track/model"
	"time"
)

// const insertQueryText = "INSERT INTO PledgeData (TimeStamp, Funding, Citizens, Fleet) VALUES($1, $2, $3, $4)"
const insertQueryText = `INSERT INTO PledgeData (Timestamp, Funding, Citizens, Fleet) VALUES (:TimeStamp, :Funding, :Citizens, :Fleet)`

type IPledgeDataService interface {
	Insert(pd *model.PledgeData) error
	Get(id int, offset int) (*model.PledgeData, error)
	GetAll(offset int, limit int) ([]*model.PledgeData, error)
	GetByTimestamp(time time.Time, offset int) (*model.PledgeData, error)
	GetAfterTimestamp(time time.Time, offset int, limit int) ([]*model.PledgeData, error)
}

type PledgeDataService struct {
	DB *sqlx.DB
}

const pledgeDataCol = "ID as id, TimeStamp as timestamp, Funding as funding, Citizens as citizens, Fleet as fleet"

func NewPledgeDataService(db *sqlx.DB) *PledgeDataService {
	return &PledgeDataService{
		DB: db,
	}
}

func (s *PledgeDataService) Insert(pd *model.PledgeData) error {
	result, err := s.DB.NamedExec(insertQueryText, pd)
	if err != nil {
		return err
	}

	// Access the last inserted ID, if needed
	lastInsertedID, _ := result.LastInsertId()
	pd.ID = int(lastInsertedID)

	return nil
}

const getPledgeDataQuery = `SELECT ` + pledgeDataCol + ` FROM PledgeData WHERE ID = $1`

func (s *PledgeDataService) Get(id int, offset int) (*model.PledgeData, error) {
	logger.Info("calling Get")
	var pledgeData model.PledgeData
	err := s.DB.Get(&pledgeData, getPledgeDataQuery, id)

	return &pledgeData, err
}

const getAllPledgeDataQuery = `SELECT ` + pledgeDataCol + ` FROM PledgeData`

func (s *PledgeDataService) GetAll(offset int, limit int) ([]*model.PledgeData, error) {
	logger.Info("calling GetAll")
	var pledgeDataList []*model.PledgeData
	err := s.DB.Select(&pledgeDataList, getAllPledgeDataQuery)

	return pledgeDataList, err
}

func (s *PledgeDataService) GetByTimestamp(time time.Time, offset int) (*model.PledgeData, error) {
	return nil, nil
}

func (s *PledgeDataService) GetAfterTimestamp(time time.Time, offset int, limit int) ([]*model.PledgeData, error) {
	return nil, nil
}
