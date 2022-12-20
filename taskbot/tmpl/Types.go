package tmpl

type User struct {
	ID string
}

type Task struct {
	Content          string
	Author           User
	Executor         *User
	NotBeingExecuted bool
}

type Pull struct {
	Tasks []Task
}
