package repository

import (
	"{{ .ProjectName }}/internal/model"
)

type {{ .StructName }}Repository interface {
	FirstById(id int64) (*model.{{ .StructName }}, error)
}

func New{{ .StructName }}Repository(repository *Repository) {{ .StructName }}Repository {
	return &{{ .StructNameLowerFirst }}Repository{
		Repository: repository,
	}
}

type {{ .StructNameLowerFirst }}Repository struct {
	*Repository
}

func (r *{{ .StructNameLowerFirst }}Repository) FirstById(id int64) (*model.{{ .StructName }}, error) {
	var {{ .StructNameLowerFirst }} model.{{ .StructName }}
	// TODO: query db
	return &{{ .StructNameLowerFirst }}, nil
}
