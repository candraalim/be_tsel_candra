package model

import (
	"context"
	"time"
)

type ReferralCode struct {
	ID          int64     `db:"id"`
	Msisdn      string    `db:"msisdn"`
	Code        string    `db:"code"`
	CreatedDate time.Time `db:"created_date"`
	Status      int       `db:"status"`
}

type ReferralCodeRepository interface {
	FindByMsisdn(ctx context.Context, msisdn string) (ReferralCode, error)
	FindByCode(ctx context.Context, code string) (ReferralCode, error)
	Insert(ctx context.Context, code *ReferralCode) error
}
