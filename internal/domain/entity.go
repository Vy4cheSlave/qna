package domain

type Question struct {
	Id   int
	Text string
}

type Answer struct {
	Id         int
	QuestionId int
	UserId     string
	Text       string
}

type User struct {
	Id   string
	Name string
}
