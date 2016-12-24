package scrape

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"bufio"
	"fmt"
	"log"
	"strings"
	"net/http"
	"io"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

type Blip struct {
	Link  string
	Title string
}

func (blip Blip) Download() File {
	parts := strings.Split(blip.Link, "/")
	fileName := "rawData/items/" + parts[len(parts) - 1]

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		resp, err := http.Get("http://www.thoughtworks.com" + blip.Link)
		checkError(err)
		body := resp.Body

		file, err := os.Create(fileName)
		checkError(err)

		io.Copy(bufio.NewWriter(file), body)
		file.Close()
		body.Close()
	}

	return File{Title: blip.Title, Path: fileName }
}

type File struct {
	Title string
	Path  string
}

func (fileToScrape File ) Scrape() ScrapedFile {
	file, err := os.Open(fileToScrape.Path)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(file))
	checkError(err)
	file.Close()

	var entries []map[string]string
	doc.Find("div.blip-timeline-item").Each(func(i int, s *goquery.Selection) {
		entry := make(map[string]string, 0)
		entry["time"] = s.Find("div.blip-timeline-item__time").First().Text()
		entry["outcome"] = strings.Trim(s.Find("div.blip-timeline-item__ring span").First().Text(), " ")
		entry["description"] = s.Find("div.blip-timeline-item__lead").First().Text()
		entries = append(entries, entry)
	})

	return ScrapedFile{File:fileToScrape, Entries:entries}
}

type ScrapedFile struct {
	File    File
	Entries []map[string]string
}

func FindBlips(pathToRadar string) []Blip {
	blips := make([]Blip, 0)

	file, err := os.Open(pathToRadar)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(file))
	checkError(err)

	doc.Find(".blip").Each(func(i int, s *goquery.Selection) {
		item := s.Find("a")
		title := item.Text()
		link, _ := item.Attr("href")
		blips = append(blips, Blip{Title: title, Link: link })
	})

	return blips
}