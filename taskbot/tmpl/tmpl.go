package tmpl

import (
	"text/template"
)

func inc(a int) int {
	return a + 1
}

func deref(a *User) User {
	if a == nil {
		return User{ID: ""}
	}
	return *a
}

const (
	SHOWTASKS = "{{range $index, $task:=.}} {{$id:=inc $index}}{{$Executor:=deref $task.Executor}} {{$id}}. {{$task.Content}} by {{$task.Author.ID}}\n{{if $task.NotBeingExecuted}} /assign_{{$id}}\n{{else}}{{if eq $Executor.ID $task.Author.ID}} assigner: я\n /unassign_{{$id}}, /resolve_{{$id}} \n\n{{else}}assigner: {{$Executor.ID}}\n\n{{end}}{{end}}{{end}}"
)

var (
	TempShow, _ = template.New("Showing").Funcs(template.FuncMap{"inc": inc, "deref": deref}).Parse(SHOWTASKS)
)
