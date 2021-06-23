package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
    "flag"
    "strings"

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
    
    fmt.Println("Please enter the Google Sheet ID where you would like to save verbs: ")
    if _, err := fmt.Scan(&sId); err != nil {
        log.Fatalf("Unable to read Sheet ID from user input: %v", err)
    }

    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

    if err != nil {
        log.Fatalf("Unable to open Sheet ID file: %v", err)
    }

    defer f.Close()
    err = json.NewEncoder(f).Encode(sId)

    if err != nil {
        log.Fatalf("Unable to cache sheet id: %v", err)
    }

    return sId
}

func removeCredentialsFiles() {
    var c string
    fmt.Println("Are you sure you want to clear the credentials? This will remove the auth credentials, Oauth token, and your Sheet ID. [y/n]")

    if _, err := fmt.Scan(&c); err != nil {
        log.Fatalf("Unable to read user input.")
    }

    c = strings.ToLower(c)
    
    if c == "y" || c == "yes" {
        paths := []string{"token.json", "credentials.json", "sheet_id.json"}
        for _, p := range paths {

            err := os.Remove(p)
            if err != nil {
                log.Fatalf("Unable to remove file \"%s\". Aborting.", p)
            }

            fmt.Printf("Removed %s.\n", p)
        }
        return
    }

    fmt.Println("Credentials were *NOT* cleared.")
}

func main() {
    // Clear tokens, auth, and sheet id subcommand
    if len(os.Args) < 2 {
		log.Fatal("Please provide at least one verb to add to your Google Sheet")
    }

    sc := os.Args[1]
    if strings.ToLower(sc) == "clear" {
        removeCredentialsFiles()
        os.Exit(0)
    }

    newIdPtr := flag.Bool("id", false, "Provide this flag if you'd like to clear the saved sheet id.")
    flag.Parse()

    verbs := flag.Args()

	if len(verbs) == 0 {
		log.Fatal("Please provide at least one verb to add to your Google Sheet")
	}


	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, "token.json")

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

    var spreadsheetId string
    if *newIdPtr == true {
        spreadsheetId = saveSheetId("sheet_id.json") 
    } else {
        spreadsheetId = getSheetId("sheet_id.json")
    }

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
