package service

import (
	"context"
	"os"
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
	values := []interface{}{timestamp, userID, desc, amount, category}
	rb := &sheets.ValueRange{Values: [][]interface{}{values}}

	_, err := g.srv.Spreadsheets.Values.Append(g.sheetID, "A:E", rb).
		ValueInputOption("RAW").Do()
	return err
}
