package handler

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/internal/service"
)

type {{ .FileName }}Handler struct {
	*Handler
	{{ .FileNameTitleLower }}Service service.{{ .FileName }}Service
}

func New{{ .FileName }}Handler(handler *Handler, {{ .FileNameTitleLower }}Service service.{{ .FileName }}Service) *{{ .FileName }}Handler {
	return &{{ .FileName }}Handler{
		Handler:      handler,
		{{ .FileNameTitleLower }}Service: {{ .FileNameTitleLower }}Service,
	}
}

func (h *{{ .FileName }}Handler) Get{{ .FileName }}(ctx *gin.Context) {

}
