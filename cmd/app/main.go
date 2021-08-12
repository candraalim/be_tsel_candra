package main

import (
	"github.com/candraalim/be_tsel_candra/config"
	"github.com/candraalim/be_tsel_candra/internal/storage/postgresql"
	"github.com/candraalim/be_tsel_candra/internal/transport/http"
	"github.com/candraalim/be_tsel_candra/internal/usecase/inquiry"
	"github.com/candraalim/be_tsel_candra/internal/usecase/referral"
)

func main() {
	cfg := config.LoadFile()

	db := postgresql.NewDatabase(cfg.Database)
	codeRepo := postgresql.SetupReferralCodeRepository(db)
	historyRepo := postgresql.SetupReferralHistoryRepository(db)
	rewardRepo := postgresql.SetupRewardRepository(db)

	inquiryUseCase := inquiry.SetupInquiryUseCase(codeRepo, historyRepo, rewardRepo)
	inquiryHandler := inquiry.SetupInquiringHandler(inquiryUseCase)

	referralUseCase := referral.SetupReferUseCase(codeRepo, historyRepo)
	referralHandler := referral.SetupReferHandler(referralUseCase)

	http.StartHttpService(cfg.Server, cfg.Auth, inquiryHandler, referralHandler)
}
