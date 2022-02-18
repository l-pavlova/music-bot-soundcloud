package mocks

import (
	"musicbot/mongo"
	"musicbot/soundcloud"
)

var (
	GetAddSongToPlayListFunc func(dataBase, col string, query, field interface{}, songTitle string) (result mongo.Playlist)
	GetAddSongFunc           func(playlistName string, song soundcloud.SongModel) (result mongo.Playlist, ok bool)
	GetNewPlaylistFunc       func(playlistName string)
	GetGetPlaylistFunc       func(playlistName string) (pl mongo.Playlist, ok bool)
)

type MongoClientMock struct {
	AddSongToPlayListFunc func(playlistName string) mongo.Playlist
	AddSongFunc           func(playlistName string, song soundcloud.SongModel) (result mongo.Playlist)
	NewPlayListFunc       func(playlistName string)
	GetPlaylistFunc       func(playlistName string) (pl mongo.Playlist, ok bool)
}

func (m *MongoClientMock) AddSongToPlaylist(dataBase, col string, query, field interface{}, songTitle string) (result mongo.Playlist) {
	return GetAddSongToPlayListFunc(dataBase, col, query, field, songTitle)
}

func (m *MongoClientMock) NewPlayList(playlistName string) {
	GetNewPlaylistFunc(playlistName)
}
func (m *MongoClientMock) AddSong(playlistName string, song soundcloud.SongModel) (result mongo.Playlist, ok bool) {
	return GetAddSongFunc(playlistName, song)
}
func (m *MongoClientMock) GetPlaylist(playlistName string) (pl mongo.Playlist, ok bool) {
	return mongo.Playlist{}, true
}
