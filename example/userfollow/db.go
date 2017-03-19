package userfollow

import (
	"os"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gomodel"
	"github.com/cosiner/gomodel/driver"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB = gomodel.NewDB()
)

func werckerEnv(key, defval string) string {
	// Host: WERCKER_MYSQL_HOST
	// Port: WERCKER_MYSQL_PORT
	// Username: WERCKER_MYSQL_USERNAME
	// Password: WERCKER_MYSQL_PASSWORD
	// Database: WERCKER_MYSQL_DATABASE
	const MYSQL = "WERCKER_MYSQL_"
	val := os.Getenv(MYSQL + key)
	if val != "" {
		return val
	}

	return defval
}

func dsn(dri gomodel.Driver) string {
	host := werckerEnv("HOST", "localhost")
	port := werckerEnv("PORT", "3306")
	username := werckerEnv("USERNAME", "root")
	password := werckerEnv("PASSWORD", "root")
	database := werckerEnv("DATABASE", "test")

	return dri.DSN(host, port, username, password, database, map[string]string{
		"charset":         "utf8",
		"clientFoundRows": "true",
	})
}

func createTables() {
	dri := driver.MySQL("mysql")
	errors.Panic(
		DB.Connect(dri, dsn(dri), 1, 1),
	)

	_, err := DB.DB.Exec(`
CREATE TABLE user (
    id int AUTO_INCREMENT,
    name varchar(50) UNIQUE NOT NULL,
    age int NOT NULL DEFAULT 0,

    followings int NOT NULL DEFAULT 0,
    followers int NOT NULL DEFAULT 0,

    PRIMARY KEY(id)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8
    `)
	errors.Panic(err)

	_, err = DB.DB.Exec(`
CREATE TABLE user_follow (
    user_id varchar(16),
    follow_user_id varchar(16),

    PRIMARY KEY(user_id, follow_user_id)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8
    `)
	errors.Panic(err)
}

func dropTables() {
	_, err := DB.DB.Exec("DROP TABLE IF EXISTS user")
	errors.Panic(err)

	_, err = DB.DB.Exec("DROP TABLE IF EXISTS user_follow")
	errors.Panic(err)
}
