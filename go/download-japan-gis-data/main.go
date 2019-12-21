package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const urlTemplate = "http://nlftp.mlit.go.jp/isj/dls/data/%s/%02d000-%s.zip"

// e.g. go run main.go 17.0a /home/xxx/gis
func main() {
	args := os.Args[1:]
	if len(args) <= 1 {
		panic(fmt.Errorf("arguments are not enough"))
	}
	if err := bulkDownload(args[0], args[1]); err != nil {
		panic(err)
	}
	fmt.Println("download finished")
}

func bulkDownload(ver, p string) (err error) {
	if ver == "" {
		return fmt.Errorf("version is empty")
	}
	if p == "" {
		return fmt.Errorf("download path is empty")
	}

	for i := 1; i <= 47; i++ {
		u := getURL(ver, i)
		fp := filepath.Join(p, fmt.Sprintf("%02d000-%s.zip", i, ver))
		if err = downloadFile(u, fp); err != nil {
			return fmt.Errorf("fail to request : %s", u)
		}
	}
	return nil
}

func downloadFile(u, p string) (err error) {
	// Create the file
	out, err := os.Create(p)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getURL(ver string, no int) string {
	return fmt.Sprintf(urlTemplate, ver, no, ver)
}
