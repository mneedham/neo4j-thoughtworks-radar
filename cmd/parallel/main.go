package main

import (
	"os"
	"encoding/csv"
	"github.com/mneedham/neo4j-thoughtworks-radar/scrape"
)

func main() {
	var filesToScrape chan scrape.File = make(chan scrape.File)
	var filesScraped chan scrape.ScrapedFile = make(chan scrape.ScrapedFile)
	defer close(filesToScrape)
	defer close(filesScraped)

	// wget https://www.thoughtworks.com/radar/a-z -O rawData/twRadar.html
	blips := scrape.FindBlips("rawData/twRadar.html")

	for _, blip := range blips {
		go func(blip scrape.Blip) { filesToScrape <- blip.Download() }(blip)
	}

	for i := 0; i < len(blips); i++ {
		select {
		case file := <-filesToScrape:
			go func(file scrape.File) { filesScraped <- file.Scrape() }(file)
		}
	}

	blipsCsvFile, _ := os.Create("import/blips.csv")
	writer := csv.NewWriter(blipsCsvFile)
	defer blipsCsvFile.Close()

	writer.Write([]string{"technology", "date", "suggestion" })
	for i := 0; i < len(blips); i++ {
		select {
		case scrapedFile := <-filesScraped:
			for _, blip := range scrapedFile.Entries {
				writer.Write([]string{scrapedFile.File.Title, blip["time"], blip["outcome"] })
			}
		}
	}
	writer.Flush()
}