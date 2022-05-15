package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
)

type College struct{}

func (c College) Query(q req.QueryCollegesService) []dao.College {
	var colleges []dao.College
	if len(q.Name) > 0 {
		dao.DB.Where("name LIKE ?", "%"+q.Name+"%").Find(&colleges)
	} else {
		dao.DB.Find(&colleges)
	}
	return colleges
}

func (c College) NewCollege(n req.NewCollegeService) (*dao.College, error) {
	college := &dao.College{Name: n.Name}
	if err := dao.DB.Create(&college).Error; err != nil {
		return nil, err
	}
	return college, nil
}
