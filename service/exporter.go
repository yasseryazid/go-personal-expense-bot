package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportToCSV(username string, records [][]string) (string, error) {
	now := time.Now()

	safeName := strings.ReplaceAll(username, " ", "_")
	if safeName == "" {
		safeName = "expenzo_user"
	}

	fileName := fmt.Sprintf(
		"expenzo_%s_%02d_%d_%02d%02d.csv",
		safeName,
		now.Month(),
		now.Year(),
		now.Hour(),
		now.Minute(),
	)

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Tanggal", "Pengeluaran", "Jumlah Pengeluaran", "Kategori"})

	total := 0

	for _, record := range records {
		if len(record) >= 5 {
			amountStr := record[3]
			amount, _ := strconv.Atoi(amountStr)

			writer.Write([]string{
				record[0],
				record[2],
				humanize.Comma(int64(amount)),
				record[4],
			})

			total += amount
		}
	}

	writer.Write([]string{"", "TOTAL", humanize.Comma(int64(total)), ""})

	return fileName, nil
}
