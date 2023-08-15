package repository

import (
	"{{ .ProjectName }}/internal/model"
)

type {{ .FileName }}Repository interface {
	FirstById(id int64) (*model.{{ .FileName }}, error)
}

func New{{ .FileName }}Repository(repository *Repository) {{ .FileName }}Repository {
	return &{{ .FileNameTitleLower }}Repository{
		Repository: repository,
	}
}

type {{ .FileNameTitleLower }}Repository struct {
	*Repository
}

func (r *{{ .FileNameTitleLower }}Repository) FirstById(id int64) (*model.{{ .FileName }}, error) {
	var {{ .FileNameTitleLower }} model.{{ .FileName }}
	// TODO: query db
	return &{{ .FileNameTitleLower }}, nil
}
