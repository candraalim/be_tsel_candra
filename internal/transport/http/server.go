package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"

	"github.com/candraalim/be_tsel_candra/config"
	"github.com/candraalim/be_tsel_candra/internal/usecase/inquiry"
	"github.com/candraalim/be_tsel_candra/internal/usecase/referral"
	"github.com/candraalim/be_tsel_candra/internal/util"
)

func StartHttpService(config *config.ServerConfig, auth *config.AuthConfig, inquiring *inquiry.InquiringHandler, referral *referral.ReferHandler) {
	server := echo.New()
	server.HideBanner = true

	setupMiddleware(server)
	setupRouter(server, auth, inquiring, referral)

	// start server
	go func() {
		if err := server.Start(config.AppAddress()); err != nil {
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
	}
}

func setupMiddleware(server *echo.Echo) {
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin,
			echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentEncoding,
			echo.HeaderAcceptEncoding, echo.HeaderXCSRFToken},
		ExposeHeaders:    []string{echo.HeaderContentLength, echo.HeaderAccessControlAllowOrigin},
		AllowCredentials: true,
	}))

	server.HTTPErrorHandler = errorHandler

	server.Validator = &DataValidator{ValidatorData: validator.New()}
}

type DataValidator struct {
	ValidatorData *validator.Validate
}

func (cv *DataValidator) Validate(i interface{}) error {
	return cv.ValidatorData.Struct(i)
}

func errorHandler(err error, c echo.Context) {
	fmt.Println(c.Request().RequestURI)
	code := util.ErrorGeneral

	if he, ok := err.(*util.ApplicationError); ok {
		code = he
	} else if he, ok := err.(*echo.HTTPError); ok {
		fmt.Println(he.Error())
		code.HttpStatus = he.Code
		code.Message = err.Error()
	} else {
		fmt.Println(err.Error())
	}

	_ = c.JSON(code.HttpStatus, code)
	return
}
