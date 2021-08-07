# google-sheets-adder

Simple tool for adding new entries to the first column of a google sheet (really customized for my uses)

# Setup

_Note_ This is just following the quickstart guide from the [Google quickstart documentation for the Sheets API with Go](https://developers.google.com/sheets/api/quickstart/go).

**Step One**, you're going to have to create a project in your GCP and enable the API. Steps found [here](https://developers.google.com/workspace/guides/create-project).

**Step Two**, save your credentials in the same directory as this repo under the filename `credentials.json`.

**Step Three**, run the program! `go run *.go` will handle it for ya. The output will request you to verify the auth by providing a link. Go to that link, verify the app (if you're ok with that), and copy and paste the code given in your browser into your terminal.

**Step Four**, you'll need to provide a Sheets ID to the program in order for it to know which Sheet to add your items to. You can find the ID from the URL of a Sheet.

ie: `https://docs.google.com/spreadsheets/d/<ID will be here>/edit#gid=0`

Copy and paste that ID into the terminal.

**All done!** Now you can proceed to add items to your google Sheet via this program.
