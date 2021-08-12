package referral

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

type ReferUseCase interface {
	ProcessReferral(ctx context.Context, request ReferRequest) (ReferResponse, error)
}

type referUseCase struct {
	codeRepository    model.ReferralCodeRepository
	historyRepository model.ReferralHistoryRepository
}

func SetupReferUseCase(referralCodeRepository model.ReferralCodeRepository,
	referralHistoryRepository model.ReferralHistoryRepository) ReferUseCase {
	if referralCodeRepository == nil {
		panic("ReferralCodeRepository is nil")
	}
	if referralHistoryRepository == nil {
		panic("ReferralHistoryRepository is nil")
	}
	return &referUseCase{
		codeRepository:    referralCodeRepository,
		historyRepository: referralHistoryRepository,
	}
}

func (r referUseCase) ProcessReferral(ctx context.Context, request ReferRequest) (resp ReferResponse, err error) {
	//validate msisdn
	request.Msisdn, err = util.ValidateAndSanitizeMsisdn(request.Msisdn)
	if err != nil {
		return ReferResponse{}, err
	}

	//todo change to config
	if len(request.Code) > 20 {
		log.Println("invalid code length")
		return ReferResponse{}, util.ErrorInvalidRequest
	}

	referralCode, err := r.codeRepository.FindByCode(ctx, request.Code)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("referral code not found: ", request.Code)
		return ReferResponse{}, util.ErrorInvalidRequest
	}
	if err != nil {
		return ReferResponse{}, err
	}

	history, err := r.historyRepository.FindByMsisdnReferee(ctx, request.Msisdn)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return ReferResponse{}, err
	}
	if history.ID > 0 {
		log.Println("msisdn already register with other referral: ", request.Msisdn)
		return ReferResponse{}, util.ErrorInvalidRequest
	}

	if code, _ := r.codeRepository.FindByMsisdn(ctx, request.Msisdn); code.ID > 0 {
		log.Println("msisdn already register: ", request.Msisdn)
		return ReferResponse{}, util.ErrorInvalidRequest
	}

	history = model.ReferralHistory{
		Msisdn:        referralCode.Msisdn,
		Code:          request.Code,
		ReferralDate:  time.Now().Format("2006-01-02"),
		MsisdnReferee: request.Msisdn,
	}
	err = r.historyRepository.Insert(ctx, &history)
	if err != nil {
		return ReferResponse{}, err
	}
	//TODO store counter to redis
	return ReferResponse{Code: util.CodeSuccess, Message: util.MessageSuccess}, nil
}
