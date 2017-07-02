package test

//go:generate gomodel $GOFILE

type User struct {
	Id   int64
	Name string
	Age  int

	Followings int
	Followers  int
}

type Follow struct {
	UserId       int64 `table:"user_follow"`
	FollowUserId int64
}
