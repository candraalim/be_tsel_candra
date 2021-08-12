package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

type referralCodeRepository struct {
	db *Database
}

func SetupReferralCodeRepository(db *Database) *referralCodeRepository {
	if db == nil {
		panic("postgresql db is nil")
	}
	return &referralCodeRepository{
		db: db,
	}
}

const (
	queryReferralCodeFindByMsisdn = "SELECT id, msisdn, code, status FROM referral_code WHERE msisdn = $1"
	queryReferralCodeFindByCode   = "SELECT id, msisdn, code, status FROM referral_code WHERE code = $1"
	queryReferralCodeInsert       = "INSERT INTO %s.referral_code (msisdn, code) VALUES ($1, $2) RETURNING id"
)

func (r referralCodeRepository) FindByMsisdn(ctx context.Context, msisdn string) (result model.ReferralCode, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = r.db.GetContext(ctx, &result, queryReferralCodeFindByMsisdn, msisdn)
	if err != nil {
		return model.ReferralCode{}, err
	}
	if result.Status == 0 {
		return model.ReferralCode{}, util.ErrorDataNotFound
	}
	return result, nil
}

func (r referralCodeRepository) FindByCode(ctx context.Context, code string) (result model.ReferralCode, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = r.db.GetContext(ctx, &result, queryReferralCodeFindByCode, code)
	if err != nil {
		return model.ReferralCode{}, err
	}
	if result.Status == 0 {
		return model.ReferralCode{}, util.ErrorDataNotFound
	}
	return result, nil
}

func (r referralCodeRepository) Insert(ctx context.Context, code *model.ReferralCode) error {
	err := r.db.GetContext(ctx, &code.ID, fmt.Sprintf(queryReferralCodeInsert, r.db.SchemaName()), code.Msisdn, code.Code)
	if err != nil {
		return err
	}
	if code.ID == 0 {
		return util.ErrorDatabase
	}
	return nil
}
