package services

import (
	"errors"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct{}

func (u User) GetUser(id uint) (*dao.User, error) {
	var user dao.User

	if err := dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u User) NewUser(b req.NewUserRequest) error {
	var user *dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ? OR username = ? OR username = ?", b.Id, b.Username, b.Username, b.Id).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		user := &dao.User{
			Model:        dao.Model{ID: b.Id},
			Username:     b.Username,
			RealName:     b.RealName,
			Role:         b.Role,
			CollegeId:    b.CollegeId,
			EntranceYear: b.EntranceYear,
		}
		user.SetPassword(b.Password)
		result := dao.DB.Create(user)
		if err := result.Error; err != nil {
			return ErrUnknown
		}
		return nil
	} else {
		return ErrConflict
	}
}

func (u User) GetUserList(pageInfo req.Page) ([]dao.User, error) {
	var users []dao.User
	if err := dao.DB.Preload(clause.Associations).Offset(pageInfo.Offset()).Limit(pageInfo.ActualSize()).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u User) DeleteUser(id uint) error {
	var user dao.User
	var count int64
	if dao.DB.Model(&dao.User{}).Find(&user, "id = ?", id).Count(&count); count <= 0 {
		return ErrNotFound
	}
	if err := dao.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
		return ErrUnknown
	}
	return nil
}

func (u User) UpdateUser(id uint, b req.UpdateUserRequest) (*resp.User, error) {
	var user dao.User
	if err := dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// 暂时就先返回 resp.User 吧
	oldUserInfo := resp.NewUser(&user)

	user.Username = b.Username
	user.RealName = b.RealName
	user.Role = b.Role
	user.CollegeId = b.CollegeId
	user.EntranceYear = b.EntranceYear
	if err := dao.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return oldUserInfo, nil
}
