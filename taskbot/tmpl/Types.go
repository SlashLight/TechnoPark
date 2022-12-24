package tmpl

type User struct {
	ID     string
	ChatID int64
}

type Task struct {
	Content  string
	Author   User
	Executor *User
}

type Pool struct {
	Tasks []Task
}
