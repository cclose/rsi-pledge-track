package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"github.org/cclose/rsi-pledge-track/model"
	"time"
)

const pledgeDataSchema = `public`
const pledgeDataTable = pledgeDataSchema + `.starcitizen_pledgedata`

const pledgeDataCol = "id as id, pledge_timestamp as pledge_timestamp, COALESCE(funding, 0) as funding, COALESCE(citizens, 0) as citizens, COALESCE(fleet, 0) as fleet"
const getPledgeDataBase = `SELECT ` + pledgeDataCol + ` FROM ` + pledgeDataTable

// const insertQueryText = "INSERT INTO PledgeData (TimeStamp, Funding, Citizens, Fleet) VALUES($1, $2, $3, $4)"
const insertQueryText = `INSERT INTO ` + pledgeDataTable + ` (pledge_timestamp, funding, citizens, fleet) VALUES (:pledge_timestamp, :funding, :citizens, :fleet) RETURNING id`

type IPledgeDataService interface {
	Insert(ctx context.Context, pd *model.PledgeData) error
	Get(id int) (*model.PledgeData, error)
	GetAll(offset int, limit int) ([]*model.PledgeData, error)
	GetByTimestamp(time time.Time, offset int) (*model.PledgeData, error)
	GetAfterTimestamp(time time.Time, offset int, limit int) ([]*model.PledgeData, error)
}

type PledgeDataService struct {
	DB *sqlx.DB
}

func NewPledgeDataService(db *sqlx.DB) *PledgeDataService {
	return &PledgeDataService{
		DB: db,
	}
}

func (s *PledgeDataService) Insert(ctx context.Context, pd *model.PledgeData) error {
	// need to use query because we have a RETURNING clause
	rows, err := s.DB.NamedQueryContext(ctx, insertQueryText, pd)
	if err != nil {
		return err
	}

	// Fetch the last inserted ID returned by the RETURNING clause
	if rows.Next() {
		var lastInsertedID int
		err = rows.Scan(&lastInsertedID)
		if err != nil {
			return err
		}

		// Assign the last inserted ID to the model
		pd.ID = lastInsertedID
	} else {
		return sql.ErrNoRows
	}

	return nil
}

const getPledgeDataQuery = getPledgeDataBase + ` WHERE id = $1`

func (s *PledgeDataService) Get(id int) (*model.PledgeData, error) {
	logger.Info("calling Get")
	var pledgeData model.PledgeData
	err := s.DB.Get(&pledgeData, getPledgeDataQuery, id)

	return &pledgeData, err
}

func (s *PledgeDataService) GetAll(offset int, limit int) ([]*model.PledgeData, error) {
	logger.Info("calling GetAll")
	var pledgeDataList []*model.PledgeData
	query := getPledgeDataBase
	var args []interface{}
	if offset != 0 {
		args = append(args, offset)
		query = fmt.Sprintf("%s OFFSET $%d", query, len(args))
	}
	if limit != 0 {
		args = append(args, limit)
		query = fmt.Sprintf("%s LIMIT $%d", query, len(args))
	}
	logger.Infof("query %s | args=%v", query, args)
	err := s.DB.Select(&pledgeDataList, query, args...)

	return pledgeDataList, err
}

func (s *PledgeDataService) GetByTimestamp(time time.Time, offset int) (*model.PledgeData, error) {
	logger.Debug("calling GetByTimestamp")
	var pledgeData model.PledgeData
	args := []interface{}{time}
	query := getPledgeDataBase
	if offset != 0 {
		query = fmt.Sprintf("%s WHERE pledge_timestamp = ($%d + INTERVAL '$%d HOUR')", query, len(args), len(args)+1)
		args = append(args, offset)
	} else {
		query = fmt.Sprintf("%s WHERE pledge_timestamp = $%d", query, len(args))
	}
	err := s.DB.Select(&pledgeData, query, args...)

	return &pledgeData, err
}

func (s *PledgeDataService) GetAfterTimestamp(time time.Time, offset int, limit int) ([]*model.PledgeData, error) {
	logger.Debug("calling AfterTimestamp")
	var pledgeDataList []*model.PledgeData
	args := []interface{}{time}
	query := getPledgeDataBase
	if offset != 0 {
		query = fmt.Sprintf("%s WHERE pledge_timestamp >= ($%d + INTERVAL '$%d HOUR')", query, len(args), len(args)+1)
		args = append(args, offset)
	} else {
		query = fmt.Sprintf("%s WHERE pledge_timestamp >= $%d", query, len(args))
	}
	if limit != 0 {
		args = append(args, limit)
		query = fmt.Sprintf(`%s LIMIT $%d`, query, len(args))
	}
	err := s.DB.Select(&pledgeDataList, query, args...)

	return pledgeDataList, err
}
