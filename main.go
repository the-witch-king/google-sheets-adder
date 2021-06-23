package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getSheetId(path string) string {
    savedId, err := sheetIdFromFile(path)

    if err != nil {
        savedId = saveSheetId(path)
    }

    return savedId
}

func sheetIdFromFile(path string) (string, error) {
    f, err := os.Open(path)
    // Unable to open file
    if err != nil {
        return "", err
    }
    var s string
    err = json.NewDecoder(f).Decode(&s)
    return s, err
}

func saveSheetId(path string) string {
    var sId string
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

    if err != nil {
        log.Fatalf("Unable to open Sheet ID file: %v", err)
    }

    fmt.Println("Please enter the Google Sheet ID where you would like to save verbs: ")

    if _, err := fmt.Scan(&sId); err != nil {
        log.Fatalf("Unable to read Sheet ID from user input: %v", err)
    }

    defer f.Close()
    err = json.NewEncoder(f).Encode(sId)

    if err != nil {
        log.Fatalf("Unable to cache sheet id: %v", err)
    }

    return sId
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide at least one verb to add to your Google Sheet")
	}

    verbs := os.Args[1:]

	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

    spreadsheetId := getSheetId("sheet_id.json")

	v := make([][]interface{}, len(verbs))
    for i, verb := range verbs {
        v[i] = make([]interface{}, 1)
        v[i][0] = verb

    }

	values := &(sheets.ValueRange{Values: v})
	_, e := srv.Spreadsheets.Values.Append(spreadsheetId, "A3", values).ValueInputOption("raw").Do()

	if e != nil {
		log.Fatalf("API Error: %v", e)
	}

    fmt.Printf("Added:\n")

    for _, v := range verbs {
        fmt.Printf(" - %s\n", v)
    }
}
