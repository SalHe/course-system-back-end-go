package services

import (
	"github.com/se2022-qiaqia/course-system/dao"
	"github.com/se2022-qiaqia/course-system/model/req"
	"strconv"
)

type Semester struct{}

func (s *Semester) GetSemesters() ([]*dao.Semester, error) {
	var semesters []*dao.Semester
	if err := dao.DB.Model(&dao.Semester{}).Find(&semesters).Error; err != nil {
		return nil, err
	}
	return semesters, nil
}

func (s *Semester) CreateSemester(semester *req.Semester) (*dao.Semester, error) {
	var count int64
	if err := dao.DB.Model(&dao.Semester{}).Where("year = ? and term = ?", semester.Year, semester.Term).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, ErrConflict
	}

	var dbSemester dao.Semester
	if err := dao.DB.Model(&dao.Semester{}).Create(&dao.Semester{
		Year: semester.Year,
		Term: semester.Term,
	}).Error; err != nil {
		return nil, err
	}

	dao.DB.Model(&dao.Semester{}).Where("year = ? and term = ?", semester.Year, semester.Term).First(&dbSemester)
	return &dbSemester, nil
}

func (s *Semester) GetCurrentSemester() (*dao.Semester, error) {
	ids, ok := Services.Setting.Get(KeyCurrentSemester)
	if !ok {
		return nil, ErrNotFound
	}
	id, _ := strconv.Atoi(ids)
	var semester dao.Semester
	if err := dao.DB.Model(&dao.Semester{}).Where("id = ?", id).First(&semester).Error; err != nil {
		return nil, err
	}
	return &semester, nil
}

func (s *Semester) SetCurrentSemester(id uint) (*dao.Semester, error) {
	var semester dao.Semester
	var count int64
	if err := dao.DB.Model(&dao.Semester{}).Where("id = ?", id).First(&semester).Count(&count).Error; err != nil || count <= 0 {
		return nil, ErrNotFound
	}

	if !Services.Setting.Set(KeyCurrentSemester, strconv.Itoa(int(id))) {
		return nil, ErrUpdateFailed
	}
	return &semester, nil
}
