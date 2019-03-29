package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Option struct {
	Text string
	Arc  string
}

type Segment struct {
	Title   string
	Story   []string
	Options []Option
}

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err.Error())
	}

	fileSize := fileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	file.Read(fileBuffer)

	jsonMap := make(map[string]Segment)
	json.Unmarshal(fileBuffer, &jsonMap)
	fmt.Println(jsonMap)

	http.ListenAndServe(":8080", defaultMux())
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, gopher!")
}
