package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/erbridge/gotwit"
	"github.com/erbridge/gotwit/callback"
	"github.com/erbridge/gotwit/twitter"

	"github.com/erbridge/botcommentator/connect4"
)

const (
	Blank = `⬜`
	Sun   = `🌞`
	Moon  = `🌚`
)

var (
	boardRegexp = regexp.MustCompile("(?sm)\n\n(.*)\n1⃣")
	pieceRegexp = regexp.MustCompile(fmt.Sprintf("%s|(%s)|(%s)", Blank, Sun, Moon))

	sunWinRegexp  = regexp.MustCompile("Sun[[:space:]]*Wins!*[[:space:]]*$")
	moonWinRegexp = regexp.MustCompile("Moon[[:space:]]*Wins!*[[:space:]]*$")

	sunStartsRegexp  = regexp.MustCompile("Sun[[:space:]]*to[[:space:]]*Play[[:space:]]*$")
	moonStartsRegexp = regexp.MustCompile("Moon[[:space:]]*to[[:space:]]*Play[[:space:]]*$")

	sunTurnRegexp  = regexp.MustCompile("Sun's[[:space:]]*Turn[[:space:]]*$")
	moonTurnRegexp = regexp.MustCompile("Moon's[[:space:]]*Turn[[:space:]]*$")

	nullGameRegexp = regexp.MustCompile("Null[[:space:]]*Game[[:space:]]*$")
)

func createMassConnect4Callback(b *gotwit.Bot) func(anaconda.Tweet) {
	return func(t anaconda.Tweet) {
		if t.User.ScreenName != "massconnect4" {
			return
		}

		text := ""

		if sunWinRegexp.MatchString(t.Text) {
			text += "Looks like the sun's shining today. A great game from two great teams."
		} else if moonWinRegexp.MatchString(t.Text) {
			text += "It's been a long night, but the moon's as bright as ever. Well played by both teams."
		} else if nullGameRegexp.MatchString(t.Text) {
			text += "They're calling the match off. Looks like they just couldn't take the pressure."
		} else if sunStartsRegexp.MatchString(t.Text) {
			text += "Good morning, folks! Time for another gripping game of Mass Connect 4. Sun goes first this time."
		} else if moonStartsRegexp.MatchString(t.Text) {
			text += "Good evening, folks! Time for another gripping game of Mass Connect 4. Moon goes first this time."
		} else if boardStrings := boardRegexp.FindStringSubmatch(t.Text); boardStrings != nil {
			nextPiece := 0
			nextTeam := ""
			lastTeam := ""
			if sunTurnRegexp.MatchString(t.Text) {
				nextPiece = 1
				nextTeam = "Sun"
				lastTeam = "Moon"
			} else if moonTurnRegexp.MatchString(t.Text) {
				nextPiece = -1
				nextTeam = "Moon"
				lastTeam = "Sun"
			}

			weightedCount, count := countWins(boardStrings[1], nextPiece)

			if prop := float32(nextPiece) * float32(weightedCount) / float32(count); prop > 0.9 {
				text += fmt.Sprintf("Things are looking good for %s.", nextTeam)
			} else if prop < -0.9 {
				text += fmt.Sprintf("%s's pulling away. Can %s come back from this?", lastTeam, nextTeam)
			}
		}

		if text == "" {
			return
		}

		text += fmt.Sprintf(" https://twitter.com/%s/status/%s", t.User.ScreenName, t.IdStr)

		b.Post(text, false)
	}
}

func countWins(boardString string, piece int) (weightedCount, count int) {
	board := connect4.NewBoard(6, 7)
	rows := strings.Split(boardString, "\n")
	for rowIdx := len(rows) - 1; rowIdx >= 0; rowIdx-- {
		matches := pieceRegexp.FindAllStringSubmatch(rows[rowIdx], -1)

		for colIdx, m := range matches {
			piece := 0
			if m[1] != "" {
				piece = 1
			} else if m[2] != "" {
				piece = -1
			} else {
				continue
			}

			board.Play(piece, colIdx)
		}
	}

	return board.CountWins(piece, 8)
}

func main() {
	var (
		con twitter.ConsumerConfig
		acc twitter.AccessConfig
	)

	f := "secrets.json"
	if _, err := os.Stat(f); err == nil {
		con, acc, _ = twitter.LoadConfigFile(f)
	} else {
		con, acc, _ = twitter.LoadConfigEnv()
	}

	b := gotwit.NewBot("BotCommentator", con, acc)

	b.RegisterCallback(callback.Post, createMassConnect4Callback(&b))

	b.Start()
	b.Stop()
}
