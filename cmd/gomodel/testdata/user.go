package testdata

type User struct {
	Name string `table:"us" column:"n"`
	Id   int    `column:"i"`
	Age  int
}
