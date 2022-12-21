package tmpl

type User struct {
	ID string
}

type Task struct {
	Content  string
	Author   User
	Executor *User
}

type Pull struct {
	Tasks []Task
}
