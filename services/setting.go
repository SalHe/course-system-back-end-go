package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
)

const (
	KeyCurrentSemester = "current_semester"
	KeyEnableRegister  = "enable_register"
)

type Setting interface {
	Get(key string) (string, bool)
	Set(key string, value string) bool
}

type settingInDB struct{}

func (s *settingInDB) Get(key string) (string, bool) {
	var setting dao.Setting
	if err := dao.DB.Model(&dao.Setting{}).Where("key = ?", key).First(&setting).Error; err != nil {
		return "", false
	}
	return setting.Value, true
}

func (s *settingInDB) Set(key string, value string) bool {
	var setting dao.Setting
	var setError error
	if err := dao.DB.Model(&dao.Setting{}).Where("key = ?", key).First(&setting).Error; err != nil {
		setError = dao.DB.Create(&dao.Setting{Key: key, Value: value}).Error
	} else {
		setError = dao.DB.Model(&dao.Setting{}).Where("key = ?", key).Update("value", value).Error
	}
	return setError == nil
}
