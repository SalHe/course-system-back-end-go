package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
)

type QueryCollegesService struct {
	Name string `json:"name"`
}

func (q QueryCollegesService) Query() []dao.College {
	var colleges []dao.College
	if len(q.Name) > 0 {
		dao.DB.Where("name LIKE ?", "%"+q.Name+"%").Find(&colleges)
	} else {
		dao.DB.Find(&colleges)
	}
	return colleges
}

type NewCollegeService struct {
	Name string `json:"name"`
}

func (n NewCollegeService) NewCollege() (*dao.College, error) {
	college := &dao.College{Name: n.Name}
	if err := dao.DB.Create(&college).Error; err != nil {
		return nil, err
	}
	return college, nil
}
