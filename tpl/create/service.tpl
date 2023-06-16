package service

import (
	"{{ .ProjectName }}/internal/model"
	"{{ .ProjectName }}/internal/repository"
)

type {{ .FileName }}Service interface {
	Get{{ .FileName }}ById(id int64) (*model.{{ .FileName }}, error)
}

type {{ .FileNameTitleLower }}Service struct {
	*Service
	{{ .FileNameTitleLower }}Repository repository.{{ .FileName }}Repository
}

func New{{ .FileName }}Service(service *Service, {{ .FileNameTitleLower }}Repository repository.{{ .FileName }}Repository) {{ .FileName }}Service {
	return &{{ .FileNameTitleLower }}Service{
		Service:        service,
		{{ .FileNameTitleLower }}Repository: {{ .FileNameTitleLower }}Repository,
	}
}

func (s *{{ .FileNameTitleLower }}Service) Get{{ .FileName }}ById(id int64) (*model.{{ .FileName }}, error) {
	return s.{{ .FileNameTitleLower }}Repository.FirstById(id)
}
