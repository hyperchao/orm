package orm

import (
	"database/sql"
)

type Model struct {
	ID       int64        `orm:"id"`
	CreateAt sql.NullTime `orm:"create_at"`
	UpdateAt sql.NullTime `orm:"update_at"`
	DeleteAt sql.NullTime `orm:"delete_at"`
}
