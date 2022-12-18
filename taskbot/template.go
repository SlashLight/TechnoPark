package main

import "text/template"

func inc(a int) int {
	return a + 1
}

const (
	SHOWTASKS = `{{range $index, &task:=.}} {{$id:=inc $index}} {{$id}}. {{$item.Content}} by {{$item.Author}}\n {{if $item.Executor == nil}} /assign_{{$id}} {{if else $item.Executor == $item.Author}} assigner: я\n /unassign_$id, /resolve_$id \n \n {{else}} assigner: $item.Author.ID {{end}} {{end}}`
)

var (
	tempShow, _ = template.New("Showing").Funcs(template.FuncMap{"inc": inc}).Parse(SHOWTASKS)
)
