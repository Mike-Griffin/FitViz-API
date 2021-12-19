package main

type Activity struct {
	ActivityID int    `json:"id"`
	Type       string `json:"type"`
}

type User struct {
	UserID      int    `json:"id"`
	CreatedTime string `json:"createtime"`
}
