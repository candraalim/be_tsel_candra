package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
)

func TestSetupReferralCodeRepository(t *testing.T) {
	assert.Panics(t, func() {
		SetupReferralCodeRepository(nil)
	})

	assert.NotPanics(t, func() {
		db, _ := setupStub(t)
		SetupReferralCodeRepository(db)
	})
}

func Test_referralCodeRepository_FindByCode(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_code*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralCodeRepository(db)
		_, err := r.FindByCode(context.Background(), "ABS123AD12")
		assert.NotNil(t, err)
	})
	t.Run("data deleted", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_code*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "msisdn", "code", "status"}).
				AddRow(3, "082100000", "ABS123AD12", 0))

		r := SetupReferralCodeRepository(db)
		_, err := r.FindByCode(context.Background(), "ABS123AD12")
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_code*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "msisdn", "code", "status"}).
				AddRow(3, "082100000", "ABS123AD12", 1))

		r := SetupReferralCodeRepository(db)
		result, err := r.FindByCode(context.Background(), "ABS123AD12")
		assert.Nil(t, err)
		assert.Equal(t, model.ReferralCode{
			ID:     3,
			Msisdn: "082100000",
			Code:   "ABS123AD12",
			Status: 1,
		}, result)
	})
}

func Test_referralCodeRepository_FindByMsisdn(t *testing.T) {
	t.Run("return context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_code*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralCodeRepository(db)
		_, err := r.FindByMsisdn(context.Background(), "082100000")
		assert.NotNil(t, err)
	})
	t.Run("data deleted", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_code*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "msisdn", "code", "status"}).
				AddRow(3, "082100000", "ABS123AD12", 0))

		r := SetupReferralCodeRepository(db)
		_, err := r.FindByMsisdn(context.Background(), "082100000")
		assert.NotNil(t, err)
	})
	t.Run("data found in db", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^SELECT (.+)referral_code*").
			WillReturnRows(sqlmock.NewRows([]string{"id", "msisdn", "code", "status"}).
				AddRow(3, "082100000", "ABS123AD12", 1))

		r := SetupReferralCodeRepository(db)
		result, err := r.FindByMsisdn(context.Background(), "082100000")
		assert.Nil(t, err)
		assert.Equal(t, model.ReferralCode{
			ID:     3,
			Msisdn: "082100000",
			Code:   "ABS123AD12",
			Status: 1,
		}, result)
	})
}

func Test_referralCodeRepository_Insert(t *testing.T) {
	t.Run("error context deadline exceed", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)referral_code*").
			WillReturnError(context.DeadlineExceeded)

		r := SetupReferralCodeRepository(db)
		err := r.Insert(context.Background(), &model.ReferralCode{Msisdn: "082100000", Code: "ABS123AD12"})
		assert.NotNil(t, err)
	})
	t.Run("failed insert, id is 0", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)referral_code*").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(0))

		r := SetupReferralCodeRepository(db)
		err := r.Insert(context.Background(), &model.ReferralCode{Msisdn: "082100000", Code: "ABS123AD12"})
		assert.NotNil(t, err)
	})
	t.Run("success insert data", func(t *testing.T) {
		db, mock := setupStub(t)
		mock.ExpectQuery("^INSERT INTO (.+)referral_code*").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(11))

		r := SetupReferralCodeRepository(db)
		data := &model.ReferralCode{Msisdn: "082100000", Code: "ABS123AD12"}
		err := r.Insert(context.Background(), data)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), data.ID)
	})
}
