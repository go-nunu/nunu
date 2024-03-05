package model

import "gorm.io/gorm"

type {{ .StructName }} struct {
	gorm.Model
}

func (m *{{ .StructName }}) TableName() string {
    return "{{ .StructNameSnakeCase }}"
}
