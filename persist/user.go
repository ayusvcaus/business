package persist

import (
	"strconv"

	"github.com/alecthomas/log4go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int       `json:"id" gorm:"primary_key"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Links    []*Link   `json:"links" gorm:"foreignkey:UserID"`
	Projects []Project `json:"projects" gorm:"many2many:project_users;"`
}

func (user *User) GetUserById(id string) (*User, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	var dbuser User
	err = Db.Where(&User{ID: i}).First(&dbuser).Error
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	return &dbuser, nil
}

func (user *User) GetUserByUsername(username string) (*User, error) {
	var dbuser User
	err := Db.Preload("Links").Preload("Projects").Where(&User{Username: username}).First(&dbuser).Error
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	return &dbuser, nil
}

func (user *User) Authenticate() bool {
	pwd := user.Password
	err := Db.Where("username = ?", user.Username).Find(&user).Error
	if err != nil {
		log4go.Error(err.Error())
		return false
	}

	if user.Username == "admin" {
		if pwd == user.Password {
			return true
		}
		return false
	}
	if CheckPasswordHash(pwd, user.Password) {
		return true
	}
	return false
}

func (user *User) Create() (*User, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	user.Password = hashedPassword
	tx := Db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return user, nil
}

func (user *User) UpdatePassword() error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		log4go.Error(err.Error())
		return err
	}
	err = Db.Where("username = ?", user.Username).Find(&user).Error
	user.Password = hashedPassword
	if err != nil {
		log4go.Error(err.Error())
		return err
	}
	tx := Db.Begin()
	if err = tx.Save(&user).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (user *User) DeleteUser() (string, error) {
	tx := Db.Begin()
	if err := tx.Where("username = ?", user.Username).Delete(&user).Error; err != nil {
		log4go.Error(err.Error())
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	return "deleted", nil
}

func (user *User) IsExist() bool {
	if !Db.NewRecord(user) {
		return true
	}
	return false
}

//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (user *User) AllUsers() ([]*User, error) {
	var users []*User
	err := Db.Preload("Links").Preload("Projects").Find(&users).Error
	if err != nil {
		log4go.Info(err)
		return nil, err
	}
	return users, nil
}

type UserError struct {
	Err string
}

func (m *UserError) Error() string {
	return m.Err
}
