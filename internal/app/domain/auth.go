package domain

type NewTokenInput struct {
	Username string
	Password string
}

type Token struct {
	ID        string
	Expiry    int64
	IssueTime int64
	Username  string
}
