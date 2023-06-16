package repository

import (
	"{{ .ProjectName }}/internal/model"
)

type {{ .FileName }}Repository interface {
	FirstById(id int64) (*model.{{ .FileName }}, error)
}
type {{ .FileNameTitleLower }}Repository struct {
	*Repository
}

func New{{ .FileName }}Repository(repository *Repository) {{ .FileName }}Repository {
	return &{{ .FileNameTitleLower }}Repository{
		Repository: repository,
	}
}

func (r *{{ .FileNameTitleLower }}Repository) FirstById(id int64) (*model.{{ .FileName }}, error) {
	var {{ .FileNameTitleLower }} model.{{ .FileName }}
	// TODO: query db
	return &{{ .FileNameTitleLower }}, nil
}
