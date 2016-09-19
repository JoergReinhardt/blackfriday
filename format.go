package agiledoc

import (
	t "text/template"
)

var Tmpl map[string]*t.Template = map[string]*t.Template{}

func init() {
	Tmpl["element"], _ = t.New("element").Parse(
		"{{,String()}}",
	)
	Tmpl["tuple"], _ = t.New("tuple").Parse(
		"{{template \"element\" ,Key()}}{{if ,Type()&NUMERIC !=0}}.) {{else if .Type()&SYMBOLIC != 0}}: {{else if .Type()&REAL !=0}}/{{end}}{{template \"element\" .Key()}}",
	)

	// range over all elements and print value
	Tmpl["list"], _ = t.New("list").Parse(
		"{{range $ok :=.Next()}}{{if ok}}{{template \"element\" .Value()}}\n{{{end}}{end}}",
	)

	// expect iterator over flat list. represent index preceeding the value
	Tmpl["idxList"], _ = t.New("idxList").Parse(
		// tuple dispatches the appropriate delimiter type
		"{{range $i,$e := .}}{{template \"tuple\" Value($i,$e).(pair)}}{{end}}",
	)

	// expect an iterator that yields tuples. range over all elements and
	// print Typle
	Tmpl["map"], _ = t.New("map").Parse(
		"{{for $ok := .Next()}}{{if $ok}}{{template \"tuple\" .Value(.Key(), .Value()),(pair)}}\n{{end}}{{end}}",
	)
}
