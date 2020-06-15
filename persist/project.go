package persist

import (
	"github.com/alecthomas/log4go"
)

type Project struct {
	ID       int     `json:"id" gorm:"primary_key`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Buget    float64 `json:"buget"`
	Users    []*User `json:"users" gorm:"many2many:project_users;"`
}

func (project *Project) Create() (*Project, error) {
	tx := Db.Begin()
	if err := tx.Create(&project).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return project, nil
}

func (project *Project) Update() (*Project, error) {
	tx := Db.Begin()
	if err := tx.Save(&project).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return project, nil
}

func (project *Project) Delete() (string, error) {
	tx := Db.Begin()
	if err := tx.Where("id = ?", project.ID).Delete(&project).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return "deleted", nil
}

func (project *Project) Projects() ([]*Project, error) {
	var projects []*Project
	err := Db.Preload("Users").Find(&projects).Error
	if err != nil {
		log4go.Info(err)
		return nil, err
	}
	return projects, nil
}

func (project *Project) GetProjectsByName(name string) ([]*Project, error) {
	var projects []*Project
	err := Db.Where("name LIKE ?", "%"+name+"%").Preload("Users").Find(&projects).Error
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	return projects, nil
}
