package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func CreateCsvWriter(fileName string) (*os.File, *csv.Writer) {
	csvFile, err := os.Create(fileName)

	if err != nil {
		log.Fatalln(err)
	}

	csvWriter := csv.NewWriter(csvFile)
	return csvFile, csvWriter
}

func WriteInCsv(csvWriter *csv.Writer, value []string) {
	err := csvWriter.Write(value)

	if err != nil {
		log.Fatalln(err)
	}
}
