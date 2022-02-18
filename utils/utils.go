package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//util function to get the song name from a command
func GetSongTitleFromCommand(command string) (string, bool) {
	if strings.Contains(command, "!play ") {
		return strings.Split(command, "!play ")[1], true
	}

	return "", false
}

//util function to get the playlist name from a command
func GetPlalistNameFromCreateCommand(command string) (string, bool) {
	if strings.Contains(command, "!playlist create ") {
		return strings.Split(command, "!playlist create ")[1], true
	}

	return "", false
}

//util function to get the playlist name and a new song name to add from a command
func GetPlalistNameAndTrackFromCommand(command string) (plName string, trackName string, ok bool) {
	if !strings.Contains(command, "-p") || !strings.Contains(command, "-t") {
		ok = false
		return
	}
	commandParams := strings.Split(command, "-p")[1]
	plName = strings.Split(commandParams, " -t ")[0]
	plName = strings.TrimSpace(plName)
	fmt.Println(plName)
	trackName = strings.Split(commandParams, " -t ")[1]
	trackName = strings.TrimSpace(trackName)
	fmt.Println(trackName)
	ok = true
	return
}

//util function to get the playlist name from a play command passed to bot
func GetPlalistNameFromPlayCommand(command string) (plName string, ok bool) {
	if !strings.Contains(command, "-p") {
		ok = false
		return
	}
	commandParams := strings.Split(command, "-p")[1]
	plName = strings.TrimSpace(commandParams)
	ok = true
	return
}

//util function that shuffles an array of strings
func Shuffle(src []string) []string {
	arr := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		arr[v] = src[i]
	}

	return arr
}
