package service

import (
	"{{ .ProjectName }}/internal/model"
	"{{ .ProjectName }}/internal/repository"
)

type {{ .StructName }}Service interface {
	Get{{ .StructName }}(id int64) (*model.{{ .StructName }}, error)
}

func New{{ .StructName }}Service(service *Service, {{ .StructNameLowerFirst }}Repository repository.{{ .StructName }}Repository) {{ .StructName }}Service {
	return &{{ .StructNameLowerFirst }}Service{
		Service:        service,
		{{ .StructNameLowerFirst }}Repository: {{ .StructNameLowerFirst }}Repository,
	}
}

type {{ .StructNameLowerFirst }}Service struct {
	*Service
	{{ .StructNameLowerFirst }}Repository repository.{{ .StructName }}Repository
}

func (s *{{ .StructNameLowerFirst }}Service) Get{{ .StructName }}(id int64) (*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.FirstById(id)
}
