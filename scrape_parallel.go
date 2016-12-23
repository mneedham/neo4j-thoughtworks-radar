package main

import (
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
		log.Fatal(err)
	}
}

func findBlips(pathToRadar string) []Blip {

	file, err := os.Open(pathToRadar)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(file))
	checkError(err)

	var blips []Blip
	doc.Find(".blip").Each(func(i int, s *goquery.Selection) {
		item := s.Find("a")
		attr, _ := item.Attr("href")
		blips = append(blips, Blip{title: item.Text(), link: attr})
	})

	return blips
}

func downloadBlip(blip Blip, filesToScrape chan <- File) {
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

	filesToScrape <- File{title: blip.title, path: fileName }
}

func scrapeFile(fileToScrape File, filesScraped chan <- ScrapedFile) {
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

	filesScraped <- ScrapedFile{file:fileToScrape, blips:blips}
}

func main() {
	var filesToScrape chan File = make(chan File)
	var filesScraped chan ScrapedFile = make(chan ScrapedFile)
	defer close(filesToScrape)
	defer close(filesScraped)

	blips := findBlips("rawData/twRadar.html")

	for _, blip := range blips {
		go downloadBlip(blip, filesToScrape)
	}

	for i := 0; i < len(blips); i++ {
		select {
		case file := <-filesToScrape:
			go scrapeFile(file, filesScraped)
		}
	}

	blipsCsvFile, _ := os.Create("import/blips.csv")
	writer := csv.NewWriter(blipsCsvFile)
	defer blipsCsvFile.Close()

	writer.Write([]string{"technology", "date", "suggestion" })
	for i := 0; i < len(blips); i++ {
		select {
		case scrapedFile := <-filesScraped:
			fmt.Println(scrapedFile.file.title)
			for _, blip := range scrapedFile.blips {
				writer.Write([]string{scrapedFile.file.title, blip["time"], blip["outcome"] })
			}
		}
	}
	writer.Flush()
}