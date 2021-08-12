package referral

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReferHandler struct {
	useCase ReferUseCase
}

func SetupReferHandler(useCase ReferUseCase) *ReferHandler {
	if useCase == nil {
		panic("referral use case is nil")
	}
	return &ReferHandler{
		useCase: useCase,
	}
}

func (h ReferHandler) ProcessReferral(e echo.Context) error {
	var request ReferRequest
	if err := e.Bind(&request); err != nil {
		return err
	}

	if err := e.Validate(&request); err != nil {
		return err
	}

	res, err := h.useCase.ProcessReferral(e.Request().Context(), request)
	//handle error response
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, res)
}
