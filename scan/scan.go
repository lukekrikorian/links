package scan

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ScanInfo struct {
	URL     string   `json:"url"`
	Type    string   `json:"type"`
	Author  string   `json:"author"`
	Title   string   `json:"title"`
	Tags    []string `json:"tags"`
	Comment string   `json:"comment"`
}

var overrides = map[string]string{
	"www.youtube.com":   "video",
	"www.wikipedia.org": "article",
	"github.com":        "repository",
	"music.apple.com":   "audio",
	"open.spotify.com":  "audio",
}

func scan(r *http.Response) ScanInfo {
	defer r.Body.Close()

	info := ScanInfo{
		URL:  r.Request.URL.String(),
		Type: "webpage",
	}

	if override, ok := overrides[r.Request.URL.Hostname()]; ok {
		info.Type = override
	}

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		log.Println(err)
		return info
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		key, ok := s.Attr("name")
		if !ok {
			return
		}

		content, ok := s.Attr("content")

		if !ok || content == "" {
			return
		}

		switch key {
		case "title":
			info.Title = content
		case "description":
			info.Comment = content
		case "author":
			info.Author = content
		case "keywords":
			info.Tags = strings.Split(content, ", ")
		}
	})

	if info.Title == "" {
		info.Title = doc.Find("title").Text()
	}

	return info
}

func HandleScan(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Url"]; !ok {
		http.Error(w, "Missing url header", http.StatusBadRequest)
		return
	}

	url := r.Header["Url"][0]
	res, err := http.Get(url)

	if err != nil || res.StatusCode != 200 {
		http.Error(w, "Couldn't load page", 500)
		return
	}

	if _, ok := res.Header["Content-Type"]; !ok {
		http.Error(w, "Unknown content type", 500)
		return
	}

	var info ScanInfo

	switch strings.Split(res.Header["Content-Type"][0], ";")[0] {
	case "text/html":
		log.Println("Scanning", url)
		info = scan(res)
	case "application/pdf":
		info = ScanInfo{Type: "pdf"}
	case "image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml":
		info = ScanInfo{Type: "image"}
	case "video/mp4",
		"video/mpeg",
		"video/webm":
		info = ScanInfo{Type: "film"}
	case "audio/mp3",
		"audio/mpeg",
		"audio/webm":
		info = ScanInfo{Type: "audio"}
	case "text/plain":
		info = ScanInfo{Type: "text"}
	case "application/msword",
		"application/epub+zip":
		info = ScanInfo{Type: "document"}
	default:
		log.Println("Unknown content type", res.Header["Content-Type"][0])
		http.Error(w, "Unknown content type", 500)
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
