package main

type Activity struct {
	ActivityID int    `json:"id"`
	Type       string `json:"type"`
	UserID     int    `json:"user_id"`
}

type User struct {
	UserID      int    `json:"user_id"`
	CreatedTime string `json:"createtime"`
}
