package service

import (
	"{{ .ProjectName }}/internal/dao"
	"{{ .ProjectName }}/internal/model"
)

type {{ .FileName }}Service struct {
	*Service
	{{ .FileNameTitleLower }}Dao *dao.{{ .FileName }}Dao
}

func New{{ .FileName }}Service(service *Service, {{ .FileNameTitleLower }}Dao *dao.{{ .FileName }}Dao) *{{ .FileName }}Service {
	return &{{ .FileName }}Service{
		Service: service,
		{{ .FileNameTitleLower }}Dao: {{ .FileNameTitleLower }}Dao,
	}
}

func (s *{{ .FileName }}Service) Get{{ .FileName }}ById(id int64) (*model.{{ .FileName }}, error) {
	return s.{{ .FileNameTitleLower }}Dao.FirstById(id)
}