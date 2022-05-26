package services

import (
	"errors"
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type User struct{}

func (u *User) GetUser(id uint) (*dao.User, error) {
	var user dao.User

	if err := dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetUserByUsername(username string) (*dao.User, error) {
	var user dao.User

	if err := dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) NewUser(b req.NewUserRequest) error {
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

func (u *User) GetUserCount() int64 {
	var count int64
	dao.DB.Model(&dao.User{}).Count(&count)
	return count
}

func (u *User) GetUserList(b req.QueryUserRequest) (int64, []dao.User, error) {
	var users []dao.User

	tx := dao.DB.Preload(clause.Associations).Model(&dao.User{})
	if b.Id != 0 {
		tx = tx.Where("id LIKE ?", "%"+strconv.Itoa(int(b.Id))+"%")
	}
	if b.Username != "" {
		tx = tx.Where("username LIKE ?", "%"+b.Username+"%")
	}
	if b.RealName != "" {
		tx = tx.Where("real_name LIKE ?", "%"+b.RealName+"%")
	}
	if len(b.Roles) > 0 {
		tx = tx.Where("role in (?)", b.Roles)
	}
	if len(b.CollegesId) > 0 {
		tx = tx.Where("college_id in (?)", b.CollegesId)
	}
	if b.EntranceYearFrom > 0 {
		tx = tx.Where("entrance_year >= ?", b.EntranceYearFrom)
	}
	if b.EntranceYearTo > 0 {
		tx = tx.Where("entrance_year <= ?", b.EntranceYearTo)
	}
	var count int64
	tx.Count(&count)
	tx = tx.Offset(b.Offset()).Limit(b.ActualSize())
	if err := tx.Find(&users).Error; err != nil {
		return 0, nil, err
	}
	return count, users, nil
}

func (u *User) DeleteUser(id uint) error {
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

func (u *User) UpdateUser(id uint, b req.UpdateUserRequest, operatedByAdmin bool) (*dao.User, error) {
	var user dao.User
	if err := dao.DB.Model(&dao.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// 暂时就先返回 resp.User 吧
	// oldUserInfo := resp.NewUser(&user)

	user.Username = b.Username
	if operatedByAdmin {
		// 目前个人只允许修改用户名
		user.RealName = b.RealName
		if user.Role != dao.RoleAdmin && b.Role != dao.RoleAdmin {
			user.Role = b.Role
		}
		user.CollegeId = b.CollegeId
		user.EntranceYear = b.EntranceYear
	}
	if err := dao.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ?", id).First(&user)

	return &user, nil
}

func (u *User) UpdatePassword(id uint, pwd req.UpdateUserPassword) error {
	var user dao.User
	if err := dao.DB.Preload(clause.Associations).Model(&dao.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	user.SetPassword(pwd.Password)
	dao.DB.Save(&user)
	return nil
}
