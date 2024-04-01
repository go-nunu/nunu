# Model

在一个 Go 语言项目中使用 GORM (GORM) 时，通常会将数据库表映射到 Go 的结构体。这些结构体定义了数据模型，用于表示数据库中的表和表之间的关系。在 GORM 项目中，通常会将这些结构体定义在一个名为 model 的目录中。

## 核心用途
这个目录的作用主要是组织和管理与数据库模型相关的代码。将所有的模型文件放置在一个目录下可以更好地组织代码，并使代码结构清晰易于维护。通常，每个模型都会对应一个结构体，该结构体的字段与数据库表的列相对应。

除了结构体定义之外，model 目录中的文件可能还包括与模型相关的其他代码，如模型的方法、模型之间的关系定义等。


下面是一个简单的 User 结构体定义的例子：
```go
type User struct {
	Id        uint   `gorm:"primarykey"`
	UserId    string `gorm:"unique;not null"`
	Nickname  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
```