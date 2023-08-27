package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/pkg/response"
	"{{ .ProjectName }}/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type {{ .FileName }}Handler interface {
	Get{{ .FileName }}ById(ctx *gin.Context)
}

func New{{ .FileName }}Handler(handler *Handler, {{ .FileNameTitleLower }}Service service.{{ .FileName }}Service) {{ .FileName }}Handler {
	return &{{ .FileNameTitleLower }}Handler{
		Handler:      handler,
		{{ .FileNameTitleLower }}Service: {{ .FileNameTitleLower }}Service,
	}
}

type {{ .FileNameTitleLower }}Handler struct {
	*Handler
	{{ .FileNameTitleLower }}Service service.{{ .FileName }}Service
}

func (h *{{ .FileNameTitleLower }}Handler) Get{{ .FileName }}ById(ctx *gin.Context) {
	var params struct {
		Id int64 `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, nil)
		return
	}

	{{ .FileNameTitleLower }}, err := h.{{ .FileNameTitleLower }}Service.Get{{ .FileName }}ById(params.Id)
	h.logger.Info("Get{{ .FileName }}ByID", zap.Any("{{ .FileNameTitleLower }}", {{ .FileNameTitleLower }}))
	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, nil)
		return
	}
	response.HandleSuccess(ctx, {{ .FileNameTitleLower }})
}
