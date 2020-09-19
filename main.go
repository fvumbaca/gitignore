package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

//go:generate go run tools/generatelist/generate.go

func main() {
	newRootCMD().ExecuteContext(context.Background())
}

const (
	baseURL = "https://www.toptal.com/developers/gitignore/api"
	minArgs = 1
)

func newRootCMD() *cobra.Command {
	var filename string
	var showBashCompletionInstead bool
	var showZshCompletionInstead bool
	cmd := cobra.Command{
		Use:   "gitignore [names...]",
		Short: "Generate .gitignore files",
		Long:  "Generate .gitignore files without leaving your shell. Powered by https://gitignore.io",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < minArgs && !(showBashCompletionInstead || showZshCompletionInstead) {
				return fmt.Errorf("requires at least %d arg(s), only received %d", minArgs, len(args))
			}
			return nil
		},
		ValidArgs: autocompleteList,
		Run: func(cmd *cobra.Command, args []string) {

			if showBashCompletionInstead {
				fatalOnErr(cmd.GenBashCompletion(os.Stdout))
				return
			}
			if showZshCompletionInstead {
				fatalOnErr(cmd.GenZshCompletion(os.Stdout))
				fmt.Fprintln(os.Stdout, "compdef _gitignore gitignore")
				return
			}

			u, err := url.Parse(baseURL)
			fatalOnErr(err)

			u.Path = path.Join(u.Path, strings.Join(args, ","))

			req, err := http.NewRequestWithContext(cmd.Context(), "GET", u.String(), nil)
			fatalOnErr(err)

			resp, err := http.DefaultClient.Do(req)
			fatalOnErr(err)

			f, err := os.Create(filename)
			fatalOnErr(err)
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			fatalOnErr(err)
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "f", ".gitignore", "Filename of gitignore file")
	cmd.Flags().BoolVar(&showBashCompletionInstead, "bash-autocomplete", false, "Print bash autocompletion instead of generating gitignore")
	cmd.Flags().BoolVar(&showZshCompletionInstead, "zsh-autocomplete", false, "Print zsh autocompletion instead of generating gitignore")
	return &cmd
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
