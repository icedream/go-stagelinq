package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"github.com/icedream/go-stagelinq/eaas"
	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Imports needed for image resizing (see commented out code for it)
// import (
// 	"image"
// 	_ "image/jpeg"
// 	_ "image/png"
// )

var (
	demoTrackFileName             = "Icedream - Whiplash (Radio Edit).m4a"
	demoLibrary                   = "12eceaa2-f81a-4b63-b196-94648a3bdd95"
	demoLibraryName               = "Demo Library"
	demoPlaylist                  = "55ab0c7c-6c35-429a-81d0-25b039a34a9f"
	demoPlaylistName              = "Demo Playlist"
	demoPlaylistTrackCount uint32 = 1
	demoTrackIDs                  = []string{
		"1 " + demoLibrary,
	}
	// HACK - imitating original Engine DJ software behavior by using Windows paths
	demoTrackURL = filepath.Join("C:", "demo", demoTrackFileName)
	// HACK - imitating original Engine DJ software behavior by adding brackets.
	demoTrackURLGRPC  = fmt.Sprintf("<%s>", filepath.ToSlash(demoTrackURL))
	demoTrackLength   = uint32(len(demoTrackBytes))
	demoTrackMetadata enginelibrary.TrackMetadata
	demoTrackArtwork  []byte
	demoToken         eaas.Token = eaas.Token{
		0x5e, 0xff, 0xae, 0x59, 0x12, 0x88, 0x29, 0x30,
		0xde, 0xad, 0xc0, 0xde, 0xc0, 0xff, 0xee, 0x00,
	}
)

//go:embed "Icedream - Whiplash (Radio Edit).m4a"
var demoTrackBytes []byte

//go:embed "Icedream - Whiplash (Radio Edit).m4a.beatgrid"
var demoBeatGrid []byte

//go:embed "Icedream - Whiplash (Radio Edit).m4a.waveform"
var demoOverviewWaveform []byte

var demoTrackPreviewArtwork []byte

func init() {
	if len(demoTrackIDs) == 0 {
		for i := 0; i < int(demoPlaylistTrackCount); i++ {
			demoTrackIDs = append(demoTrackIDs, uuid.New().String())
		}
	}

	demoTrackMetadata.DateAdded = timestamppb.Now()
	if metadata, err := tag.ReadFrom(bytes.NewReader(demoTrackBytes)); err == nil {
		if metadata.Picture() != nil {
			demoTrackArtwork = metadata.Picture().Data
			demoTrackPreviewArtwork = demoTrackArtwork
			// // If you wanna be nice to the hardware, you can have the server
			// // shrink down the artwork. I don't think even the original Engine
			// // DJ software does that though.
			// img, _, err := image.Decode(bytes.NewReader(demoTrackArtwork))
			// if err == nil {
			// 	img = resize.Resize(240, 240, img, resize.Lanczos2)
			// }
			// var b bytes.Buffer
			// if err := jpeg.Encode(&b, img, &jpeg.Options{Quality: 70}); err == nil {
			// 	demoTrackPreviewArtwork = b.Bytes()
			// }
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
			s = strings.Trim(s, "\x00 ")
			demoTrackMetadata.Key = &s
		}
		if v, ok := metadata.Raw()["LABEL"]; ok {
			s := fmt.Sprint(v)
			s = strings.Trim(s, "\x00 ")
			demoTrackMetadata.Label = &s
		}
		if v, ok := metadata.Raw()["REMIXER"]; ok {
			s := fmt.Sprint(v)
			s = strings.Trim(s, "\x00 ")
			demoTrackMetadata.Remixer = &s
		}
		if v := uint32(metadata.Year()); v > 0 {
			demoTrackMetadata.Year = &v
		}
	}
}
