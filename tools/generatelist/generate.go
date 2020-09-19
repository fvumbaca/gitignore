package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const baseURL = "https://www.toptal.com/developers/gitignore/api/list"
const modTemplate = `
package main

var autocompleteList = []string{
	{{- range .}}
	"{{quote .}}",
	{{- end }}
}
`

var funcmap = map[string]interface{}{
	"quote": quoteItem,
}

func quoteItem(item string) string {
	return strings.TrimSpace(strings.ReplaceAll(item, "\n", ""))
}

func main() {
	resp, err := http.DefaultClient.Get(baseURL)
	fatalOnErr(err)

	list, err := ioutil.ReadAll(resp.Body)
	fatalOnErr(err)

	tpl, err := template.New("template").Funcs(funcmap).Parse(modTemplate)
	fatalOnErr(err)

	f, err := os.Create("autocomplete.go")
	fatalOnErr(err)
	defer f.Close()

	err = tpl.Execute(f, strings.Split(string(list), ","))
	fatalOnErr(err)
}

func fatalOnErr(err error) {
	if err != nil {
		printFatal(err)
	}
}

func printFatal(a ...interface{}) {
	fmt.Fprint(os.Stderr, "error: ")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
