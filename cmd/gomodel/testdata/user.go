package testdata

type UserInfo struct {
	Name string `table:"userInfo"`
	Id   int    `column:"user_id"`
	Age  int
}
