package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.ListenAndServe(":7000",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "watcher is running at 7000\nPackage Args:", os.Args)
			},
		),
	)
}
