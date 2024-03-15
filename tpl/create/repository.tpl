package repository

import (
    "context"
	"{{ .ProjectName }}/internal/model"
)

type {{ .StructName }}Repository interface {
	FirstBy(ctx context.Context, query interface{}, args ...interface{}) (*model.{{ .StructName }}, error)
}

func New{{ .StructName }}Repository(repository *Repository) {{ .StructName }}Repository {
	return &{{ .StructNameLowerFirst }}Repository{
		Repository: repository,
	}
}

type {{ .StructNameLowerFirst }}Repository struct {
	*Repository
}

func (r *{{ .StructNameLowerFirst }}Repository) FirstBy(ctx context.Context, query interface{}, args ...interface{}) (*model.{{ .StructName }}, error) {
	var {{ .StructNameLowerFirst }} model.{{ .StructName }}
	// eg: db.Where(query, args...).First(&{{ .StructNameLowerFirst }}).Error
	return &{{ .StructNameLowerFirst }}, nil
}