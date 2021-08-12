package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
)

func TestSetupReferralHistoryRepository(t *testing.T) {
	assert.Panics(t, func() {
		SetupReferralHistoryRepository(nil)
	})

	assert.NotPanics(t, func() {
		db, _ := setupStub(t)
		SetupReferralHistoryRepository(db)
	})
}

func Test_referralHistoryRepository_CountByMsisdn(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT COUNT\\(id\\) FROM referral_history*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralHistoryRepository(db)
		_, err := r.CountByMsisdn(context.Background(), "082100000")
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT COUNT\\(id\\) FROM referral_history*").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(13))

		r := SetupReferralHistoryRepository(db)
		result, err := r.CountByMsisdn(context.Background(), "082100000")
		assert.Nil(t, err)
		assert.Equal(t, 13, result)
	})
}

func Test_referralHistoryRepository_FindByMsisdn(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_history*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralHistoryRepository(db)
		_, err := r.FindByMsisdn(context.Background(), "082100000", 0, 10)
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_history*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "msisdn", "code", "msisdn_referee", "referral_date"}).
				AddRow(13, "082100000", "ABS123AD12", "08210001000", "2021-08-11").
				AddRow(11, "082100000", "ABS123AD12", "08210001002", "2021-08-10"))

		r := SetupReferralHistoryRepository(db)
		result, err := r.FindByMsisdn(context.Background(), "082100000", 0, 10)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))
	})
}

func Test_referralHistoryRepository_FindByMsisdnReferee(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_history*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralHistoryRepository(db)
		_, err := r.FindByMsisdnReferee(context.Background(), "082100000")
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_history*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "msisdn", "code", "msisdn_referee", "referral_date"}).
				AddRow(3, "082100001", "ABS123AD12", "082100000", "2021-08-11"))

		r := SetupReferralHistoryRepository(db)
		result, err := r.FindByMsisdnReferee(context.Background(), "082100000")
		assert.Nil(t, err)
		assert.Equal(t, model.ReferralHistory{
			ID:            3,
			Msisdn:        "082100001",
			Code:          "ABS123AD12",
			MsisdnReferee: "082100000",
			ReferralDate:  "2021-08-11",
		}, result)
	})
}

func Test_referralHistoryRepository_GetTotalByMsisdnAndMonth(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT COUNT\\(id\\) FROM referral_history*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralHistoryRepository(db)
		_, err := r.GetTotalByMsisdnAndMonth(context.Background(), "082100000", "2021-08")
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT COUNT\\(id\\) FROM referral_history*").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(13))

		r := SetupReferralHistoryRepository(db)
		result, err := r.GetTotalByMsisdnAndMonth(context.Background(), "082100000", "2021-08")
		assert.Nil(t, err)
		assert.Equal(t, 13, result)
	})
}

func Test_referralHistoryRepository_Insert(t *testing.T) {
	t.Run("error context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)referral_history*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralHistoryRepository(db)
		err := r.Insert(context.Background(), &model.ReferralHistory{Msisdn: "082100000", Code: "ABS123AD12", MsisdnReferee: "0821000001", ReferralDate: "2021-08-11"})
		assert.NotNil(t, err)
	})
	t.Run("failed insert, id is 0", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)referral_history*").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(0))

		r := SetupReferralHistoryRepository(db)
		err := r.Insert(context.Background(), &model.ReferralHistory{Msisdn: "082100000", Code: "ABS123AD12", MsisdnReferee: "0821000001", ReferralDate: "2021-08-11"})
		assert.NotNil(t, err)
	})
	t.Run("success insert data", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)referral_history*").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(11))

		r := SetupReferralHistoryRepository(db)
		data := &model.ReferralHistory{Msisdn: "082100000", Code: "ABS123AD12", MsisdnReferee: "0821000001", ReferralDate: "2021-08-11"}
		err := r.Insert(context.Background(), data)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), data.ID)
	})
}
