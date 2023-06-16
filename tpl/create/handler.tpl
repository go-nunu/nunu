package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/service"
	"{{ .ProjectName }}/pkg/helper/resp"
	"go.uber.org/zap"
	"net/http"
)

type {{ .FileName }}Handler interface {
	Get{{ .FileName }}ById(ctx *gin.Context)
	Update{{ .FileName }}(ctx *gin.Context)
}

type {{ .FileNameTitleLower }}Handler struct {
	*Handler
	{{ .FileNameTitleLower }}Service service.{{ .FileName }}Service
}

func New{{ .FileName }}Handler(handler *Handler, {{ .FileNameTitleLower }}Service service.{{ .FileName }}Service) {{ .FileName }}Handler {
	return &{{ .FileNameTitleLower }}Handler{
		Handler:     handler,
		{{ .FileNameTitleLower }}Service: {{ .FileNameTitleLower }}Service,
	}
}

func (h *{{ .FileNameTitleLower }}Handler) Get{{ .FileName }}ById(ctx *gin.Context) {
	var params struct {
		Id int64 `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		resp.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	{{ .FileNameTitleLower }}, err := h.{{ .FileNameTitleLower }}Service.Get{{ .FileName }}ById(params.Id)
	h.logger.Info("Get{{ .FileName }}ByID", zap.Any("{{ .FileNameTitleLower }}", {{ .FileNameTitleLower }}))
	if err != nil {
		resp.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	resp.HandleSuccess(ctx, {{ .FileNameTitleLower }})
}

func (h *{{ .FileNameTitleLower }}Handler) Update{{ .FileName }}(ctx *gin.Context) {
	resp.HandleSuccess(ctx, nil)
}
