package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/service"
)

type {{ .FileName }}Handler interface {
	Get{{ .FileName }}(ctx *gin.Context)
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

func (h *{{ .FileNameTitleLower }}Handler) Get{{ .FileName }}(ctx *gin.Context) {

}
