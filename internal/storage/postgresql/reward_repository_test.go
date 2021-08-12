package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
)

func setupStub(t *testing.T) (*Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub database", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	database := &Database{
		DB: sqlxDB,
	}
	return database, mock
}

func TestSetupRewardRepository(t *testing.T) {
	assert.Panics(t, func() {
		SetupRewardRepository(nil)
	})

	assert.NotPanics(t, func() {
		db, _ := setupStub(t)
		SetupRewardRepository(db)
	})
}

func Test_rewardRepository_FindByTotalReferral(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)reward*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupRewardRepository(db)
		_, err := r.FindByTotalReferral(context.Background(), 10)
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)reward*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "total_referral", "reward_description", "status"}).
				AddRow(3, 6, "bonus 20 GB", 1))

		r := SetupRewardRepository(db)
		result, err := r.FindByTotalReferral(context.Background(), 10)
		assert.Nil(t, err)
		assert.Equal(t, model.Reward{
			ID:            3,
			TotalReferral: 6,
			Description:   "bonus 20 GB",
			Status:        1,
		}, result)
	})
}

func Test_rewardRepository_Delete(t *testing.T) {
	t.Run("error context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectExec("^UPDATE (.+)reward*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupRewardRepository(db)
		err := r.Delete(context.Background(), 11)
		assert.NotNil(t, err)
	})
	t.Run("data not found", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectExec("^UPDATE (.+)reward*").
			WillReturnResult(sqlmock.NewResult(0, 0))

		r := SetupRewardRepository(db)
		err := r.Delete(context.Background(), 11)
		assert.NotNil(t, err)
	})
	t.Run("success delete data", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectExec("^UPDATE (.+)reward*").
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := SetupRewardRepository(db)
		err := r.Delete(context.Background(), 11)
		assert.Nil(t, err)
	})
}

func Test_rewardRepository_Insert(t *testing.T) {
	t.Run("error context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)reward*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupRewardRepository(db)
		err := r.Insert(context.Background(), &model.Reward{TotalReferral: 3, Description: "bonus 3 GB"})
		assert.NotNil(t, err)
	})
	t.Run("failed insert, id is 0", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)reward*").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(0))

		r := SetupRewardRepository(db)
		err := r.Insert(context.Background(), &model.Reward{TotalReferral: 3, Description: "bonus 3 GB"})
		assert.NotNil(t, err)
	})
	t.Run("success insert data", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)reward*").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(11))

		r := SetupRewardRepository(db)
		data := &model.Reward{TotalReferral: 3, Description: "bonus 3 GB"}
		err := r.Insert(context.Background(), data)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), data.ID)
	})
}

func Test_rewardRepository_Update(t *testing.T) {
	t.Run("error context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectExec("^UPDATE (.+)reward*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupRewardRepository(db)
		err := r.Update(context.Background(), model.Reward{ID: 2})
		assert.NotNil(t, err)
	})
	t.Run("data not found", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectExec("^UPDATE (.+)reward*").
			WillReturnResult(sqlmock.NewResult(0, 0))

		r := SetupRewardRepository(db)
		err := r.Update(context.Background(), model.Reward{ID: 2, Description: "bonus 3 GB"})
		assert.NotNil(t, err)
	})
	t.Run("success delete data", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectExec("^UPDATE (.+)reward*").
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := SetupRewardRepository(db)
		err := r.Update(context.Background(), model.Reward{ID: 2, TotalReferral: 4})
		assert.Nil(t, err)
	})
}
