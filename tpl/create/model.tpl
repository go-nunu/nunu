package model

import "gorm.io/gorm"

type {{ .StructName }} struct {
	gorm.Model
}

func ({{ .StructNameFirstChar }} *{{ .StructName }}) TableName() string {
    return "{{ .StructNameSnakeCase }}"
}
