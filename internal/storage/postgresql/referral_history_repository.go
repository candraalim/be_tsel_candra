package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

type referralHistoryRepository struct {
	db *Database
}

func SetupReferralHistoryRepository(db *Database) *referralHistoryRepository {
	if db == nil {
		panic("postgresql db is nil")
	}
	return &referralHistoryRepository{
		db: db,
	}
}

const (
	queryHistoryFindByMsisdn = `SELECT id, msisdn, code, referral_date, msisdn_referee, referral_date, created_date FROM referral_history 
								WHERE msisdn = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	queryHistoryCountByMsisdn = "SELECT COUNT(id) FROM referral_history WHERE msisdn = $1"
	queryHistoryFindByReferee = `SELECT id, msisdn, code, referral_date, msisdn_referee, referral_date, created_date FROM referral_history
								 WHERE msisdn_referee = $1`
	queryHistoryTotalMonthByMsisdn = "SELECT COUNT(id) FROM referral_history WHERE msisdn = $1 AND referral_date LIKE $2"
	queryHistoryInsert             = `INSERT INTO %s.referral_history (msisdn, code, referral_date, msisdn_referee) 
									  VALUES ($1, $2, $3, $4) RETURNING id`
)

func (r referralHistoryRepository) FindByMsisdn(ctx context.Context, msisdn string, offset, limit int) (result []model.ReferralHistory, err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err = r.db.SelectContext(ctx, &result, queryHistoryFindByMsisdn, msisdn, limit, offset)
	return result, err
}

func (r referralHistoryRepository) CountByMsisdn(ctx context.Context, msisdn string) (total int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err = r.db.GetContext(ctx, &total, queryHistoryCountByMsisdn, msisdn)
	return total, err
}

func (r referralHistoryRepository) FindByMsisdnReferee(ctx context.Context, msisdnReferee string) (result model.ReferralHistory, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = r.db.GetContext(ctx, &result, queryHistoryFindByReferee, msisdnReferee)
	return result, err
}

func (r referralHistoryRepository) GetTotalByMsisdnAndMonth(ctx context.Context, msisdn, month string) (total int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	err = r.db.GetContext(ctx, &total, queryHistoryTotalMonthByMsisdn, msisdn, month+"%")
	return total, err
}

func (r referralHistoryRepository) Insert(ctx context.Context, referral *model.ReferralHistory) error {
	err := r.db.GetContext(ctx, &referral.ID, fmt.Sprintf(queryHistoryInsert, r.db.SchemaName()), referral.Msisdn,
		referral.Code, referral.ReferralDate, referral.MsisdnReferee)
	if err != nil {
		return err
	}
	if referral.ID == 0 {
		return util.ErrorDatabase
	}
	return nil
}
