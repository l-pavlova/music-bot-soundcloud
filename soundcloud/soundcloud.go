package soundcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

//this is taken from a soundcloud standart user account because you can't register your app since 2019.
//auth token seems to not expire, works for all requests needed by the bot
var auth string = "OAuth 2-293253-1082025049-TPPFIYYYKTyes"

type PublisherMetadata struct {
	ID            int    `json:"id"`
	Urn           string `json:"urn"`
	Artist        string `json:"artist"`
	ContainsMusic bool   `json:"contains_music"`
}

//generated with https://mholt.github.io/json-to-go/
//only the essential fields that might be needed are left
type SongModel struct {
	ArtworkURL        string            `json:"artwork_url"`
	Caption           interface{}       `json:"caption"`
	Title             interface{}       `json:"title"`
	CreatedAt         time.Time         `json:"created_at"`
	Description       string            `json:"description"`
	Duration          int               `json:"duration"`
	Genre             string            `json:"genre"`
	ID                int               `json:"id"`
	LabelName         interface{}       `json:"label_name"`
	LikesCount        int               `json:"likes_count"`
	PublisherMetadata PublisherMetadata `json:"publisher_metadata"`
	State             string            `json:"state"`
	Streamable        bool              `json:"streamable"`
	TrackFormat       string            `json:"track_format"`
	URI               string            `json:"uri"`
	PermalinkURL      string            `json:"permalink_url"`
	Media             struct {
		Transcodings []struct {
			URL      string `json:"url"`
			Preset   string `json:"preset"`
			Duration int    `json:"duration"`
			Snipped  bool   `json:"snipped"`
			Format   struct {
				Protocol string `json:"protocol"`
				MimeType string `json:"mime_type"`
			} `json:"format"`
			Quality string `json:"quality"`
		} `json:"transcodings"`
	} `json:"media"`
	Policy string `json:"policy"`
}

type Songs struct {
	Songs []SongModel `json:"collection"`
}

type SongCDN struct {
	SongURL string `json:"url"`
}

type SoundCloud struct {
}

func NewSoundcloud() *SoundCloud {
	return &SoundCloud{}
}

//this method searches a song by its name in the sound cloud api
func (sc *SoundCloud) SearchSong(name string) SongModel {
	query := fmt.Sprintf("https://api-v2.soundcloud.com/search?q=`%s`&sc_a_id=f55635101bfdf1e8418a36ef0ee8e86f23d9f257&variant_ids=2451&facet=model&user_id=565035-794848-92508-940751&client_id=BmI0Zgypr3dPccFBK9QLjkCpCgvowlzQ&limit=20&offset=0&linked_partitioning=1&app_version=1643966166&app_locale=en", name)
	soundcloudClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	fmt.Println(query)

	req, _ := http.NewRequest("GET", query, nil)
	req.Header.Set("Authorization", auth)

	res, err := soundcloudClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(err)
	}

	var songs Songs

	jsonErr := json.Unmarshal(body, &songs)
	if jsonErr != nil {

		fmt.Println(err)
	}
	if len(songs.Songs) == 0 {
		return SongModel{}
	}

	ind := 0
	for len(songs.Songs[ind].Media.Transcodings) < 2 {
		ind++
		if ind >= len(songs.Songs) {
			return SongModel{}
		}
	}

	return songs.Songs[ind] //returns the progressive stream link
}

//this method gets the song stream url from the sound cloud api, passing the auth param
func (sc *SoundCloud) GetSongFromSoundcloudURL(url string) (string, bool) {
	soundcloudClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", auth)
	res, getErr := soundcloudClient.Do(req)
	if getErr != nil {
		fmt.Println(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
		fmt.Println(res.Body)
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			fmt.Println(err)
		}

		var song SongCDN

		jsonErr := json.Unmarshal(body, &song)
		if jsonErr != nil {
			fmt.Println(err)
		}

		return song.SongURL, true
	}

	return "", false
}

//util function to check if a structure is empty
func (x SongModel) IsStructureEmpty() bool {
	return reflect.DeepEqual(x, SongModel{})
}
