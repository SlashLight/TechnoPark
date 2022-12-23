package tmpl

func inc(a int) int {
	return a + 1
}

func deref(a *User) User {
	if a == nil {
		return User{}
	}
	return *a
}

func isMe(sender string) func(user string) bool {
	return func(user string) bool {
		return sender == user
	}
}

func isActive(task Task) bool {
	if task.Executor == nil {
		return true
	}

	return false
}
