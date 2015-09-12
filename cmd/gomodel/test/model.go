package test

//go:generate gomodel -model $GOFILE
//go:generate gomodel -sql -t sql.tmpl $GOFILE astconv.go
//go:generate gofmt -w model_gen.go

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
