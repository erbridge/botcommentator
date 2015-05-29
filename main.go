package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
	"github.com/erbridge/gotwit"
	"github.com/erbridge/gotwit/callback"
	"github.com/erbridge/gotwit/twitter"
)

const (
	Blank = `⬜`
	Sun   = `🌞`
	Moon  = `🌚`
)

var (
	sunWinRegexp  = regexp.MustCompile("Sun[[:space:]]*Wins!![[:space:]]*$")
	moonWinRegexp = regexp.MustCompile("Moon[[:space:]]*Wins!![[:space:]]*$")
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
		} else {
			return
		}

		text += " " + fmt.Sprintf("https://twitter.com/%s/status/%s", t.User.ScreenName, t.IdStr)

		b.Post(text, false)
	}
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
