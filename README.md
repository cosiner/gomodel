# gomodel [![wercker status](https://app.wercker.com/status/9c6ef0eec7d6d217bd831bbdc3a3ace2/s "wercker status")](https://app.wercker.com/project/bykey/9c6ef0eec7d6d217bd831bbdc3a3ace2) [![GoDoc](https://godoc.org/github.com/cosiner/gomodel?status.png)](http://godoc.org/github.com/cosiner/gomodel)
gomodel provide another method to interact with database.   
Instead of reflection, use bitset represent fields of CRUD, sql/stmt cache and generate model code for you, high performance.

# Install
```sh
$ go get github.com/cosiner/gomodel
$ cd /path/of/gomodel/cmd/gomodel
$ go install # it will install the gomodel binary file to your $GOPATH/bin
$ gomodel -cp # copy model.tmpl to default path $HOME/.config/go/model.tmpl
              # or just put it to your model package, gomodel will search it first 
$ # if need customed template, just copy model.tmpl to your models directory
```

The [gomodel cmd tool and SQL convertion for structures](https://github.com/cosiner/gomodel/tree/master/cmd/gomodel).

There is a detailed example [User-Follow](https://github.com/cosiner/gomodel/tree/master/example/userfollow).
#### User
```Go

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
```
### Follow
```Go
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
    return f.txDo(func(tx gomodel.Tx, f *Follow) error {
        stmt, err := tx.PrepareById(insertUserFollowSQL)
        c, err := gomodel.CloseUpdate(stmt, err, gomodel.FieldVals(f, followFieldsAll, f.FollowUserId)...)

        err = dberrs.NoAffects(c, err, ErrNoUser)
        err = dberrs.DuplicateKeyError(err, dberrs.PRIMARY_KEY, ErrFollowed)

        return f.updateFollowInfo(tx, err, 1)
    })
}

func (f *Follow) Delete() error {
    return f.txDo(func(tx gomodel.Tx, f *Follow) error {
        c, err := tx.Delete(f, followFieldsAll)
        err = dberrs.NoAffects(c, err, ErrNonFollow)

        return f.updateFollowInfo(tx, err, -1)
    })
}

func (f *Follow) updateFollowInfo(tx gomodel.Tx, err error, c int) error {
    if err == nil {
        _, err = tx.ArgsIncrBy(userInstance, USER_FOLLOWINGS, USER_ID, c, f.UserId)
        if err == nil {
            _, err = tx.ArgsIncrBy(userInstance, USER_FOLLOWERS, USER_ID, c, f.FollowUserId)
        }
    }
    return err
}

```

# LICENSE
MIT.
