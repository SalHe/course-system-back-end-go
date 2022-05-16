package services

import (
	"errors"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Public struct{}

func (p Public) Login(l req.LoginCredit) (user *dao.User, err error) {
	if err = dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ? OR username = ?", l.Username, l.Username).First(&user).Error; err != nil {
		return nil, err
	}

	if user.ComparePassword(l.Password) {
		return user, nil
	}
	return nil, ErrWrongPassword
}

func (p Public) Register(b req.RegisterInfo) (bool, error) {
	var user dao.User
	if err := dao.DB.Unscoped().Model(&dao.User{}).Where("id = ? OR username = ?", b.Id, b.Username).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user = dao.User{
			Model: dao.Model{
				ID: b.Id,
			},
			Username: b.Username,
		}
		user.SetPassword(b.Password)
		if err = dao.DB.Create(&user).Error; err != nil {
			return false, ErrUnknown
		}
		return true, nil
	} else {
		return false, ErrConflict
	}
}
