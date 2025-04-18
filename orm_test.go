package orm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UserInfo struct {
	Uid        int64      `orm:"uid"`
	Username   string     `orm:"username"`
	Department string     `orm:"department"`
	CreateAt   *time.Time `orm:"created"`
}

type CustomTagUserInfo struct {
	Uid      int64      `foobar:"uid"`
	Username string     `foobar:"username"`
	CreateAt *time.Time `foobar:"created"`
}

type BaseEntity struct {
	Uid int64 `orm:"uid"`
}

type TimedEntity struct {
	BaseEntity
	CreateAt *time.Time `orm:"created"`
}

type NestedUserInfo struct {
	*TimedEntity
	Username string `orm:"username"`
}

func Test_Init(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)

	var name string
	r, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='userinfo'")
	assert.Nil(t, err)
	if r.Next() {
		err = r.Scan(&name)
		assert.Nil(t, err)
		assert.Equal(t, name, "userinfo")
		return
	}

	_, err = db.Exec(`
	CREATE TABLE userinfo (
	 	uid INTEGER PRIMARY KEY AUTOINCREMENT,
	 	username VARCHAR(64) NULL,
	 	department VARCHAR(64) NULL,
	 	created DATE NULL
	 )`)
	assert.Nil(t, err)

	_, err = db.Exec(`
	 INSERT INTO userinfo(username, department, created) VALUES(?,?,?)
	 `, "astaxie", "研发部门", nil)
	assert.Nil(t, err)
}

func Test_RawUse(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)

	rows, err := db.QueryContext(context.Background(), "select * from userinfo where uid = ?", 1)
	assert.Nil(t, err)

	var uid interface{}
	var username interface{}
	var department interface{}
	var createAt interface{}

	if rows.Next() {
		err = rows.Scan(&uid, &username, &department, &createAt)
	}
	assert.Nil(t, err)
}

func Test_GetOne(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)

	u, err := GetOne[UserInfo](context.Background(), db, "select * from userinfo where uid = ?", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), u.Uid)
	fmt.Printf("%+v\n", u)
}

func Test_GetOne_Pointer(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)

	u, err := GetOne[*UserInfo](context.Background(), db, "select * from userinfo where uid = ?", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), (*u).Uid)
	fmt.Printf("%+v\n", u)
}

func Test_GetOne_CustomTag(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)
	SetTagName("foobar")
	u, err := GetOne[CustomTagUserInfo](context.Background(), db, "select * from userinfo where uid = ?", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), u.Uid)
	fmt.Printf("%+v\n", u)
	SetTagName("orm")
}

func Test_GetOne_CustomTag_WithOpt(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)
	u, err := GetOne[CustomTagUserInfo](context.Background(), db, "select * from userinfo where uid =?", 1, WithTagName("foobar"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), u.Uid)
	fmt.Printf("%+v\n", u)
}

func Test_GetOne_NestedStruct(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)
	u, err := GetOne[NestedUserInfo](context.Background(), db, "select * from userinfo where uid =?", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), u.Uid)
	fmt.Printf("%+v\n", u)
}
