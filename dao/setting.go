package dao

type Setting struct {
	Model
	Key   string `gorm:"not null,unique"`
	Value string
}
