package model

import "gorm.io/gorm"

type {{ .FileName }} struct {
	gorm.Model
}

func (m *{{ .FileName }}) TableName() string {
    return "{{ .FileNameTitleLower }}"
}
