package tmpl

func inc(a int) int {
	return a + 1
}

func deref(a *User) User {
	if a == nil {
		return User{ID: ""}
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
		return false
	}

	return true
}

const (
	SHOWTASKS = "{{range $index, $task:=.}}\n    {{$id:=inc $index}}{{$Executor:=deref $task.Executor}} {{$id}}. {{$task.Content}} by {{$task.Author.ID}}\\n\n    {{if isActive $task}}\n        /assign_{{$id}}\\n\\n\n    {{else}}{{if isMe $Executor.ID}}\n        assigner: я\\n\n        /unassign_{{$id}}, /resolve_{{$id}}\\n\\n\n    {{else}}\n        assigner: {{$Executor.ID}}\\n\\n\n\n    {{end}}{{end}}\n{{end}}"
)

var (
//TempShow = template.Must(template.New("Templates").Funcs(template.FuncMap{"inc": inc, "deref": deref, "isMe": isMe()}).ParseFiles("./tmpl/templates.txt"))
)
