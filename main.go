package main

import (
	"flag"
	"musicbot/bot"
	"musicbot/soundcloud"
)

// Variables used for command line parameters
var (
	Token      string
	Bot        bot.Bot
	SoundCloud soundcloud.SoundCloud
)

//constants for user interaction
const (
	SONG_NOT_FOUND         = "Couldn't find the song, try another one :("
	SONG_PLAYING           = "Now playing"
	FOOTER_TEXT            = "Music bot"
	JOIN_CHANNEL_FIRST     = "Join to voice channel first for bot to join"
	JOIN_CHANNEL_FIRST_BOT = "Join to voice channel first for bot to play a song"
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	Bot = *bot.NewBot()
}

func main() {

	Bot.Start(Token)
}

//todo for v2: / commands
