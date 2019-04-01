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

	http.ListenAndServe(":8080", defaultMux(&jsonMap))
}

func introHandler(introSegment *Segment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		displaySegData(w, introSegment)
	})
}

func optionHandler(segmentMap *map[string]Segment) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if len(r.Form.Get("option")) == 0 {
			fmt.Fprintf(w, "You selected nothing.  The end.")
		} else {
			segment := (*segmentMap)[r.Form.Get("option")]
			displaySegData(w, &segment)
		}
	})
}

func displaySegData(w http.ResponseWriter, segment *Segment) {

	fmt.Fprintf(w, "<html><body><h1>%s</h1>", segment.Title)
	for _, bodyStr := range segment.Story {
		fmt.Fprintf(w, "<p>%s</p>", bodyStr)
	}
	if 0 < len(segment.Options) {
		fmt.Fprintf(w, "<form action=/choose method=\"get\">")
		for _, option := range segment.Options {
			fmt.Fprintf(w, "<input type=\"radio\" name=\"option\" value=\"%s\">%s</input><br>", option.Arc, option.Text)
		}
		fmt.Fprintf(w, "<input type=\"submit\" value=\"Submit\">")
		fmt.Fprintf(w, "</form>")
	}
	fmt.Fprintf(w, "</body></html>")
}

func defaultMux(segmentMap *map[string]Segment) *http.ServeMux {

	mux := http.NewServeMux()

	introSegment := (*segmentMap)["intro"]
	introHandler := introHandler(&introSegment)
	mux.Handle("/", introHandler)

	optionHandler := optionHandler(segmentMap)
	mux.Handle("/choose", optionHandler)
	return mux
}
