package main

import (
	"encoding/csv"
	"os"
	"github.com/mneedham/neo4j-thoughtworks-radar/scrape"
	"fmt"
)

func main() {
	var filesCompleted chan scrape.ScrapedFile = make(chan scrape.ScrapedFile)
	defer close(filesCompleted)

	// wget https://www.thoughtworks.com/radar/a-z -O rawData/twRadar.html 
	blips := scrape.FindBlips("rawData/twRadar.html")

	var filesToScrape []scrape.File
	for _, blip := range blips {
		filesToScrape = append(filesToScrape, blip.Download())
	}

	var filesScraped []scrape.ScrapedFile
	for _, file := range filesToScrape {
		filesScraped = append(filesScraped, file.Scrape())
	}

	blipsCsvFile, _ := os.Create("import/blipsSingle.csv")
	writer := csv.NewWriter(blipsCsvFile)
	defer blipsCsvFile.Close()

	writer.Write([]string{"technology", "date", "suggestion" })
	for _, scrapedFile := range filesScraped {
		fmt.Println(scrapedFile.File.Title)
		for _, blip := range scrapedFile.Entries {
			writer.Write([]string{scrapedFile.File.Title, blip["time"], blip["outcome"] })
		}
	}
	writer.Flush()
}
