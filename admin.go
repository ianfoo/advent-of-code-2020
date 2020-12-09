package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ianfoo/advent-of-code-2020/internal/leaderboard"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "bootstrap",
				Aliases: []string{"b"},
				Usage:   "Bootstrap a new puzzle directory",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "day",
						Usage:   "Number of day to bootstrap",
						Aliases: []string{"d"},
					},
					&cli.UintFlag{
						Name:    "year",
						Usage:   "Event year",
						Aliases: []string{"y"},
						Value:   uint(time.Now().Year()),
					},
					&cli.StringFlag{
						Name:    "template",
						Usage:   "File to copy into destination",
						Aliases: []string{"t"},
						Value:   filepath.Join("templates", "puzzle.go.tmpl"),
					},
					&cli.StringFlag{
						Name:  "puzzle-root",
						Usage: "Top-level puzzle directory",
						Value: "puzzles",
					},
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Force bootstrap even if directory/file already exists",
					},
				},
				Action: BootstrapNewDay,
			},
			{
				Name:    "leaderboard",
				Aliases: []string{"lb"},
				Usage:   "Show the current standings of private leaderboard",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:  "id",
						Usage: "Private leaderboard ID",
					},
					&cli.StringFlag{
						Name:    "token",
						Usage:   "Session token value",
						EnvVars: []string{"AOC_SESSION_TOKEN"},
					},
					&cli.UintFlag{
						Name:  "year",
						Usage: "Event year for leaderboard",
						Value: uint(time.Now().Year()),
					},
				},
				Action: DisplayLeaderboard,
			},
		},
	}

	return app.Run(os.Args)
}

// BootstrapNewDay creates a new directory and starting file for a new day
// of Advent of Code, rendering a template with optional placeholders.
// The us
func BootstrapNewDay(c *cli.Context) error {

	determineLikelyDay := func() (int, error) {
		// Advent of Code releases new puzzles at midnight in the Eastern time
		// zone of the USA.
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			return 0, fmt.Errorf("cannot determine day to init: %w", err)
		}

		now := time.Now().In(loc)

		// If it's within an hour of the new puzzle dropping, set up for
		// the next day. Otherwise set up for the current day.
		day := now.Day()
		if now.Hour() == 23 {
			fmt.Println("Assuming bootstrap is for next day because of proximity")
			day++
		}

		return day, nil
	}

	var (
		year         = int(c.Uint("year"))
		day          = int(c.Uint("day"))
		templatePath = c.String("template")
		puzzleRoot   = c.String("puzzle-root")
		noClobber    = !c.Bool("force")
	)
	if day == 0 {
		likely, err := determineLikelyDay()
		if err != nil {
			return err
		}
		day = likely
	}

	var (
		dayStr    = fmt.Sprintf("%02d", day)
		dayDir    = "day-" + dayStr
		targetDir = filepath.Join(puzzleRoot, strconv.Itoa(year), dayDir)
	)
	if err := os.MkdirAll(targetDir, os.FileMode(os.FileMode(0755))); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}
	var (
		fileName = dayDir + ".go"
		filePath = filepath.Join(targetDir, fileName)
	)
	if _, err := os.Stat(filePath); err == nil && noClobber {
		return fmt.Errorf("target file %s already exists: aborting", filePath)
	}
	fmt.Printf("rendering template %s into %s\n", templatePath, filePath)
	tmpl, err := template.New("aoc-template").ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("cannot read template: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	tmplData := struct {
		Day, Year int
	}{
		Day:  day,
		Year: year,
	}
	templateName := filepath.Base(templatePath)
	if err := tmpl.ExecuteTemplate(f, templateName, tmplData); err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	return nil
}

func DisplayLeaderboard(c *cli.Context) error {
	var (
		year          = c.Uint("year")
		leaderboardID = c.Uint("id")
		token         = c.String("token")

		lb  leaderboard.Leaderboard
		err error
	)

	// Read from stdin if missing required params.
	if leaderboardID == 0 && token == "" {
		lb, err = leaderboard.FromReader(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading leaderboard: %w", err)
		}
		fmt.Println(lb)
		return nil
	}

	// Fetch from internet if ID and session spcified.
	lb, err = leaderboard.Fetch(http.DefaultClient, year, leaderboardID, token)
	if err != nil {
		return fmt.Errorf("fetching leaderboard: %w", err)
	}

	fmt.Println(lb)
	return nil
}
