package utils

import "testing"

func TestGetSongTitleFromCommand(t *testing.T) {
	command := "!play Maccarena"
	got, ok := GetSongTitleFromCommand(command)
	if got != "Maccarena" {
		t.Errorf("GetSongTitleFromCommand(%s) = %s; want Maccarena", command, got)
	}
	if !ok {
		t.Errorf("GetSongTitleFromCommand(%s) = %s, %t; want true", command, got, ok)
	}
}

func TestGetSongTitleFromCommandInvalidParams(t *testing.T) {
	command := "!playsadsad Maccarena"
	got, ok := GetSongTitleFromCommand(command)
	if got != "" {
		t.Errorf("GetSongTitleFromCommand(%s) = %s; want empty string", command, got)
	}
	if ok {
		t.Errorf("GetSongTitleFromCommand(%s) = %s, %t; want false", command, got, ok)
	}
}

func TestShuffle(t *testing.T) {
	arr := []string{"one", "two", "three", "four", "five", "six"}
	got := Shuffle(arr)
	if len(arr) != len(got) {
		t.Errorf("Shuffle returned array of len %d; want array of len %d", len(got), len(arr))
	}

	areEqual := func(a []string, b []string) bool {
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	if areEqual(arr, got) {
		t.Errorf("Shuffle didn't shuffle enough")
	}
}

func TestGetPlalistNameFromPlayCommandInvalid(t *testing.T) {
	command := "!playsadsad Maccarena"
	got, ok := GetPlalistNameFromPlayCommand(command)
	if got != "" {
		t.Errorf("GetPlalistNameFromPlayCommand(%s) = %s; want empty string", command, got)
	}
	if ok {
		t.Errorf("GetPlalistNameFromPlayCommand(%s) = %s, %t; want false", command, got, ok)
	}
}

func TestGetPlaylistNameFromCreateCommandInvalid(t *testing.T) {
	command := "!playlist -sad"
	got, ok := GetPlalistNameFromCreateCommand(command)
	if got != "" {
		t.Errorf("GetPlalistNameFromCreateCommand(%s) = %s; want empty string", command, got)
	}
	if ok {
		t.Errorf("GetPlalistNameFromCreateCommand(%s) = %s, %t; want false", command, got, ok)
	}
}

func TestGetPlaylistNameFromCreateCommand(t *testing.T) {
	command := "!playlist create myPlaylist"
	got, ok := GetPlalistNameFromCreateCommand(command)
	if got != "myPlaylist" {
		t.Errorf("GetPlalistNameFromCreateCommand(%s) = %s; want myPlaylist", command, got)
	}
	if !ok {
		t.Errorf("GetPlalistNameFromCreateCommand(%s) = %s, %t; want true", command, got, ok)
	}
}

func TestGetPlalistNameAndTrackFromCommandInvalid(t *testing.T) {
	command := "!notadd notsong "
	plName, track, ok := GetPlalistNameAndTrackFromCommand(command)
	if ok {
		t.Errorf("GetPlalistNameAndTrackFromCommand(%s)= _,_, %t; want false", command, ok)
	}
	if plName != "" || track != "" {
		t.Errorf("GetPlalistNameAndTrackFromCommand(%s) = %s, %s; want empty string", command, plName, track)
	}
}
