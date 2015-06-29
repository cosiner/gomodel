package userfollow

import (
	"github.com/cosiner/gomodel"
)

const (
	FOLLOW_USERID uint64 = 1 << iota
	FOLLOW_FOLLOWUSERID

	followFieldEnd               = iota
	followFieldsAll              = 1<<followFieldEnd - 1
	followFieldsExcpUserId       = followFieldsAll & (^FOLLOW_USERID)
	followFieldsExcpFollowUserId = followFieldsAll & (^FOLLOW_FOLLOWUSERID)

	FollowUserIdCol       = "user_id"
	FollowFollowUserIdCol = "follow_user_id"
)

var (
	// DB is the global DB instance
	followInstance = &Follow{}
	FollowTable    = DB.Table(followInstance)
)

func (f *Follow) Table() string {
	return "user_follow"
}

func (f *Follow) Columns() []string {
	return []string{
		FollowUserIdCol, FollowFollowUserIdCol,
	}
}

func (f *Follow) Vals(fields uint64, vals []interface{}) {
	if fields != 0 {
		if fields == followFieldsAll {
			vals[0] = f.UserId
			vals[1] = f.FollowUserId

		} else {
			index := 0
			if fields&FOLLOW_USERID != 0 {
				vals[index] = f.UserId
				index++
			}
			if fields&FOLLOW_FOLLOWUSERID != 0 {
				vals[index] = f.FollowUserId
				index++
			}
		}
	}
}

func (f *Follow) Ptrs(fields uint64, ptrs []interface{}) {
	if fields != 0 {
		if fields == followFieldsAll {
			ptrs[0] = &(f.UserId)
			ptrs[1] = &(f.FollowUserId)

		} else {
			index := 0
			if fields&FOLLOW_USERID != 0 {
				ptrs[index] = &(f.UserId)
				index++
			}
			if fields&FOLLOW_FOLLOWUSERID != 0 {
				ptrs[index] = &(f.FollowUserId)
				index++
			}
		}
	}
}

func (f *Follow) txDo(do func(gomodel.Tx, *Follow) error) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return
	}

	defer tx.DeferDone(&err)

	err = do(tx, f)
	return
}

type (
	followStore struct {
		Values []Follow
		Fields uint64
	}
)

func (a *followStore) Make(count int) {
	if a.Values == nil {
		a.Values = make([]Follow, count)
	} else {
		a.Values = a.Values[:count]
	}
}

func (a *followStore) Ptrs(index int, ptrs []interface{}) bool {
	if len := len(a.Values); index == len {
		values := make([]Follow, 2*len)
		copy(values, a.Values)
		a.Values = values
	}

	a.Values[index].Ptrs(a.Fields, ptrs)
	return true
}

const (
	USER_ID uint64 = 1 << iota
	USER_NAME
	USER_AGE
	USER_FOLLOWINGS
	USER_FOLLOWERS

	userFieldEnd             = iota
	userFieldsAll            = 1<<userFieldEnd - 1
	userFieldsExcpId         = userFieldsAll & (^USER_ID)
	userFieldsExcpName       = userFieldsAll & (^USER_NAME)
	userFieldsExcpAge        = userFieldsAll & (^USER_AGE)
	userFieldsExcpFollowings = userFieldsAll & (^USER_FOLLOWINGS)
	userFieldsExcpFollowers  = userFieldsAll & (^USER_FOLLOWERS)

	UserIdCol         = "id"
	UserNameCol       = "name"
	UserAgeCol        = "age"
	UserFollowingsCol = "followings"
	UserFollowersCol  = "followers"
)

var (
	// DB is the global DB instance
	userInstance = &User{}
	UserTable    = DB.Table(userInstance)
)

func (u *User) Table() string {
	return "user"
}

func (u *User) Columns() []string {
	return []string{
		UserIdCol, UserNameCol, UserAgeCol, UserFollowingsCol, UserFollowersCol,
	}
}

func (u *User) Vals(fields uint64, vals []interface{}) {
	if fields != 0 {
		if fields == userFieldsAll {
			vals[0] = u.Id
			vals[1] = u.Name
			vals[2] = u.Age
			vals[3] = u.Followings
			vals[4] = u.Followers

		} else {
			index := 0
			if fields&USER_ID != 0 {
				vals[index] = u.Id
				index++
			}
			if fields&USER_NAME != 0 {
				vals[index] = u.Name
				index++
			}
			if fields&USER_AGE != 0 {
				vals[index] = u.Age
				index++
			}
			if fields&USER_FOLLOWINGS != 0 {
				vals[index] = u.Followings
				index++
			}
			if fields&USER_FOLLOWERS != 0 {
				vals[index] = u.Followers
				index++
			}
		}
	}
}

func (u *User) Ptrs(fields uint64, ptrs []interface{}) {
	if fields != 0 {
		if fields == userFieldsAll {
			ptrs[0] = &(u.Id)
			ptrs[1] = &(u.Name)
			ptrs[2] = &(u.Age)
			ptrs[3] = &(u.Followings)
			ptrs[4] = &(u.Followers)

		} else {
			index := 0
			if fields&USER_ID != 0 {
				ptrs[index] = &(u.Id)
				index++
			}
			if fields&USER_NAME != 0 {
				ptrs[index] = &(u.Name)
				index++
			}
			if fields&USER_AGE != 0 {
				ptrs[index] = &(u.Age)
				index++
			}
			if fields&USER_FOLLOWINGS != 0 {
				ptrs[index] = &(u.Followings)
				index++
			}
			if fields&USER_FOLLOWERS != 0 {
				ptrs[index] = &(u.Followers)
				index++
			}
		}
	}
}

func (u *User) txDo(do func(gomodel.Tx, *User) error) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return
	}

	defer tx.DeferDone(&err)

	err = do(tx, u)
	return
}

type (
	userStore struct {
		Values []User
		Fields uint64
	}
)

func (a *userStore) Make(count int) {
	if a.Values == nil {
		a.Values = make([]User, count)
	} else {
		a.Values = a.Values[:count]
	}
}

func (a *userStore) Ptrs(index int, ptrs []interface{}) bool {
	if len := len(a.Values); index == len {
		values := make([]User, 2*len)
		copy(values, a.Values)
		a.Values = values
	}

	a.Values[index].Ptrs(a.Fields, ptrs)
	return true
}

var (
	insertUserFollowSQL = gomodel.NewSqlId(func(gomodel.Tabler) string {
		return "insert into user_follow(user_id, follow_user_id) select ?, ? from DUAL where exists (select id from user where id = ?)"
	})
)
