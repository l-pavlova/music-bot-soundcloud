package mocks

import (
	"musicbot/soundcloud"
)

var (
	GetGetSongFromSoundcloudURLFunc func(url string) (string, bool)
	GetSearchSongFunc               func(name string) soundcloud.SongModel
)

type SoundCloudClientMock struct {
	GetSongFromSoundcloudURLFunc func(url string) (string, bool)
	SearchSongFunc               func(name string) soundcloud.SongModel
}

func (s *SoundCloudClientMock) GetSongFromSoundcloudURL(url string) (string, bool) {
	return GetGetSongFromSoundcloudURLFunc(url)
}

func (s *SoundCloudClientMock) SearchSong(name string) soundcloud.SongModel {
	return GetSearchSongFunc(name)
}
