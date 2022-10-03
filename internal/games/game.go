package games

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/gocarina/gocsv"
)

type Game struct {
	Date     GameTime `csv:"Date"`
	Team1    string   `csv:"Team 1"`
	Match    string   `csv:"Match"`
	Team2    string   `csv:"Team 2"`
	Group    string   `csv:"Group"`
	Location string   `csv:"Location"`
	Stage    string   `csv:"Stage"`
	Place    string   `csv:"Place"`
}

type GameTime struct {
	time.Time
}

func (date *GameTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2-Jan-06 3:04 PM", csv)
	return err
}

// TODO: maybe not embed

//go:embed games.csv
var rawGames []byte

func GetGames() ([]Game, error) {
	allGames := []Game{}
	if err := gocsv.UnmarshalBytes(rawGames, &allGames); err != nil {
		return nil, fmt.Errorf("error unmarshalling games.csv: %w", err)
	}

	return allGames, nil
}
