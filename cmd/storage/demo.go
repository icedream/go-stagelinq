package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	demoTrackFileName             = "Icedream - Whiplash (Radio Edit).flac"
	demoLibrary                   = "12eceaa2-f81a-4b63-b196-94648a3bdd95"
	demoLibraryName               = "Demo Library"
	demoPlaylist                  = "55ab0c7c-6c35-429a-81d0-25b039a34a9f"
	demoPlaylistName              = "Demo Playlist"
	demoPlaylistTrackCount uint32 = 1
	demoTrackIDs           []string
	demoTrackURL           = "/demo/" + demoTrackFileName
	demoTrackLength        = uint32(len(demoTrackBytes))
	demoTrackMetadata      enginelibrary.TrackMetadata
	demoTrackArtwork       []byte
)

//go:embed "Icedream - Whiplash (Radio Edit).m4a"
var demoTrackBytes []byte

func init() {
	for i := 0; i < int(demoPlaylistTrackCount); i++ {
		demoTrackIDs = append(demoTrackIDs, uuid.New().String())
	}

	demoTrackMetadata.DateAdded = timestamppb.Now()
	if metadata, err := tag.ReadFrom(bytes.NewReader(demoTrackBytes)); err == nil {
		if metadata.Picture() != nil {
			demoTrackArtwork = metadata.Picture().Data
		}
		if v := metadata.Artist(); len(v) > 0 {
			demoTrackMetadata.Artist = &v
		}
		if v := metadata.Title(); len(v) > 0 {
			demoTrackMetadata.Title = &v
		}
		if v := metadata.Album(); len(v) > 0 {
			demoTrackMetadata.Album = &v
		}
		if v, ok := metadata.Raw()["BPM"]; ok {
			if f, err := strconv.ParseFloat(fmt.Sprint(v), 64); err == nil {
				demoTrackMetadata.Bpm = &f
			}
		}
		if v := metadata.Comment(); len(v) > 0 {
			demoTrackMetadata.Comment = &v
		}
		if v := metadata.Composer(); len(v) > 0 {
			demoTrackMetadata.Composer = &v
		}
		if v := metadata.Genre(); len(v) > 0 {
			demoTrackMetadata.Genre = &v
		}
		if v, ok := metadata.Raw()["KEY"]; ok {
			s := fmt.Sprint(v)
			demoTrackMetadata.Key = &s
		}
		if v, ok := metadata.Raw()["LABEL"]; ok {
			s := fmt.Sprint(v)
			demoTrackMetadata.Label = &s
		}
		if v, ok := metadata.Raw()["REMIXER"]; ok {
			s := fmt.Sprint(v)
			demoTrackMetadata.Remixer = &s
		}
		if v := uint32(metadata.Year()); v > 0 {
			demoTrackMetadata.Year = &v
		}

	}
}
