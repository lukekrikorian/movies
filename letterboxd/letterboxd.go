package letterboxd

import (
	"encoding/csv"
	"fmt"
	"os"
)

const (
	DateAdded = iota
	Title
	Year
	URL
)

func Watchlist(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return records
}
