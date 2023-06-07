package model

import "gorm.io/gorm"

type {{ .FileName }} struct {
	gorm.Model
}

// func ({{ .FileNameFirstChar }} *{{ .FileName }}) TableName() string {
// 	return "{{ .FileNameTitleLower }}"
// }
