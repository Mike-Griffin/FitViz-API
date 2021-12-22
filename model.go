package main

import "database/sql"

type Activity struct {
	ActivityID int           `json:"id"`
	Type       string        `json:"type"`
	UserID     sql.NullInt16 `json:"user_id"`
}

type User struct {
	UserID      int    `json:"user_id"`
	CreatedTime string `json:"createtime"`
}
