package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

type rewardRepository struct {
	db *Database
}

func SetupRewardRepository(db *Database) *rewardRepository {
	if db == nil {
		panic("postgresql db is nil")
	}
	return &rewardRepository{
		db: db,
	}
}

const (
	queryRewardFindByTotalReferral = `SELECT id, total_referral, reward_description, status FROM reward 
									  WHERE total_referral >= $1 AND status = 1 ORDER BY total_referral ASC LIMIT 1 `
	queryRewardInsert     = "INSERT INTO %s.reward(total_referral, reward_description) VALUES ($1, $2) RETURNING id"
	queryRewardSoftDelete = "UPDATE %s.reward SET status = 0 WHERE id = $1"
)

func (r rewardRepository) FindByTotalReferral(ctx context.Context, totalReferral int) (result model.Reward, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = r.db.GetContext(ctx, &result, queryRewardFindByTotalReferral, totalReferral)
	return result, err
}

func (r rewardRepository) Insert(ctx context.Context, model *model.Reward) (err error) {
	err = r.db.GetContext(ctx, &model.ID, fmt.Sprintf(queryRewardInsert, r.db.SchemaName()), model.TotalReferral, model.Description)
	if err != nil {
		return err
	}
	if model.ID == 0 {
		return util.ErrorDatabase
	}
	return nil
}

func (r rewardRepository) Update(ctx context.Context, model model.Reward) (err error) {
	qb := strings.Builder{}
	qb.WriteString("UPDATE ")
	qb.WriteString(r.db.SchemaName())
	qb.WriteString(".reward SET ")

	var cols []string
	var args []interface{}
	if model.TotalReferral > 0 {
		args = append(args, model.TotalReferral)
		cols = append(cols, fmt.Sprintf("total_referral=$%d", len(args)))
	}
	if model.Description != "" {
		args = append(args, model.Description)
		cols = append(cols, fmt.Sprintf("reward_description=$%d", len(args)))
	}
	cols = append(cols, "updated_date=NOW()")

	qb.WriteString(strings.Join(cols, ","))

	args = append(args, model.ID)
	qb.WriteString(fmt.Sprintf(" WHERE id=$%d", len(args)))

	result, err := r.db.ExecContext(ctx, qb.String(), args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return util.ErrorDataNotFound
	}
	return nil
}

func (r rewardRepository) Delete(ctx context.Context, ID int64) (err error) {
	result, err := r.db.ExecContext(ctx, fmt.Sprintf(queryRewardSoftDelete, r.db.SchemaName()), ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return util.ErrorDataNotFound
	}
	return nil
}
