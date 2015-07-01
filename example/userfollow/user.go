package userfollow

import (
	"github.com/cosiner/gohper/utils/httperrs"
	"github.com/cosiner/gomodel"
	"github.com/cosiner/gomodel/dberrs"
)

var (
	ErrDuplicateUserName = httperrs.Conflict.NewS("user name already exists")
	ErrNoUser            = httperrs.NotFound.NewS("user not found")
)

type User struct {
	Id   int64
	Name string
	Age  int

	Followings int
	Followers  int
}

func (u *User) Add() error {
	u.Followings = 0
	u.Followers = 0

	id, err := DB.Insert(u, userFieldsExcpId, gomodel.RES_ID)
	err = dberrs.DuplicateKeyError(err, UserNameCol, ErrDuplicateUserName)
	if err == nil {
		u.Id = id
	}

	return err
}

func DeleteUserById(id int64) error {
	c, err := DB.ArgsDelete(userInstance, USER_ID, id)

	return dberrs.NoAffects(c, err, ErrNoUser)
}

func (u *User) Update() error {
	c, err := DB.Update(u, USER_NAME|USER_AGE, USER_ID)

	return dberrs.NoAffects(c, err, ErrNoUser)
}

func (u *User) ById() error {
	err := DB.One(u, userFieldsExcpId, USER_ID)

	return dberrs.NoRows(err, ErrNoUser)
}

func UsersByAge(age, start, count int) ([]User, error) {
	users := userStore{
		Fields: userFieldsAll,
	}

	err := DB.ArgsLimit(&users, userInstance, userFieldsAll, USER_AGE, age, start, count)
	return users.Values, dberrs.NoRows(err, ErrNoUser)
}

func AllUsersByAge(age int) ([]User, error) {
	users := userStore{
		Fields: userFieldsAll,
	}

	err := DB.ArgsAll(&users, userInstance, userFieldsAll, USER_AGE, age)
	return users.Values, dberrs.NoRows(err, ErrNoUser)
}
