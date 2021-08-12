package inquiry

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/candraalim/be_tsel_candra/internal/storage/model"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

type InquiringUseCase interface {
	GetReferralCode(ctx context.Context, msisdn string) (ReferralCodeResponse, error)
	GetCurrentReferralReward(ctx context.Context, msisdn string) (ReferralRewardResponse, error)
	GetListReferral(ctx context.Context, msisdn string, page, limit int) (ReferralHistoryResponse, error)
}

type inquiryUseCase struct {
	codeRepository    model.ReferralCodeRepository
	historyRepository model.ReferralHistoryRepository
	rewardRepository  model.RewardRepository
}

func SetupInquiryUseCase(referralCodeRepository model.ReferralCodeRepository,
	referralHistoryRepository model.ReferralHistoryRepository,
	rewardRepository model.RewardRepository) InquiringUseCase {
	if referralCodeRepository == nil {
		panic("ReferralCodeRepository is nil")
	}
	if referralHistoryRepository == nil {
		panic("ReferralHistoryRepository is nil")
	}
	if rewardRepository == nil {
		panic("RewardRepository is nil")
	}
	return &inquiryUseCase{
		codeRepository:    referralCodeRepository,
		historyRepository: referralHistoryRepository,
		rewardRepository:  rewardRepository,
	}
}

func (i inquiryUseCase) GetReferralCode(ctx context.Context, msisdn string) (response ReferralCodeResponse, err error) {
	//validate msisdn
	msisdn, err = util.ValidateAndSanitizeMsisdn(msisdn)
	if err != nil {
		return ReferralCodeResponse{}, err
	}

	//lookup data from table referral_code
	result, err := i.codeRepository.FindByMsisdn(ctx, msisdn)
	if err == nil {
		//data found return it immediately
		return ReferralCodeResponse{
			Code:    util.CodeSuccess,
			Message: util.MessageSuccess,
			Data: struct {
				ReferralCode string `json:"referralCode"`
			}{
				ReferralCode: result.Code,
			},
		}, nil
	}
	//if error is not no rows, return database error
	if !errors.Is(err, sql.ErrNoRows) {
		return ReferralCodeResponse{}, err
	}

	//generate new referral code
	referralCode, err := i.generateReferralCode(ctx)
	if err != nil {
		return ReferralCodeResponse{}, err
	}
	//store to db
	err = i.codeRepository.Insert(context.Background(), &model.ReferralCode{
		Msisdn: msisdn,
		Code:   referralCode,
	})
	if err != nil {
		return ReferralCodeResponse{}, err
	}
	return ReferralCodeResponse{
		Code:    util.CodeSuccess,
		Message: util.MessageSuccess,
		Data: struct {
			ReferralCode string `json:"referralCode"`
		}{
			ReferralCode: referralCode,
		},
	}, nil
}

//todo move length referral code to config
func (i inquiryUseCase) generateReferralCode(ctx context.Context) (referralCode string, err error) {
	// generate unique referral code, make sure it unique by lookup to db
	// if return duplicate then try to generate it 10 times, after that return error
	// retry in next request
	for c := 0; c < 10; c++ {
		code, er := randomHex(5)
		if er != nil {
			log.Println("failed to generate random", er.Error())
			return "", util.ErrorGenerateReferralCode
		}
		if _, er := i.codeRepository.FindByCode(ctx, code); errors.Is(er, sql.ErrNoRows) {
			referralCode = code
			break
		} else if er != nil {
			break
		}
	}
	if referralCode == "" {
		log.Println("unable to generate unique referral code")
		return referralCode, util.ErrorGenerateReferralCode
	}
	return referralCode, nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func (i inquiryUseCase) GetCurrentReferralReward(ctx context.Context, msisdn string) (resp ReferralRewardResponse, err error) {
	//validate msisdn
	msisdn, err = util.ValidateAndSanitizeMsisdn(msisdn)
	if err != nil {
		return ReferralRewardResponse{}, err
	}

	//todo get total from redis
	//get total referral
	month := time.Now().Format("2006-01")
	total, err := i.historyRepository.GetTotalByMsisdnAndMonth(ctx, msisdn, month)
	if err != nil {
		return ReferralRewardResponse{}, err
	}
	resp = ReferralRewardResponse{
		Code:    util.CodeSuccess,
		Message: util.MessageSuccess,
		Data: struct {
			TotalReferral int    `json:"totalReferral"`
			Reward        string `json:"reward"`
		}{
			TotalReferral: total,
		},
	}
	//when no referral happen, return success but empty reward
	if total == 0 {
		return resp, nil
	}

	//lookup reward based on total referral
	reward, err := i.rewardRepository.FindByTotalReferral(ctx, total)
	if err != nil {
		return ReferralRewardResponse{}, err
	}
	resp.Data.Reward = reward.Description
	return resp, nil
}

func (i inquiryUseCase) GetListReferral(ctx context.Context, msisdn string, page, limit int) (resp ReferralHistoryResponse, err error) {
	msisdn, err = util.ValidateAndSanitizeMsisdn(msisdn)
	if err != nil {
		return ReferralHistoryResponse{}, err
	}

	var (
		entities []model.ReferralHistory
		total    int
	)
	if page < 1 {
		page = 1
	}

	eg := errgroup.Group{}
	eg.Go(func() error {
		var er error
		offset := 0
		if page > 1 {
			offset = (page - 1) * limit
		}
		if limit < 1 {
			limit = 10
		}
		entities, er = i.historyRepository.FindByMsisdn(ctx, msisdn, offset, limit)
		return er
	})
	eg.Go(func() error {
		var er error
		total, er = i.historyRepository.CountByMsisdn(ctx, msisdn)
		return er
	})

	if err = eg.Wait(); err != nil {
		return ReferralHistoryResponse{}, err
	}

	return assemblerReferralHistory(entities, total, page, limit), nil
}

func assemblerReferralHistory(entities []model.ReferralHistory, total, page, limit int) ReferralHistoryResponse {
	var totalPage int
	if total > 0 {
		totalPage = ((total - (total % limit)) / limit) + 1
	}

	list := make([]ReferralHistory, len(entities))
	for i, v := range entities {
		list[i] = ReferralHistory{
			Msisdn:       v.MsisdnReferee,
			ReferralDate: v.ReferralDate,
			DateTime:     v.CreatedDate.UnixNano() / 1000000,
		}
	}

	return ReferralHistoryResponse{
		Code:    util.CodeSuccess,
		Message: util.MessageSuccess,
		Data: ReferralHistoryData{
			List: list,
			Meta: Meta{
				TotalPage:   totalPage,
				TotalRecord: total,
				Page:        page,
				Size:        len(entities),
				Limit:       limit,
				FirstPage:   page == 1 || page == 0,
				LastPage:    page == totalPage,
			},
		},
	}
}
