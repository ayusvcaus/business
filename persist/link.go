package persist

import (
	"github.com/alecthomas/log4go"
)

type Link struct {
	ID      int    `json:"id" gorm:"primary_key`
	Title   string `json:"title"`
	Address string `json:"address"`
	UserID  int    `json:"-"`
}

func (link *Link) Create() (*Link, error) {
	tx := Db.Begin()
	if err := tx.Create(&link).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return link, nil
}

func (link *Link) Update() (*Link, error) {
	tx := Db.Begin()
	if err := tx.Save(&link).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return link, nil
}

func (link *Link) Delete() error {
	tx := Db.Begin()
	if err := tx.Where("id = ?", link.ID).Delete(&link).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (link *Link) Links() ([]*Link, error) {
	var links []*Link
	err := Db.Find(&links).Error
	if err != nil {
		log4go.Info(err)
		return nil, err
	}
	return links, nil
}

func (link *Link) GetLinksByTitle(title string) ([]*Link, error) {
	var links []*Link
	err := Db.Where("title LIKE ?", "%"+title+"%").Find(&links).Error
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	return links, nil
}
