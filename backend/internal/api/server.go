package api

import (
	"fmt"
	"net/http"

	"github.com/X-Calibre/MasjidPi/backend/internal/version"
)

// Start launches the HTTP server.
func Start(addr string) error {
	http.HandleFunc("/", home)

	fmt.Printf("HTTP server listening on %s\n", addr)

	return http.ListenAndServe(addr, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(
		w,
		"%s is running\nVersion: %s\n",
		version.AppName,
		version.Version,
	)
}
