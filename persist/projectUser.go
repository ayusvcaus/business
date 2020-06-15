package persist

import (
	"github.com/alecthomas/log4go"
)

type Project_users struct {
	ProjectID int `json:"id" gorm:"primary_key`
	UserID    int `json:"id" gorm:"primary_key`
}

func (project_users *Project_users) AssociateProject2User() (string, error) {
	tx := Db.Begin()
	if err := tx.Save(&project_users).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return "mapping succeeded!", nil
}

func (project_users *Project_users) Delete() (string, error) {
	tx := Db.Begin()
	err := tx.Where("project_id = ? and user_id = ?", project_users.ProjectID, project_users.UserID).Delete(&project_users).Error
	if err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return "deleted", nil
}
