package models

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64     `orm:"column(id);auto;pk"`
	Username     string    `orm:"column(username);unique"`
	PasswordHash string    `orm:"column(password_hash)"`
	WeightUnit   string    `orm:"column(weight_unit)"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now"`
}

func (u *User) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(&User{})
}

func CreateUser(username, password, weightUnit string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     username,
		PasswordHash: string(hash),
		WeightUnit:   weightUnit,
	}

	o := orm.NewOrm()
	if _, err := o.Insert(user); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	o := orm.NewOrm()
	user := &User{Username: username}
	if err := o.Read(user, "Username"); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func GetUserByID(id int64) (*User, error) {
	o := orm.NewOrm()
	user := &User{ID: id}
	if err := o.Read(user); err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateWeightUnit(userID int64, unit string) error {
	if unit != "lb" && unit != "kg" {
		return errors.New("weight_unit must be lb or kg")
	}
	o := orm.NewOrm()
	_, err := o.QueryTable(&User{}).Filter("ID", userID).Update(orm.Params{"weight_unit": unit})
	return err
}

func DeleteUserByUsername(username string) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(&User{}).Filter("Username", username).Delete()
	return err
}

func DeleteUserByID(id int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable(&User{}).Filter("ID", id).Delete()
	return err
}
