package mongo

import (
	"fmt"
	"musicbot/soundcloud"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	MONGO_HOST         = "mongodb://localhost:27017"
	MONGO_DBNAME       = "playlists"
	COLECTION_PLAYLIST = "Playlist"
)

type Playlist struct {
	PlaylistName string
	Songs        []string
}

type Mongo struct {
}

func NewMongo() *Mongo {
	return &Mongo{}
}

//this method adds a new empty user defined playlist to the database, with name playlistName
func (m *Mongo) NewPlayList(playlistName string) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Release resource when main function is returned.
	defer Close(client, ctx, cancel)

	var songs []string
	pl := Playlist{PlaylistName: playlistName, Songs: songs}
	insertOneResult, err := InsertOne(client, ctx, MONGO_DBNAME, COLECTION_PLAYLIST, pl)
	if err != nil {
		panic(err)
	}

	fmt.Print(insertOneResult)
}

//this method adds a new song to the playlist with the given name
func (m *Mongo) AddSongToPlaylist(dataBase, col string, query, field interface{}, songTitle string) (result Playlist) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Release resource when main function is returned.
	defer Close(client, ctx, cancel)

	playlist := FindOne(client, ctx, dataBase, col, query, field)
	fmt.Println(playlist)
	var pl Playlist
	playlist.Decode(&pl)

	fmt.Print(pl)
	pl.Songs = append(pl.Songs, songTitle)
	update := bson.D{
		{"$set", bson.D{
			{"songs", pl.Songs},
		}},
	}

	u, err := UpdateOne(client, ctx, dataBase, col, query, update)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(u.ModifiedCount)

	return
}

func (m *Mongo) AddSong(playlistName string, song soundcloud.SongModel) (Playlist, bool) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer Close(client, ctx, cancel)

	filter := bson.M{
		"playlistname": playlistName,
	}
	fmt.Print(filter)

	return m.AddSongToPlaylist(MONGO_DBNAME, COLECTION_PLAYLIST, filter, nil, fmt.Sprint(song.Media.Transcodings[1].URL)), true
}

func (m *Mongo) GetPlaylist(playlistName string) (pl Playlist, ok bool) {
	client, ctx, cancel, err := Connect(MONGO_HOST)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer Close(client, ctx, cancel)

	filter := bson.M{
		"playlistname": playlistName,
	}
	fmt.Print(filter)

	playlist := FindOne(client, ctx, MONGO_DBNAME, COLECTION_PLAYLIST, filter, nil)
	er := playlist.Decode(&pl)

	if er != nil {
		ok = false
		return
	}
	ok = true

	return
}
