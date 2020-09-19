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

func main() {
	newRootCMD().ExecuteContext(context.Background())
}

const baseURL = "https://www.toptal.com/developers/gitignore/api"

func newRootCMD() *cobra.Command {
	var filename string
	cmd := cobra.Command{
		Use:  "gitignore [names...]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			u, err := url.Parse(baseURL)
			fatalOnErr(err)

			u.Path = path.Join(u.Path, strings.Join(args, ","))

			req, err := http.NewRequestWithContext(cmd.Context(), "GET", u.String(), nil)
			fatalOnErr(err)

			resp, err := http.DefaultClient.Do(req)
			fatalOnErr(err)

			// f, err := os.Open(filename)
			f, err := os.Create(filename)
			fatalOnErr(err)
			defer f.Close()

			_, err = io.Copy(f, resp.Body)
			fatalOnErr(err)
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "f", ".gitignore", "Filename of gitignore file")
	return &cmd
}

func fatalOnErr(err error) {
	if err != nil {
		printFatal(err)
	}
}

func printFatal(a ...interface{}) {
	fmt.Fprint(os.Stderr, "error:")
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
