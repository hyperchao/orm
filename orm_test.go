package orm

import (
	"context"
	"fmt"
	"reflect"
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

func Test_X(t *testing.T) {
	db, err := Open("sqlite3", "./foo.db")
	checkErr(err)

	rows, err := db.QueryContext(context.Background(), "select * from userinfo where uid = ?", 1)
	checkErr(err)

	var uid interface{}
	var username interface{}
	var department interface{}
	var createAt interface{}

	if rows.Next() {
		err = rows.Scan(uid, username, department, createAt)
	}
	checkErr(err)
}

func Test_A(t *testing.T) {
	// db, err := sql.Open("sqlite3", "./foo.db")
	// checkErr(err)
	// _, err = db.Exec(`
	// CREATE TABLE userinfo (
	// 	uid INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	username VARCHAR(64) NULL,
	// 	department VARCHAR(64) NULL,
	// 	created DATE NULL
	// )`)
	// checkErr(err)

	// res, err := db.Exec(`
	// INSERT INTO userinfo(username, department, created) VALUES(?,?,?)
	// `, "astaxie", "研发部门", nil)
	// checkErr(err)

	// id, err := res.LastInsertId()
	// checkErr(err)
	// t.Log(id)

	// rows, err := db.Query("select * from userinfo where uid = 1")
	// checkErr(err)
	// cols, err := rows.Columns()
	// checkErr(err)
	// t.Log(cols)

	// for rows.Next() {
	// var uid int
	// var username, department string
	// var created *time.Time
	// v := reflect.ValueOf(uid).Addr().Interface()
	// err = rows.Scan(v, &username, &department, &created)
	// checkErr(err)
	// t.Log(uid, username, department, *created)
	// }

	db, err := Open("sqlite3", "./foo.db")
	checkErr(err)

	u := UserInfo{}
	ok, err := db.GetOne(context.Background(), &u, "select * from userinfo where uid = ?", 1)
	checkErr(err)
	fmt.Printf("%v %+v\n", ok, u)

}

type Inner struct {
	A int `orm:"aa"`
	B int
}

type Embed struct {
	C int `orm:"cc"`
	D int
}

type S struct {
	Embed
	I *Inner `orm:"age"`
	T struct {
		X int
		Y int `orm:"user"`
	} `orm:"user"`
}

func Test_B(t *testing.T) {
	x := reflect.TypeOf(S{})
	r := getColumnIndexMapping(x)
	fmt.Printf("%+v\n", r)
}

func traverse(prefix string, t reflect.Type) {
	fmt.Printf("%stype = %v\n", prefix, t)
	if t.Kind() != reflect.Struct {
		return
	}
	prefix += "  "
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("%sname = %s, index = %v, pkgPath = %v, anonymous = %v, offset = %d, type = %v \n", prefix, f.Name, f.Index, f.PkgPath, f.Anonymous, f.Offset, f.Type)
		traverse(prefix, f.Type)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Test_C(t *testing.T) {
	a := make([]int, 2)
	a[0] = 5
	v := reflect.ValueOf(a).Index(0).Addr()
	fmt.Printf("%T %v", v.Interface(), v.Interface())
}
