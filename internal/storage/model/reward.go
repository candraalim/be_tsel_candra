package model

import (
	"context"
	"time"
)

type Reward struct {
	ID            int64     `db:"id"`
	TotalReferral int       `db:"total_referral"`
	Description   string    `db:"reward_description"`
	CreatedDate   time.Time `db:"created_date"`
	UpdatedDate   time.Time `db:"updated_date"`
	Status        int       `db:"status"`
}

type RewardRepository interface {
	FindByTotalReferral(ctx context.Context, totalReferral int) (Reward, error)
	Insert(ctx context.Context, model *Reward) (err error)
	Update(ctx context.Context, model Reward) (err error)
	Delete(ctx context.Context, ID int64) (err error)
}
