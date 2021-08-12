package inquiry

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/candraalim/be_tsel_candra/internal/mock/storage"
	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

func TestSetupInquiryUseCase(t *testing.T) {
	assert.Panics(t, func() {
		SetupInquiryUseCase(nil, nil, nil)
	})
	assert.Panics(t, func() {
		SetupInquiryUseCase(&mocks.ReferralCodeRepository{}, nil, nil)
	})
	assert.Panics(t, func() {
		SetupInquiryUseCase(&mocks.ReferralCodeRepository{}, &mocks.ReferralHistoryRepository{}, nil)
	})
	assert.NotPanics(t, func() {
		SetupInquiryUseCase(&mocks.ReferralCodeRepository{}, &mocks.ReferralHistoryRepository{}, &mocks.RewardRepository{})
	})
}

func Test_inquiryUseCase_GetCurrentReferralReward(t *testing.T) {
	t.Run("invalid msisdn", func(t *testing.T) {
		i := inquiryUseCase{}
		_, err := i.GetCurrentReferralReward(context.Background(), "080000acbcd")
		assert.NotNil(t, err)
	})
	t.Run("get total referral history return error", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("GetTotalByMsisdnAndMonth", mock.Anything, mock.Anything, mock.Anything).Return(0, context.DeadlineExceeded)

		i := inquiryUseCase{
			historyRepository: historyMock,
		}
		_, err := i.GetCurrentReferralReward(context.Background(), "62821000000")
		assert.NotNil(t, err)
	})
	t.Run("get total referral history 0", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("GetTotalByMsisdnAndMonth", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)

		i := inquiryUseCase{
			historyRepository: historyMock,
		}
		resp, err := i.GetCurrentReferralReward(context.Background(), "62821000000")
		assert.Nil(t, err)
		assert.Equal(t, util.CodeSuccess, resp.Code)
		assert.Equal(t, 0, resp.Data.TotalReferral)
		assert.Empty(t, resp.Data.Reward)
	})
	t.Run("unable to get reward data", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("GetTotalByMsisdnAndMonth", mock.Anything, mock.Anything, mock.Anything).Return(3, nil)

		rewardMock := &mocks.RewardRepository{}
		rewardMock.On("FindByTotalReferral", mock.Anything, mock.Anything).Return(model.Reward{}, context.DeadlineExceeded)

		i := inquiryUseCase{
			historyRepository: historyMock,
			rewardRepository:  rewardMock,
		}
		_, err := i.GetCurrentReferralReward(context.Background(), "62821000000")
		assert.NotNil(t, err)
	})
	t.Run("success get reward", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("GetTotalByMsisdnAndMonth", mock.Anything, mock.Anything, mock.Anything).Return(4, nil)

		rewardMock := &mocks.RewardRepository{}
		rewardMock.On("FindByTotalReferral", mock.Anything, mock.Anything).Return(model.Reward{Description: "bonus 20GB"}, nil)

		i := inquiryUseCase{
			historyRepository: historyMock,
			rewardRepository:  rewardMock,
		}
		resp, err := i.GetCurrentReferralReward(context.Background(), "62821000000")
		assert.Nil(t, err)
		assert.Equal(t, util.CodeSuccess, resp.Code)
		assert.Equal(t, 4, resp.Data.TotalReferral)
		assert.Equal(t, "bonus 20GB", resp.Data.Reward)
	})
}

func Test_inquiryUseCase_GetListReferral(t *testing.T) {
	t.Run("invalid msisdn", func(t *testing.T) {
		i := inquiryUseCase{}
		_, err := i.GetListReferral(context.Background(), "080000acbcd", 0, 0)
		assert.NotNil(t, err)
	})
	t.Run("error get count", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.ReferralHistory{}, nil)
		historyMock.On("CountByMsisdn", mock.Anything, mock.Anything).Return(0, context.DeadlineExceeded)

		i := inquiryUseCase{
			historyRepository: historyMock,
		}
		_, err := i.GetListReferral(context.Background(), "0800001231321", 0, 0)
		assert.NotNil(t, err)
	})
	t.Run("error get list", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.ReferralHistory{}, context.DeadlineExceeded)
		historyMock.On("CountByMsisdn", mock.Anything, mock.Anything).Return(0, nil)

		i := inquiryUseCase{
			historyRepository: historyMock,
		}
		_, err := i.GetListReferral(context.Background(), "0800001231321", 2, 10)
		assert.NotNil(t, err)
	})
	t.Run("return empty list", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.ReferralHistory{}, nil)
		historyMock.On("CountByMsisdn", mock.Anything, mock.Anything).Return(0, nil)

		i := inquiryUseCase{
			historyRepository: historyMock,
		}
		resp, err := i.GetListReferral(context.Background(), "0800001231321", 3, 10)
		assert.Nil(t, err)
		assert.Empty(t, resp.Data.List)
	})
	t.Run("return success", func(t *testing.T) {
		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.ReferralHistory{{
			ID:            123,
			Msisdn:        "62800000002",
			Code:          "AABBCC112233",
			ReferralDate:  "2021-08-12",
			MsisdnReferee: "62800000001",
			CreatedDate:   time.Now(),
		}}, nil)
		historyMock.On("CountByMsisdn", mock.Anything, mock.Anything).Return(1, nil)

		i := inquiryUseCase{
			historyRepository: historyMock,
		}
		resp, err := i.GetListReferral(context.Background(), "0800001231321", 1, 10)
		assert.Nil(t, err)
		assert.NotEmpty(t, resp.Data.List)
	})
}

func Test_inquiryUseCase_GetReferralCode(t *testing.T) {
	t.Run("invalid msisdn", func(t *testing.T) {
		i := inquiryUseCase{}
		_, err := i.GetReferralCode(context.Background(), "080000acbcd")
		assert.NotNil(t, err)
	})
	t.Run("db timeout", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, context.DeadlineExceeded)

		i := inquiryUseCase{
			codeRepository: codeMock,
		}
		_, err := i.GetReferralCode(context.Background(), "080000123131")
		assert.NotNil(t, err)
	})
	t.Run("referral code found in db", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{Code: "AA11BB22CC33"}, nil)

		i := inquiryUseCase{
			codeRepository: codeMock,
		}
		resp, err := i.GetReferralCode(context.Background(), "080000123131")
		assert.Nil(t, err)
		assert.Equal(t, util.CodeSuccess, resp.Code)
		assert.Equal(t, "AA11BB22CC33", resp.Data.ReferralCode)
	})
	t.Run("success generate referral code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)
		codeMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

		i := inquiryUseCase{
			codeRepository: codeMock,
		}
		resp, err := i.GetReferralCode(context.Background(), "080000123131")
		assert.Nil(t, err)
		assert.Equal(t, util.CodeSuccess, resp.Code)
		assert.NotEmpty(t, resp.Data.ReferralCode)
	})
	t.Run("unable to store referrel code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)
		codeMock.On("Insert", mock.Anything, mock.Anything).Return(context.DeadlineExceeded)

		i := inquiryUseCase{
			codeRepository: codeMock,
		}
		_, err := i.GetReferralCode(context.Background(), "080000123131")
		assert.NotNil(t, err)
	})
	t.Run("generate not unique", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{}, nil)
		codeMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

		i := inquiryUseCase{
			codeRepository: codeMock,
		}
		_, err := i.GetReferralCode(context.Background(), "080000123131")
		assert.NotNil(t, err)
	})
	t.Run("error in middle generate referral code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{}, context.DeadlineExceeded)
		codeMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

		i := inquiryUseCase{
			codeRepository: codeMock,
		}
		_, err := i.GetReferralCode(context.Background(), "080000123131")
		assert.NotNil(t, err)
	})
}
