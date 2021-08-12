package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/candraalim/be_tsel_candra/config"
	"github.com/candraalim/be_tsel_candra/internal/usecase/inquiry"
	"github.com/candraalim/be_tsel_candra/internal/usecase/referral"
)

func setupRouter(server *echo.Echo, auth *config.AuthConfig, inquiring *inquiry.InquiringHandler, referral *referral.ReferHandler) {

	// health check
	server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "services up and running... "+time.Now().Format(time.RFC3339))
	})

	var basicAuthFunc middleware.BasicAuthValidator = func(username, password string, context echo.Context) (b bool, e error) {
		if username == auth.Username && password == auth.Password {
			b = true
			return
		}
		return
	}

	basicAuth := middleware.BasicAuth(basicAuthFunc)

	group := server.Group("/1.0/referral", basicAuth)
	{
		group.GET("/:msisdn/code", inquiring.GetReferralCode)
		group.GET("/:msisdn", inquiring.GetListReferral)
		group.GET("/:msisdn/reward", inquiring.GetCurrentReferralReward)
	}
	{
		group.POST("", referral.ProcessReferral)
	}
}
