package controllers

import (
	"backend-ta/internal/services"
	internal_err "backend-ta/pkg/errors"
	"backend-ta/pkg/http/server/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatisticsController struct {
	statsSrv services.StatisticsService
}

func NewStatisticsController(statsSrv services.StatisticsService) StatisticsController {
	return StatisticsController{statsSrv: statsSrv}
}

func (ctl *StatisticsController) GetDashboardData(ctx *gin.Context) {
	res, err := ctl.statsSrv.GetDashboardData(ctx)
	if err != nil {
		http_response.SendError(ctx, internal_err.NewDefaultError(http.StatusInternalServerError, err.Error()))
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success fetching dashboard data", res)
}
