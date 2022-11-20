package games

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/gocarina/gocsv"
)

type Game struct {
	Date       GameTime `csv:"Date"`
	Team1      string   `csv:"Team 1"`
	Match      string   `csv:"Match"`
	Team2      string   `csv:"Team 2"`
	Group      string   `csv:"Group"`
	Location   string   `csv:"Location"`
	Stage      string   `csv:"Stage"`
	Place      string   `csv:"Place"`
	GoalsTeam1 int      `csv:"Goals Team 1"`
	GoalsTeam2 int      `csv:"Goals Team 2"`
	Winner     int      `csv:"Winner"`
}

func (g *Game) String() string {
	if time.Now().After(g.Date.Time) {
		return fmt.Sprintf(
			"*%s* as *%s* | %s\n%s %s\n%s %dx%d %s\n%s\n\n",
			g.Date.Format("2/01"),
			g.Date.Format("15:04"),
			g.Place,
			g.Stage,
			g.Group,
			g.Team1,
			g.GoalsTeam1,
			g.GoalsTeam2,
			g.Team2,
			g.Location,
		)
	}

	return fmt.Sprintf(
		"*%s* as *%s* | %s\n%s %s\n%s x %s\n%s\n\n",
		g.Date.Format("2/01"),
		g.Date.Format("15:04"),
		g.Place,
		g.Stage,
		g.Group,
		g.Team1,
		g.Team2,
		g.Location,
	)
}

type GameTime struct {
	time.Time
}

func (date *GameTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2-Jan-06 3:04 PM", csv)
	return err
}

//go:embed games.csv
var rawGames []byte

func GetGames() ([]Game, error) {
	allGames := []Game{}
	if err := gocsv.UnmarshalBytes(rawGames, &allGames); err != nil {
		return nil, fmt.Errorf("error unmarshalling games.csv: %w", err)
	}

	return allGames, nil
}
