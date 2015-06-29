package userfollow

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestMain(m *testing.M) {
	createTables()
	m.Run()
	dropTables()
}

func TestUser(t *testing.T) {
	tt := testing2.Wrap(t)

	u1 := User{
		Name: "User1",
		Age:  20,
	}
	tt.Nil(u1.Add()) // add user1

	u2 := u1 // duplicate user name
	tt.Eq(ErrDuplicateUserName, u2.Add())
	u2.Name = "User2"
	tt.Nil(u2.Add()) // add user2

	f1 := Follow{
		UserId:       u1.Id,
		FollowUserId: 1024000, // user id doesn't exists
	}
	tt.Eq(ErrNoUser, f1.Add())
	f1.FollowUserId = u2.Id
	tt.Nil(f1.Add()) // user1 follow user2

	f2 := Follow{
		UserId:       u2.Id,
		FollowUserId: u1.Id,
	}

	tt.
		Nil(f2.Add()). // user2 follow user1
		Eq(ErrFollowed, f2.Add())

	tt.
		Nil(u1.ById()). // query latest user1 info
		Eq(1, u1.Followers).
		Eq(1, u1.Followings).
		Nil(u2.ById()). // query latest user2 info
		Eq(1, u2.Followers).
		Eq(1, u2.Followings)

	tt.
		Nil(f1.Delete()). // delete follow relationships
		Nil(f2.Delete())

	tt.
		Nil(u1.ById()). // query latest user1 info
		Eq(0, u1.Followings).
		Eq(0, u1.Followings).
		Nil(u2.ById()). // query latest user2 info
		Eq(0, u2.Followings).
		Eq(0, u2.Followings)

	testing2.
		Expect(nil).Arg(u1.Id). // delete user1, user2
		Expect(nil).Arg(u2.Id).
		Run(t, DeleteUserById)
}
