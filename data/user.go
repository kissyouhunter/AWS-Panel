package data

import (
	"errors"
	"github.com/Yuzuki616/Aws-Panel/utils"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type UserData struct {
	gorm.Model
	Status   int    `gorm:"Status"` //0 正常,1 封禁
	Username string `gorm:"Username"`
	Email    string `gorm:"Email"`
	Password string `gorm:"Password"`
	IsAdmin  int    `gorm:"IsAdmin"` //0 否,1 是
}

func VerifyUser(Username string, Password string) error {
	var user UserData
	Db.Where("Username = ?", Username).First(&user)
	if user.ID == 0 || !utils.VerifyPasswordHash(Password, user.Password) {
		return errors.New("用户名或密码错误")
	}
	if user.Status == 1 {
		return errors.New("用户已封禁")
	}
	return nil
}

func CreateUser(Username, Email, Password string, IsAdmin int) error {
	var user UserData
	Db.Where("Username = ?", Username).First(&user)
	if user.ID != 0 {
		return errors.New("用户已存在")
	}
	user.Username = Username
	user.Email = Email
	Password, err := utils.GenPasswordHash(Password)
	if err != nil {
		return err
	}
	user.Password = Password
	user.Status = 0
	user.IsAdmin = IsAdmin
	Db.Create(&user)
	return nil
}

func ChangeUsername(OldUsername, NewUsername, Password string) error {
	var user UserData
	Db.Where("Username = ? and Password = ?", OldUsername, Password)
	if user.ID == 0 {
		return errors.New("用户名或密码错误")
	}
	user.Username = NewUsername
	Db.Save(&user)
	return nil
}

func ChangeUserPassword(username string, OldPasswordMd5, NewPasswordMd5 string) error {
	var user UserData
	Db.Where("Username = ?", username).First(&user)
	if user.ID == 0 {
		return errors.New("用户不存在或密码错误")
	}
	if !utils.VerifyPasswordHash(OldPasswordMd5, user.Password) {
		return errors.New("用户不存在或密码错误")
	}
	user.Password = NewPasswordMd5
	Db.Save(&user)
	return nil
}

func IsAdmin(Username string) bool {
	var user UserData
	Db.Where("Username = ?", Username).First(&user)
	if user.ID == 0 {
		return false
	}
	if user.IsAdmin == 1 {
		return true
	}
	return false
}

func DeleteUser(Username string) error {
	var user UserData
	Db.Where("Username = ?", Username).First(&user)
	if user.ID == 0 {
		return errors.New("用户不存在")
	}
	if user.IsAdmin == 1 {
		return errors.New("不能删除管理员账户")
	}
	Db.Delete(&user)
	return nil
}

func BanUser(Username string) error {
	var user UserData
	Db.Where("Username = ?", Username).First(&user)
	if user.ID == 0 {
		return errors.New("用户不存在")
	}
	if user.IsAdmin == 1 {
		return errors.New("不能封禁管理员账户")
	}
	if user.Status == 0 {
		user.Status = 1
		Db.Save(&user)
	} else {
		return errors.New("该用户已被封禁")
	}
	return nil
}

func UnBanUser(Username string) error {
	var user UserData
	Db.Where("Username = ?", Username).First(&user)
	if user.ID == 0 {
		return errors.New("用户不存在")
	}
	if user.IsAdmin == 1 {
		return errors.New("不能解封管理员账户")
	}
	if user.Status == 1 {
		user.Status = 0
		Db.Save(&user)
	} else {
		return errors.New("该用户未被封禁")
	}
	return nil
}

type UserInfo struct {
	UserName string
	Status   int
	IsAdmin  int
}

func GetUserList() ([]UserInfo, error) {
	var u []UserData
	Db.Find(&u)
	if len(u) == 0 {
		return nil, errors.New("未找到任何用户")
	}
	var users []UserInfo
	for i := range u {
		users = append(users, UserInfo{
			UserName: u[i].Username,
			Status:   u[i].Status,
			IsAdmin:  u[i].IsAdmin,
		})
	}
	return users, nil
}
