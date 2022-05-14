package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
)

// QueryCollegesService 查询学院信息的筛选条件
type QueryCollegesService struct {
	Name string `json:"name" description:"学院名称"` // 学院名称，将会模糊搜索
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

// NewCollegeService 新学院信息
type NewCollegeService struct {
	Name string `json:"name" description:"学院名称"` // 学院名称
}

func (n NewCollegeService) NewCollege() (*dao.College, error) {
	college := &dao.College{Name: n.Name}
	if err := dao.DB.Create(&college).Error; err != nil {
		return nil, err
	}
	return college, nil
}
