package model

type NewTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	ID        string `json:"id"`
	Expiry    int64  `json:"expiry"`
	IssueTime int64  `json:"issue_time"`
	Username  string `json:"username"`
}
