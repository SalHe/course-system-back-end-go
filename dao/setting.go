package dao

type Setting struct {
	Model
	Name  string `gorm:"not null,unique"`
	Value string
}
