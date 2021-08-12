package referral

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/candraalim/be_tsel_candra/internal/mock/storage"
	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

func TestSetupReferUseCase(t *testing.T) {
	assert.Panics(t, func() {
		SetupReferUseCase(nil, nil)
	})
	assert.Panics(t, func() {
		SetupReferUseCase(&mocks.ReferralCodeRepository{}, nil)
	})
	assert.NotPanics(t, func() {
		SetupReferUseCase(&mocks.ReferralCodeRepository{}, &mocks.ReferralHistoryRepository{})
	})
}

func Test_referUseCase_ProcessReferral(t *testing.T) {
	t.Run("invalid msisdn", func(t *testing.T) {
		i := referUseCase{}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Msisdn: "628000abcd",
		})
		assert.NotNil(t, err)
	})
	t.Run("invalid length referral code", func(t *testing.T) {
		i := referUseCase{}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "sjdanlfhsabduhlasuhdbashddasdaa",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("unknown referral code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{}, sql.ErrNoRows)

		i := referUseCase{codeRepository: codeMock}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("error db when validate referral code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{}, context.DeadlineExceeded)

		i := referUseCase{codeRepository: codeMock}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("error db when validate msisdn referee", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{Code: "ABCABC123A",
			Msisdn: "628000001111"}, nil)

		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdnReferee", mock.Anything, mock.Anything).Return(model.ReferralHistory{}, context.DeadlineExceeded)

		i := referUseCase{codeRepository: codeMock, historyRepository: historyMock}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("msisdn referee already refer with other referral code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{Code: "ABCABC123A",
			Msisdn: "628000001111"}, nil)

		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdnReferee", mock.Anything, mock.Anything).Return(model.ReferralHistory{ID: 1122}, nil)

		i := referUseCase{codeRepository: codeMock, historyRepository: historyMock}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("msisdn already found in db referral_code", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{Code: "ABCABC123A",
			Msisdn: "628000001111"}, nil)
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{
			ID:     123,
			Code:   "A12BC92A",
			Msisdn: "6280001100001",
		}, nil)

		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdnReferee", mock.Anything, mock.Anything).Return(model.ReferralHistory{}, sql.ErrNoRows)

		i := referUseCase{codeRepository: codeMock, historyRepository: historyMock}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("db error when insert", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{Code: "ABCABC123A",
			Msisdn: "628000001111"}, nil)
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, nil)

		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdnReferee", mock.Anything, mock.Anything).Return(model.ReferralHistory{}, sql.ErrNoRows)
		historyMock.On("Insert", mock.Anything, mock.Anything).Return(sql.ErrConnDone)

		i := referUseCase{codeRepository: codeMock, historyRepository: historyMock}
		_, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.NotNil(t, err)
	})
	t.Run("success", func(t *testing.T) {
		codeMock := &mocks.ReferralCodeRepository{}
		codeMock.On("FindByCode", mock.Anything, mock.Anything).Return(model.ReferralCode{Code: "ABCABC123A",
			Msisdn: "628000001111"}, nil)
		codeMock.On("FindByMsisdn", mock.Anything, mock.Anything).Return(model.ReferralCode{}, nil)

		historyMock := &mocks.ReferralHistoryRepository{}
		historyMock.On("FindByMsisdnReferee", mock.Anything, mock.Anything).Return(model.ReferralHistory{}, sql.ErrNoRows)
		historyMock.On("Insert", mock.Anything, mock.Anything).Return(nil)

		i := referUseCase{codeRepository: codeMock, historyRepository: historyMock}
		resp, err := i.ProcessReferral(context.Background(), ReferRequest{
			Code:   "ABCABC123A",
			Msisdn: "6280001100001",
		})
		assert.Nil(t, err)
		assert.Equal(t, ReferResponse{
			Code:    util.CodeSuccess,
			Message: util.MessageSuccess,
		}, resp)
	})
}
