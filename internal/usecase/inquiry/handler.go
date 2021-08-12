package inquiry

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type InquiringHandler struct {
	useCase InquiringUseCase
}

func SetupInquiringHandler(useCase InquiringUseCase) *InquiringHandler {
	if useCase == nil {
		panic("inquiring use case is nil")
	}
	return &InquiringHandler{
		useCase: useCase,
	}
}

func (h InquiringHandler) GetReferralCode(e echo.Context) error {
	msisdn := e.Param("msisdn")

	res, err := h.useCase.GetReferralCode(e.Request().Context(), msisdn)
	//handle error response
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, res)
}

func (h InquiringHandler) GetListReferral(e echo.Context) error {
	msisdn := e.Param("msisdn")

	var page, limit int
	if pg, err := strconv.Atoi(e.QueryParam("page")); err == nil {
		page = pg
	}
	if lim, err := strconv.Atoi(e.QueryParam("limit")); err == nil {
		limit = lim
	}

	res, err := h.useCase.GetListReferral(e.Request().Context(), msisdn, page, limit)
	//handle error response
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, res)
}

func (h InquiringHandler) GetCurrentReferralReward(e echo.Context) error {
	msisdn := e.Param("msisdn")

	res, err := h.useCase.GetCurrentReferralReward(e.Request().Context(), msisdn)
	//handle error response
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, res)
}
