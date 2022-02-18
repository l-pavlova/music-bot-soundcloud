package bot

import (
	"fmt"
	"musicbot/mocks"
	"musicbot/mongo"
	"musicbot/soundcloud"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {

	mocks.GetAddSongFunc = func(string, soundcloud.SongModel) (mongo.Playlist, bool) {
		return mongo.Playlist{}, true
	}
	mocks.GetSearchSongFunc = func(string) soundcloud.SongModel {
		return soundcloud.SongModel{ //some default mock value
			ArtworkURL:        "",
			Caption:           nil,
			Title:             "title",
			CreatedAt:         time.Time{},
			Description:       "",
			Duration:          0,
			Genre:             "",
			ID:                0,
			LabelName:         nil,
			LikesCount:        0,
			PublisherMetadata: soundcloud.PublisherMetadata{},
			State:             "",
			Streamable:        false,
			TrackFormat:       "",
			URI:               "",
			PermalinkURL:      "",
			Media: struct {
				Transcodings []struct {
					URL      string "json:\"url\""
					Preset   string "json:\"preset\""
					Duration int    "json:\"duration\""
					Snipped  bool   "json:\"snipped\""
					Format   struct {
						Protocol string "json:\"protocol\""
						MimeType string "json:\"mime_type\""
					} "json:\"format\""
					Quality string "json:\"quality\""
				} "json:\"transcodings\""
			}{},
			Policy: "",
		}
	}

	mocks.GetGetSongFromSoundcloudURLFunc = func(url string) (string, bool) {
		return "https://cf-hls-media.sndcdn.com/media/0/31762/uCTJQknkneSb.128.mp3?Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiKjovL2NmLWhscy1tZWRpYS5zbmRjZG4uY29tL21lZGlhLyovKi91Q1RKUWtua25lU2IuMTI4Lm1wMyIsIkNvbmRpdGlvbiI6eyJEYXRlTGVzc1RoYW4iOnsiQVdTOkVwb2NoVGltZSI6MTY0NDE4NzM1M319fV19&Signature=YhG4-XUGngG6eM5IZB0AlgDoDLiJsX8lyCVa8E1E-7berxMPsbYsEU~OsIg52OHkq5P0fduEZ3~l5JLkzlBukbKeZ2YWVaWYSLfKQsvO1fr5wAG08f40yel3ZQ5F~0IPkpuycH69509cVNTc7G44sFEejkAeZ3e8hi9R-eGPRxCvxG-MwwUUmVFPD4QbnMAfJjgkK-Ov0WgkV4oTFyapVfU0~Kq574tn0cj7vhCEGUmgUr85OxTYOreRevmIogjVW6mD0ANiPWCkjwlYxJJ-xX1e~U3bzTPXyZbn~eN2dU5ydk~8pHvMM0asmSaqdf6CB149rbLk4s51JPYQYWyRzA__&Key-Pair-Id=APKAI6TU7MMXM5DG6EPQ",
			true
	}

	mocks.GetNewPlaylistFunc = func(playlistName string) {}
}
func TestAddSongToPlaylist(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	bot.sc = &mocks.SoundCloudClientMock{}
	command := "!add song -p peti -t maccarena"
	ok := bot.AddSong(command)
	if !ok {
		t.Errorf("bot.AddSong(%s) = %t; want Maccarena", command, ok)
	}
}

func TestAddSongToPlaylistInvalidCommand(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	bot.sc = &mocks.SoundCloudClientMock{}

	command := "dobavqm s greshna komanda"
	ok := bot.AddSong(command)
	if ok {
		t.Errorf("bot.AddSong(%s) = %t; want not okay result", command, ok)
	}
}

func TestGetPlaylist(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	bot.sc = &mocks.SoundCloudClientMock{}

	_, ok := bot.GetPlaylist("playlist bane")
	if !ok {
		t.Errorf("bot.GetPlaylist returned %t; want true", ok)
	}
}

func TestGetSongInvalidName(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	bot.sc = &mocks.SoundCloudClientMock{}
	s := *&discordgo.Session{}
	_, _, ok := bot.GetSong(&s, "", "", "invalid song name command")

	if ok {
		t.Errorf("bot.GetSong() returned true; want false")
	}
}

func TestGetSongValid(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	bot.sc = &mocks.SoundCloudClientMock{}
	s := *&discordgo.Session{}

	command := "!play maccarena"
	_, _, ok := bot.GetSong(&s, "", "", command)
	if !ok {
		fmt.Errorf("bot.GetSong() couldnt get song from a valid command")
	}
}

func TestNewPlayListInvalid(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	command := "!notplaylist notcreate testplaylist"
	ok := bot.NewPlayList(command)
	if ok {
		t.Errorf("bot.NewPlayList() returned true; wanted false")
	}
}
func TestNewPlayList(t *testing.T) {
	bot := *NewBot()
	bot.mongo = &mocks.MongoClientMock{}
	command := "!playlist create testplaylist"
	ok := bot.NewPlayList(command)
	if !ok {
		t.Errorf("bot.NewPlayList() returned false; wanted true")
	}
}
