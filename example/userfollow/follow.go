package userfollow

import (
	"github.com/cosiner/gohper/utils/httperrs"
	"github.com/cosiner/gomodel"
	"github.com/cosiner/gomodel/dberrs"
)

var (
	ErrFollowed  = httperrs.Conflict.NewS("already followed")
	ErrNonFollow = httperrs.NotFound.NewS("non follow")
)

type Follow struct {
	UserId       int64 `table:"user_follow"`
	FollowUserId int64
}

//gomodel insertUserFollowSQL = [
//  INSERT INTO Follow(UserId, FollowUserId)
//      SELECT ?, ? FROM DUAL
//      WHERE EXISTS(SELECT Id FROM User WHERE Id=?)
//]
func (f *Follow) Add() error {
	return f.txDo(DB, func(tx *gomodel.Tx, f *Follow) error {
		c, err := tx.UpdateById(insertUserFollowSQL, gomodel.FieldVals(f, followFieldsAll, f.FollowUserId)...)

		err = dberrs.NoAffects(c, err, ErrNoUser)
		err = dberrs.DuplicatePrimaryKeyError(tx, err, ErrFollowed)

		return f.updateFollowInfo(tx, err, 1)
	})
}

func (f *Follow) Delete() error {
	return f.txDo(DB, func(tx *gomodel.Tx, f *Follow) error {
		c, err := tx.Delete(f, followFieldsAll)
		err = dberrs.NoAffects(c, err, ErrNonFollow)

		return f.updateFollowInfo(tx, err, -1)
	})
}

func (f *Follow) updateFollowInfo(tx *gomodel.Tx, err error, c int) error {
	if err == nil {
		_, err = tx.ArgsIncrBy(UserInstance, USER_FOLLOWINGS, USER_ID, c, f.UserId)
		if err == nil {
			_, err = tx.ArgsIncrBy(UserInstance, USER_FOLLOWERS, USER_ID, c, f.FollowUserId)
		}
	}
	return err
}
