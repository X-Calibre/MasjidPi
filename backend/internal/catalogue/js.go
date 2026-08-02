package catalogue

import (
	"io"
	"net/http"
	"os"
)

func DownloadJS(filename string) error {

	resp, err := http.Get("https://www.livemasjid.com/js/livemasjid.js?version=4.1")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}
