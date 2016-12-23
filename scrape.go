package main

import (
	_ "github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"os"
	"bufio"
	"net/http"
	"io"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"encoding/csv"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

type Blip struct {
	link  string
	title string
}

type File struct {
	title string
	path  string
}

type ScrapedFile struct {
	file  File
	blips []map[string]string
}

func findBlips(pathToRadar string) []Blip {
	blips := make([]Blip, 0)

	file, err := os.Open(pathToRadar)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(file))
	checkError(err)

	doc.Find(".blip").Each(func(i int, s *goquery.Selection) {
		item := s.Find("a")
		title := item.Text()
		link, _ := item.Attr("href")
		blips = append(blips, Blip{title: title, link: link })
	})

	return blips
}

func downloadBlip(blip Blip) File {
	parts := strings.Split(blip.link, "/")
	fileName := "rawData/items/" + parts[len(parts) - 1]

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		resp, err := http.Get("http://www.thoughtworks.com" + blip.link)
		checkError(err)
		body := resp.Body

		file, err := os.Create(fileName)
		checkError(err)

		io.Copy(bufio.NewWriter(file), body)
		file.Close()
		body.Close()
	}

	return File{title: blip.title, path: fileName }
}

func scrapeFile(fileToScrape File) ScrapedFile {
	file, err := os.Open(fileToScrape.path)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(file))
	checkError(err)
	file.Close()

	var blips []map[string]string
	doc.Find("div.blip-timeline-item").Each(func(i int, s *goquery.Selection) {
		blip := make(map[string]string, 0)
		blip["time"] = s.Find("div.blip-timeline-item__time").First().Text()
		blip["outcome"] = strings.Trim(s.Find("div.blip-timeline-item__ring span").First().Text(), " ")
		blip["description"] = s.Find("div.blip-timeline-item__lead").First().Text()
		blips = append(blips, blip)
	})

	return ScrapedFile{file:fileToScrape, blips:blips}
}

func main() {
	var filesCompleted chan ScrapedFile = make(chan ScrapedFile)
	defer close(filesCompleted)

	blips := findBlips("rawData/twRadar.html")

	var filesToScrape []File
	for _, blip := range blips {
		filesToScrape = append(filesToScrape, downloadBlip(blip))
	}

	var filesScraped []ScrapedFile
	for _, file := range filesToScrape {
		filesScraped = append(filesScraped, scrapeFile(file))
	}

	blipsCsvFile, _ := os.Create("import/blipsSingle.csv")
	writer := csv.NewWriter(blipsCsvFile)
	defer blipsCsvFile.Close()

	writer.Write([]string{"technology", "date", "suggestion" })
	for _, scrapedFile := range filesScraped {
		fmt.Println(scrapedFile.file.title)
		for _, blip := range scrapedFile.blips {
			writer.Write([]string{scrapedFile.file.title, blip["time"], blip["outcome"] })
		}
	}
	writer.Flush()
}
