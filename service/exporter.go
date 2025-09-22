package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportToCSV(username string, records [][]string) (string, error) {
	currentYear, currentMonth, _ := time.Now().Date()

	safeName := strings.ReplaceAll(username, " ", "_")
	if safeName == "" {
		safeName = "expenzo_user"
	}

	fileName := fmt.Sprintf("expenzo_%s_%02d_%d.csv", safeName, currentMonth, currentYear)

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Tanggal", "Pengeluaran", "Jumlah Pengeluaran", "Kategori"})

	for _, record := range records {
		if len(record) >= 5 {
			writer.Write([]string{
				record[0], // Timestamp
				record[2], // Description
				record[3], // Amount
				record[4], // Category
			})
		}
	}

	return fileName, nil
}
