package orm

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UserInfo struct {
	Uid        int64      `orm:"uid,primary,autoincrement"`
	Username   string     `orm:"username"`
	Department string     `orm:"department"`
	CreateAt   *time.Time `orm:"created"`

	Leader *UserInfo // test circular reference
}

type CustomTagUserInfo struct {
	Uid      int64      `foobar:"uid"`
	Username string     `foobar:"username"`
	CreateAt *time.Time `foobar:"created"`
}

func initDb(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.Nil(t, err)

	_, err = db.Exec(`
	CREATE TABLE userinfo (
	 	uid INTEGER PRIMARY KEY AUTOINCREMENT,
	 	username VARCHAR(64) NULL,
	 	department VARCHAR(64) NULL,
	 	created DATE NULL
	 )`)

	assert.Nil(t, err)

	return db
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

	_, err = db.Exec(`
	 INSERT INTO userinfo(username, department, created) VALUES(?,?,?)
	 `, "qqqq", "研发部门", nil)
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
}

func BenchmarkGetOne(b *testing.B) {
	db, _ := sql.Open("sqlite3", "./foo.db")
	for i := 0; i < b.N; i++ {
		u, _ := GetOne[UserInfo](context.Background(), db, "select * from userinfo where uid in ?", []int64{1})
		_ = u
	}
}

func Test_GetOne_Pointer(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)

	u, err := GetOne[*UserInfo](context.Background(), db, "select * from userinfo where uid = ?", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), (*u).Uid)
}

func Test_GetOne_CustomTag(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)
	SetTagName("foobar")
	u, err := GetOne[CustomTagUserInfo](context.Background(), db, "select * from userinfo where uid = ?", 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), u.Uid)
	SetTagName("orm")
}

func Test_GetOne_CustomTag_WithOpt(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)
	u, err := GetOne[CustomTagUserInfo](context.Background(), db, "select * from userinfo where uid =?", 1, WithTagName("foobar"))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), u.Uid)
}

func Test_GetOne_Non_Struct(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)
	result, err := GetOne[int](context.Background(), db, "select * from userinfo where uid =?", 1)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, *result)
}

func Test_GetMany(t *testing.T) {
	db, err := sql.Open("sqlite3", "./foo.db")
	assert.Nil(t, err)

	users, err := GetMany[UserInfo](context.Background(), db, "select * from userinfo where uid in ?", []int{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))

	users, err = GetMany[UserInfo](context.Background(), db, "select * from userinfo where uid not in ?", ([]int)(nil))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(users))

	users, err = GetMany[UserInfo](context.Background(), db, "select * from userinfo where uid in ? and department = ?", []int{1, 2}, "研发部门")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))

}

func Test_InsertOne(t *testing.T) {
	db := initDb(t)
	userinfo := UserInfo{
		Username:   "astaxie",
		Department: "<UNK>",
		CreateAt:   nil,
	}
	err := InsertOne(context.Background(), "userinfo", db, userinfo)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), userinfo.Uid)
	assert.Equal(t, userinfo.Username, "astaxie")
	assert.Equal(t, userinfo.Department, "<UNK>")
	assert.Nil(t, userinfo.CreateAt)

	err = InsertOne(context.Background(), "userinfo", db, &userinfo)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), userinfo.Uid)

	userinfo2, err := GetOne[UserInfo](context.Background(), db, "select * from userinfo where uid = ?", 2)
	assert.Nil(t, err)
	assert.Equal(t, userinfo.Uid, userinfo2.Uid)
	assert.Equal(t, userinfo.Username, userinfo2.Username)
	assert.Equal(t, userinfo.Department, userinfo2.Department)
	assert.Equal(t, userinfo.CreateAt, userinfo2.CreateAt)
}

func Test_InsertMany(t *testing.T) {
	db := initDb(t)
	userinfo := &UserInfo{
		Username:   "astaxie",
		Department: "<UNK>",
		CreateAt:   nil,
	}
	userinfos := []*UserInfo{userinfo, userinfo, userinfo}
	err := InsertMany(context.Background(), "userinfo", db, userinfos, WithBatchSize(2))
	assert.Nil(t, err)

	userinfos, err = GetMany[UserInfo](context.Background(), db, "select * from userinfo")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(userinfos))
}
