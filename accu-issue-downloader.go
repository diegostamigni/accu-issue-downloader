package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
)

const (
	lastIssue = 130
)

func main() {
	// let's create the folder containing the issues
	dirname, err := createOverloadFolder()
	if err != nil {
		log.Fatal(err)
	}
	counter := 0
	for counter <= lastIssue {
		filename := fmt.Sprintf("Overload%d.pdf", counter)
		body, errD := downloadIssue(filename, dirname)
		if errD != nil {
			log.Println(errD)
		} else {
			if errW := writeIssueToDisk(filename, dirname, body); errW != nil {
				log.Println(errW)
			}
		}
		counter++
	}
}

func createOverloadFolder() (string, error) {
	dirname := "Overloads"
	if _, err := os.Stat(dirname); err == nil {
		return dirname, nil
	}
	return dirname, os.Mkdir(dirname, 0700)
}

func downloadIssue(filename string, folder string) (io.ReadCloser, error) {
	fileout := fmt.Sprintf("%s/%s", folder, filename)
	if _, err := os.Stat(fileout); err == nil {
		return nil, 
			fmt.Errorf("Issue '%s' already exist, skipping it.", fileout)
	}
	url := fmt.Sprintf("http://accu.org/var/uploads/journals/%s", filename)
	log.Printf("Downloading issue '%s' from '%s'", filename, url)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {
		log.Printf("Issue '%s' successfully downloaded", filename)
		return resp.Body, nil
	}
	return nil, 
		fmt.Errorf("Issue '%s' not downloaded because or error: %v " + 
			"with http StatusCode %v", filename, err, resp.StatusCode)
}

func writeIssueToDisk(filename string, 
	dirname string, body io.ReadCloser) error {
	out, errC := os.Create(dirname + "/" + filename)
	defer out.Close()
	defer body.Close()
	if errC != nil {
		return fmt.Errorf("Cannot write '%s' to disk", filename)
	}
	io.Copy(out, body)
	return nil
}
