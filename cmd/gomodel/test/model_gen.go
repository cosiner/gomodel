package test

import (
	"github.com/cosiner/gomodel"
)

const (
	USER_ID uint64 = 1 << iota
	USER_NAME
	USER_AGE
	USER_FOLLOWINGS
	USER_FOLLOWERS

	UserFieldEnd             = iota
	UserFieldsAll            = 1<<UserFieldEnd - 1
	UserFieldsExcpId         = UserFieldsAll & (^USER_ID)
	UserFieldsExcpName       = UserFieldsAll & (^USER_NAME)
	UserFieldsExcpAge        = UserFieldsAll & (^USER_AGE)
	UserFieldsExcpFollowings = UserFieldsAll & (^USER_FOLLOWINGS)
	UserFieldsExcpFollowers  = UserFieldsAll & (^USER_FOLLOWERS)

	UserTable         = "user"
	UserIdCol         = "id"
	UserNameCol       = "name"
	UserAgeCol        = "age"
	UserFollowingsCol = "followings"
	UserFollowersCol  = "followers"
)

var (
	UserInstance = new(User)
)

func (uu *User) Table() string {
	return UserTable
}

func (uu *User) Columns() []string {
	return []string{
		UserIdCol, UserNameCol, UserAgeCol, UserFollowingsCol, UserFollowersCol,
	}
}

func (uu *User) Vals(fields uint64, vals []interface{}) {
	if fields != 0 {
		if fields == UserFieldsAll {
			vals[0] = uu.Id
			vals[1] = uu.Name
			vals[2] = uu.Age
			vals[3] = uu.Followings
			vals[4] = uu.Followers

		} else {
			index := 0
			if fields&USER_ID != 0 {
				vals[index] = uu.Id
				index++
			}
			if fields&USER_NAME != 0 {
				vals[index] = uu.Name
				index++
			}
			if fields&USER_AGE != 0 {
				vals[index] = uu.Age
				index++
			}
			if fields&USER_FOLLOWINGS != 0 {
				vals[index] = uu.Followings
				index++
			}
			if fields&USER_FOLLOWERS != 0 {
				vals[index] = uu.Followers
				index++
			}
		}
	}
}

func (uu *User) Ptrs(fields uint64, ptrs []interface{}) {
	if fields != 0 {
		if fields == UserFieldsAll {
			ptrs[0] = &(uu.Id)
			ptrs[1] = &(uu.Name)
			ptrs[2] = &(uu.Age)
			ptrs[3] = &(uu.Followings)
			ptrs[4] = &(uu.Followers)

		} else {
			index := 0
			if fields&USER_ID != 0 {
				ptrs[index] = &(uu.Id)
				index++
			}
			if fields&USER_NAME != 0 {
				ptrs[index] = &(uu.Name)
				index++
			}
			if fields&USER_AGE != 0 {
				ptrs[index] = &(uu.Age)
				index++
			}
			if fields&USER_FOLLOWINGS != 0 {
				ptrs[index] = &(uu.Followings)
				index++
			}
			if fields&USER_FOLLOWERS != 0 {
				ptrs[index] = &(uu.Followers)
				index++
			}
		}
	}
}

func (uu *User) TxDo(exec gomodel.Executor, do func(*gomodel.Tx, *User) error) error {
	var (
		tx  *gomodel.Tx
		err error
	)
	switch r := exec.(type) {
	case *gomodel.Tx:
		tx = r
	case *gomodel.DB:
		tx, err = r.Begin()
		if err != nil {
			return err
		}
		defer tx.Close()
	default:
		panic("unexpected underlay type of gomodel.Executor")
	}

	err = do(tx, uu)
	tx.Success(err == nil)
	return err
}

type (
	UserStore struct {
		Values []User
		Fields uint64
	}
)

func (s *UserStore) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]User, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *UserStore) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *UserStore) Ptrs(index int, ptrs []interface{}) {
	s.Values[index].Ptrs(s.Fields, ptrs)
}

func (s *UserStore) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]User, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of UserStore")
}

func (a *UserStore) Clear() {
	if a.Values != nil {
		a.Values = a.Values[:0]
	}
}

const (
	FOLLOW_USERID uint64 = 1 << iota
	FOLLOW_FOLLOWUSERID

	FollowFieldEnd               = iota
	FollowFieldsAll              = 1<<FollowFieldEnd - 1
	FollowFieldsExcpUserId       = FollowFieldsAll & (^FOLLOW_USERID)
	FollowFieldsExcpFollowUserId = FollowFieldsAll & (^FOLLOW_FOLLOWUSERID)

	FollowTable           = "user_follow"
	FollowUserIdCol       = "user_id"
	FollowFollowUserIdCol = "follow_user_id"
)

var (
	FollowInstance = new(Follow)
)

func (ff *Follow) Table() string {
	return FollowTable
}

func (ff *Follow) Columns() []string {
	return []string{
		FollowUserIdCol, FollowFollowUserIdCol,
	}
}

func (ff *Follow) Vals(fields uint64, vals []interface{}) {
	if fields != 0 {
		if fields == FollowFieldsAll {
			vals[0] = ff.UserId
			vals[1] = ff.FollowUserId

		} else {
			index := 0
			if fields&FOLLOW_USERID != 0 {
				vals[index] = ff.UserId
				index++
			}
			if fields&FOLLOW_FOLLOWUSERID != 0 {
				vals[index] = ff.FollowUserId
				index++
			}
		}
	}
}

func (ff *Follow) Ptrs(fields uint64, ptrs []interface{}) {
	if fields != 0 {
		if fields == FollowFieldsAll {
			ptrs[0] = &(ff.UserId)
			ptrs[1] = &(ff.FollowUserId)

		} else {
			index := 0
			if fields&FOLLOW_USERID != 0 {
				ptrs[index] = &(ff.UserId)
				index++
			}
			if fields&FOLLOW_FOLLOWUSERID != 0 {
				ptrs[index] = &(ff.FollowUserId)
				index++
			}
		}
	}
}

func (ff *Follow) TxDo(exec gomodel.Executor, do func(*gomodel.Tx, *Follow) error) error {
	var (
		tx  *gomodel.Tx
		err error
	)
	switch r := exec.(type) {
	case *gomodel.Tx:
		tx = r
	case *gomodel.DB:
		tx, err = r.Begin()
		if err != nil {
			return err
		}
		defer tx.Close()
	default:
		panic("unexpected underlay type of gomodel.Executor")
	}

	err = do(tx, ff)
	tx.Success(err == nil)
	return err
}

type (
	FollowStore struct {
		Values []Follow
		Fields uint64
	}
)

func (s *FollowStore) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]Follow, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *FollowStore) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *FollowStore) Ptrs(index int, ptrs []interface{}) {
	s.Values[index].Ptrs(s.Fields, ptrs)
}

func (s *FollowStore) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]Follow, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of FollowStore")
}

func (a *FollowStore) Clear() {
	if a.Values != nil {
		a.Values = a.Values[:0]
	}
}

// Generated SQL
var (
	astConv = gomodel.NewSqlId(func(gomodel.Executor) string {
		return "insert into user_follow(user_id, follow_user_id) select ?, ? from DUAL where exists (select id from user where id = ?)"
	})
)
