// Package models model
package models

import (
    "database/sql"

    "github.com/GoAdminGroup/go-admin/modules/db"
)

// Base is base model structure.
type Base struct {
    TableName string

    Conn db.Connection
    Tx   *sql.Tx
}

// SetConn 设置数据库连接
func (b Base) SetConn(con db.Connection) Base {
    b.Conn = con
    return b
}

// Table 设置表名
func (b Base) Table(table string) *db.SQL {
    return db.Table(table).WithDriver(b.Conn)
}
