package bot

import (
	"fmt"
	"musicbot/mongo"
	"musicbot/soundcloud"
	"musicbot/utils"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

var (
	guildsInfo      map[string](chan bool)
	songsPlaying    chan bool
	next            chan bool
	counter         chan int
	playlistPlaying string
)

//constants for user interaction
const (
	PLAYLIST_NOT_FOUND              = "Couldn't find the playlist, try another one :("
	BOT_JOINED                      = "Bot joined the voice channel!"
	BOT_LEFT                        = "Bot left the voice channel!"
	SONG_PLAYING_FROM_PLAYLIST      = "Now playing from playlist"
	SHOULD_CONNECT_TO_VOICE_CHANNEL = "You should connect to a voice channel before you play song! Try calling !join first"
	AVAILABLE_COMMANDS              = "Available commands for the music bot to work with"
	INVALID_COMMAND                 = "Please type a correct command, press !help to display all commands"
	SONG_NOT_FOUND                  = "Couldn't find the song, try another one :("
	SONG_PLAYING                    = "Now playing"
	FOOTER_TEXT                     = "Music bot"
	JOIN_CHANNEL_FIRST              = "Join to voice channel first for bot to join"
	JOIN_CHANNEL_FIRST_BOT          = "Join to voice channel first for bot to play a song"
)

type SoundCloudClient interface {
	GetSongFromSoundcloudURL(url string) (string, bool)
	SearchSong(name string) soundcloud.SongModel
}

type MongoClient interface {
	NewPlayList(playlistName string)
	AddSong(playlistName string, song soundcloud.SongModel) (pl mongo.Playlist, ok bool)
	GetPlaylist(playlistName string) (pl mongo.Playlist, ok bool)
}

type Bot struct {
	sc         SoundCloudClient
	mongo      MongoClient
	Connection *discordgo.VoiceConnection
}

//creates a new bot
func NewBot() *Bot {
	guildsInfo = make(map[string](chan bool))
	songsPlaying = make(chan bool)
	next = make(chan bool, 1)
	counter = make(chan int)
	return &Bot{
		Connection: &discordgo.VoiceConnection{},
		sc:         &*soundcloud.NewSoundcloud(),
		mongo:      &*mongo.NewMongo(),
	}
}

func (b *Bot) Start(Token string) {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating discord session,", err)
		return
	}

	dg.AddHandler(b.messageCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	err = dg.Open()

	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("bot is now running. press ctrl - c to exit.")

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!gopher" {
		b.SendMessage(s, m.ChannelID, "Title much very", "", "https://cdn.discordapp.com/attachments/468414088850964480/941890121425383464/sexy_gopher_kills_crab.png", "")
		return
	}

	if m.Content == "!disconnect" {
		b.LeaveUserVoiceChannel(m.ChannelID, m.GuildID, s)
		return
	}

	if m.Content == "!join" {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				b.JoinUserVoiceChannel(s, g.ID, vs.ChannelID)
				return
			}
		}

		b.SendMessage(s, m.ChannelID, JOIN_CHANNEL_FIRST, "", "", "")
		return
	}

	if strings.Contains(m.Content, "!stop") {
		fmt.Println(m.GuildID)
		b.StopSong(m.GuildID)
		return
	}

	if strings.Contains(m.Content, "!next") {
		b.NextSong(m.GuildID)
		return
	}

	if strings.Contains(m.Content, "!stop all") {
		fmt.Println(m.GuildID)
		songs, _ := b.GetPlaylist(playlistPlaying)

		b.StopAll(m.GuildID, len(songs.Songs))
		return
	}

	if strings.Contains(m.Content, "!playlist create ") {
		b.NewPlayList(m.Content)
		return
	}

	if strings.Contains(m.Content, "!add song ") {
		b.AddSong(m.Content)
		return
	}

	if strings.Contains(m.Content, "!playlist play -p") {
		if b.Connection != nil {
			b.PlayAll(m.Content, m.ChannelID, m.GuildID, s, false)
		}
		return
	}

	if strings.Contains(m.Content, "!playlist shuffle -p ") {
		if b.Connection != nil {
			b.PlayAll(m.Content, m.ChannelID, m.GuildID, s, true)
		}
		return
	}

	if strings.Contains(m.Content, "!play ") {
		if b.Connection != nil {
			if !b.Connection.Ready {
				b.SendMessage(s, m.ChannelID, JOIN_CHANNEL_FIRST_BOT, "", "", "")
				return
			}
			b.PlaySingular(s, m.GuildID, m.ChannelID, m.Content)
			<-songsPlaying
			return
		}

		b.SendMessage(s, m.ChannelID, SONG_NOT_FOUND, "", "", "")
		return
	}

	if strings.Contains(m.Content, "!help") {
		b.ShowHelp(m.ChannelID, m.GuildID, s)
	}
}

func (b *Bot) play(guildId string, songCDN string) {
	defer func() {
		next <- true
		songsPlaying <- true
		counter <- 1
	}()

	dgvoice.PlayAudioFile(b.Connection, songCDN, guildsInfo[guildId])
}

func (b *Bot) SendMessage(s *discordgo.Session, channelId, title, desc, imgUrl, url string) {
	embeded := embed.NewEmbed().SetFooter(FOOTER_TEXT).SetTitle(title).SetDescription(desc).SetImage(imgUrl).SetURL(url).MessageEmbed
	s.ChannelMessageSendEmbed(channelId, embeded)
}

//this method plays a song once a stream url is passed to it
func (b *Bot) PlaySong(s *discordgo.Session, guildId string, channelId string, songCDN string) bool {
	fmt.Println(guildId)
	fmt.Println("playing:")
	fmt.Println(songCDN)
	guildsInfo[guildId] = make(chan bool, 2)
	if !b.Connection.Ready {
		return false
	}

	go b.play(guildId, songCDN)

	return true
}

//play a random song, gets the first match matching the searched command
func (b *Bot) PlaySingular(s *discordgo.Session, guildId string, channelId, command string) bool {
	songCDN, song, ok := b.GetSong(s, guildId, channelId, command)
	if !ok {
		return false
	}
	b.StopSong(guildId)
	if b.PlaySong(s, guildId, channelId, songCDN) {
		embedTitle := fmt.Sprintf("%s %s", SONG_PLAYING, song.Title)
		b.SendMessage(s, channelId, embedTitle, "", song.ArtworkURL, song.PermalinkURL)
	}

	<-next
	return true
}

func (b *Bot) GetSong(s *discordgo.Session, guildId string, channelId string, command string) (string, soundcloud.SongModel, bool) {
	b.StopSong(guildId)
	songName, ok := utils.GetSongTitleFromCommand(command)
	if !ok {
		return "", soundcloud.SongModel{}, false
	}
	fmt.Println(songName)
	song := b.sc.SearchSong(songName)

	if song.IsStructureEmpty() {
		b.SendMessage(s, channelId, SONG_NOT_FOUND, "", "", "")
		return "", soundcloud.SongModel{}, false
	}

	if len(song.Media.Transcodings) == 0 {
		return "", soundcloud.SongModel{}, false
	}

	var index int = 1
	if len(song.Media.Transcodings) == 1 {
		index = 0
	}

	songCDN, ok := b.sc.GetSongFromSoundcloudURL(song.Media.Transcodings[index].URL)
	if !ok {
		b.SendMessage(s, channelId, SONG_NOT_FOUND, "", "", "")
	}
	return songCDN, song, true
}

//this method stops a song when a guild id is passed to it
func (b *Bot) StopSong(guildID string) {
	go func() {
		guildsInfo[guildID] <- true

		//	next <- true
		//	b.Connection.Speaking(false)
	}()
}

//this method stops a song when a guild id is passed to it
func (b *Bot) NextSong(guildID string) {
	go func() {
		guildsInfo[guildID] <- true
		//next <- true
	}()
}

//will work for version two
//this method goes and reads from channel all the songs that will
//be played from a playlist in different goroutines and stops them
func (b *Bot) StopAll(guildID string, playlistSize int) {
	passed := 0
	for elem := range counter {
		fmt.Print(elem)
		passed++
	}
	for i := 0; i < playlistSize-passed+1; i++ {
		go func() {
			guildsInfo[guildID] <- true
			next <- true
			b.Connection.Speaking(false)
		}()
	}
}

//this method gets the discord session, guild id and channel id to join the voice channel
func (b *Bot) JoinUserVoiceChannel(s *discordgo.Session, guildID, channelID string) *discordgo.VoiceConnection {

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil
	}

	b.Connection = vc
	b.SendMessage(s, channelID, BOT_JOINED, "", "", "")
	return vc
}

//this method makes the bot disconnect from the voice channel
func (b *Bot) LeaveUserVoiceChannel(channelID string, guildId string, s *discordgo.Session) {
	b.Connection.Disconnect()

	b.SendMessage(s, channelID, BOT_LEFT, "", "", "")
}

//play all songs from a play list with a given name, if shuffle is true, songs are played at a random order
func (b *Bot) PlayAll(command string, channelID string, guildId string, s *discordgo.Session, shuffle bool) {
	playlistName, ok := utils.GetPlalistNameFromPlayCommand(command)
	if !ok {
		b.SendMessage(s, channelID, INVALID_COMMAND, "", "", "")
		return
	}
	playlistPlaying = playlistName
	playlist, ok := b.GetPlaylist(playlistName)
	if !ok {
		b.SendMessage(s, channelID, PLAYLIST_NOT_FOUND, "", "", "")
		return
	}

	if !b.Connection.Ready {
		b.SendMessage(s, channelID, SHOULD_CONNECT_TO_VOICE_CHANNEL, "", "", "")
		return
	}

	if shuffle {
		playlist.Songs = utils.Shuffle(playlist.Songs)
	}

	b.StopSong(guildId)

	for _, item := range playlist.Songs {
		songCDN, ok := b.sc.GetSongFromSoundcloudURL(item)

		if !ok {
			b.SendMessage(s, channelID, SONG_NOT_FOUND, "", "", "")
			return
		}

		if b.PlaySong(s, guildId, channelID, songCDN) {
			b.SendMessage(s, channelID, SONG_PLAYING_FROM_PLAYLIST, "", "", "")
		}

		<-next
	}
}

//create a new playlist
func (b *Bot) NewPlayList(playlistName string) bool {
	name, ok := utils.GetPlalistNameFromCreateCommand(playlistName)
	if !ok {
		return false
	}
	b.mongo.NewPlayList(name)
	return true
}

//add a song to a playlist by name
func (b *Bot) AddSong(command string) bool {
	playlistName, songName, okay := utils.GetPlalistNameAndTrackFromCommand(command)
	if !okay {
		return okay
	}
	song := b.sc.SearchSong(songName)
	playlist, ok := b.mongo.AddSong(playlistName, song)
	fmt.Printf("added song to %s", playlist)
	return ok
}

//retrieve a playlist by name, if none is found with that name, ok is false
func (b *Bot) GetPlaylist(playlistName string) (pl mongo.Playlist, ok bool) {
	return b.mongo.GetPlaylist(playlistName)
}

func (b *Bot) ShowHelp(channelID string, guildId string, s *discordgo.Session) {
	desc := "🔗 `!join` - joins a voice channel if there are users in it\n\n 🔌`!disconnect` - disconnects from a voice channel \n\n `▶️!play <song title>`  - plays a track from Soundcloud matching the title/author\n\n 🛑`!stop ` - stops playing a track \n\n 🆕`!playlist create <playlist name>` - creates an empty playlist with the given name\n\n ➕`!add song -p <playlist name> -t <track title>` - adds a song to a playlist \n\n ▶️`!playlist play -p <playlist name>` - plays a playlist normally \n\n `▶️!playlist shuffle -p <playlist name>` -plays a playlist on shuffle\n\n ➡️`!next` - skips a song from the playlist"
	b.SendMessage(s, channelID, AVAILABLE_COMMANDS, desc, "", "")
}
