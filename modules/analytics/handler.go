package analytics

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

// NewAnalyticsHandler provides a new analytics Handler instance
func NewAnalyticsHandler(i *do.Injector) (Handler, error) {
	return &analyticsHandler{
		service: do.MustInvoke[Service](i),
	}, nil
}

func (a *analyticsHandler) Get(ctx echo.Context) error {
	data, err := a.service.GetAnalytics()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, data)
}

func (a *analyticsHandler) WebRequestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if err := a.service.IncrementWebRequests(); err != nil {
			ctx.Logger().Error(err)
		}
		return next(ctx)
	}
}
