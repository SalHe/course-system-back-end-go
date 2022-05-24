package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
)

type Start struct{}

func (s *Start) IsInitialized() bool {
	var count int64
	err := dao.DB.Model(&dao.User{}).Count(&count).Error
	return err == nil && count >= 1
}

func (s *Start) InitSystem(b req.InitRequest) error {
	if s.IsInitialized() {
		return ErrConflict
	} else {
		if b.Id == 0 || b.Username == "" || b.Password == "" {
			// 这里应该交给 req.BindAndValidate 处理
			return ErrInvalidParams
		}

		var adminCollegeName = "_____ADMIN"

		user := &dao.User{
			Model:    dao.Model{ID: b.Id},
			Username: b.Username,
			RealName: b.RealName,
			Role:     dao.RoleAdmin,
			College: &dao.College{
				Name: adminCollegeName,
			},
		}
		user.SetPassword(b.Password)
		result := dao.DB.Create(user)
		if err := result.Error; err != nil {
			return ErrUnknown
		}
		return nil
	}
}
