package model

import (
	"context"
	"time"
)

type ReferralHistory struct {
	ID            int64     `db:"id"`
	Msisdn        string    `db:"msisdn"`
	Code          string    `db:"code"`
	ReferralDate  string    `db:"referral_date"`
	MsisdnReferee string    `db:"msisdn_referee"`
	CreatedDate   time.Time `db:"created_date"`
}

type ReferralHistoryRepository interface {
	FindByMsisdn(ctx context.Context, msisdn string, offset, limit int) (result []ReferralHistory, err error)
	CountByMsisdn(ctx context.Context, msisdn string) (total int, err error)
	FindByMsisdnReferee(ctx context.Context, msisdnReferee string) (result ReferralHistory, err error)
	GetTotalByMsisdnAndMonth(ctx context.Context, msisdn, month string) (total int, err error)
	Insert(ctx context.Context, referral *ReferralHistory) error
}
