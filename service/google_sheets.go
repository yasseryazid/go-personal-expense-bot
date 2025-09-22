package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheets struct {
	srv     *sheets.Service
	sheetID string
}

func NewGoogleSheets(credFile, sheetID string) (*GoogleSheets, error) {
	ctx := context.Background()
	b, err := os.ReadFile(credFile)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &GoogleSheets{srv: srv, sheetID: sheetID}, nil
}

func (g *GoogleSheets) Save(category string, amount int, desc string, userID int64) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	readRange := "A:A"
	resp, err := g.srv.Spreadsheets.Values.Get(g.sheetID, readRange).Do()
	if err != nil {
		return err
	}
	nextRow := len(resp.Values) + 1

	formula := fmt.Sprintf("=SUM($D$2:D%d)", nextRow)

	values := []interface{}{timestamp, userID, desc, amount, category, formula}
	rb := &sheets.ValueRange{Values: [][]interface{}{values}}

	_, err = g.srv.Spreadsheets.Values.Append(g.sheetID, "A:F", rb).
		ValueInputOption("USER_ENTERED").Do()
	return err
}

func (g *GoogleSheets) GetMonthlyTotalByUser(userID int64) (int, error) {
	readRange := "A:D"
	resp, err := g.srv.Spreadsheets.Values.Get(g.sheetID, readRange).Do()
	if err != nil {
		return 0, err
	}

	currentYear, currentMonth, _ := time.Now().Date()
	total := 0

	for i, row := range resp.Values {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			continue
		}

		timestamp := fmt.Sprintf("%v", row[0])
		userStr := fmt.Sprintf("%v", row[1])
		amountStr := fmt.Sprintf("%v", row[3])

		t, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			continue
		}

		if t.Year() == currentYear && t.Month() == currentMonth {
			if fmt.Sprintf("%d", userID) == userStr {
				if num, convErr := strconv.Atoi(amountStr); convErr == nil {
					total += num
				}
			}
		}
	}
	return total, nil
}

func (g *GoogleSheets) GetMonthlyDataByUser(userID int64) ([][]string, error) {
	readRange := "A:E"
	resp, err := g.srv.Spreadsheets.Values.Get(g.sheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	currentYear, currentMonth, _ := time.Now().Date()
	var records [][]string

	for i, row := range resp.Values {
		if i == 0 {
			continue
		}
		if len(row) < 5 {
			continue
		}

		timestamp := fmt.Sprintf("%v", row[0])
		userStr := fmt.Sprintf("%v", row[1])

		// parse tanggal
		t, err := time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			continue
		}

		if t.Year() == currentYear && t.Month() == currentMonth {
			if userStr == fmt.Sprintf("%d", userID) {
				record := []string{
					timestamp,
					userStr,
					fmt.Sprintf("%v", row[2]),
					fmt.Sprintf("%v", row[3]),
					fmt.Sprintf("%v", row[4]),
				}
				records = append(records, record)
			}
		}
	}
	return records, nil
}
