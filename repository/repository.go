package repository

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// SetDB データベースとの接続情報をグローバル変数に渡す
func SetDB(d *sqlx.DB) {
	db = d
}