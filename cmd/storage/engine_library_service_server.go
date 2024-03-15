package main

import (
	"context"
	_ "embed"
	"fmt"
	"math"
	"strconv"

	"github.com/icedream/go-stagelinq/eaas/proto/enginelibrary"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ enginelibrary.EngineLibraryServiceServer = &EngineLibraryServiceServer{}

type EngineLibraryServiceServer struct {
	enginelibrary.UnimplementedEngineLibraryServiceServer
}

// EventStream implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) EventStream(ctx context.Context, req *enginelibrary.EventStreamRequest) (*enginelibrary.EventStreamResponse, error) {
	return &enginelibrary.EventStreamResponse{
		Event: []*enginelibrary.Event{},
	}, nil
}

// GetCredentials implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetCredentials(context.Context, *enginelibrary.GetCredentialsRequest) (*enginelibrary.GetCredentialsResponse, error) {
	panic("unimplemented")
}

// GetHistoryPlayedTracks implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetHistoryPlayedTracks(context.Context, *enginelibrary.GetHistoryPlayedTracksRequest) (*enginelibrary.GetHistoryPlayedTracksResponse, error) {
	return &enginelibrary.GetHistoryPlayedTracksResponse{
		Tracks: []*enginelibrary.HistoryPlayedTrack{},
	}, nil
}

// GetHistorySessions implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetHistorySessions(context.Context, *enginelibrary.GetHistorySessionsRequest) (*enginelibrary.GetHistorySessionsResponse, error) {
	return &enginelibrary.GetHistorySessionsResponse{
		Sessions: []*enginelibrary.HistorySession{},
	}, nil
}

// GetLibraries implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetLibraries(ctx context.Context, req *enginelibrary.GetLibrariesRequest) (*enginelibrary.GetLibrariesResponse, error) {
	return &enginelibrary.GetLibrariesResponse{
		Libraries: []*enginelibrary.Library{
			{
				Id:    &demoLibrary,
				Title: &demoLibraryName,
			},
		},
	}, nil
}

// GetLibrary implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetLibrary(ctx context.Context, req *enginelibrary.GetLibraryRequest) (*enginelibrary.GetLibraryResponse, error) {
	switch req.GetLibraryId() {
	case demoLibrary:
		return &enginelibrary.GetLibraryResponse{
			Playlists: []*enginelibrary.PlaylistMetadata{
				{
					Id:         &demoPlaylist,
					Title:      &demoPlaylistName,
					TrackCount: &demoPlaylistTrackCount,
					Playlists:  []*enginelibrary.PlaylistMetadata{},
					ListType:   enginelibrary.ListType_LIST_TYPE_PLAY.Enum(),
				},
			},
		}, nil
	default:
		return nil, status.Error(codes.NotFound, "library not found")
	}
}

// GetSearchFilters implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetSearchFilters(ctx context.Context, req *enginelibrary.GetSearchFiltersRequest) (*enginelibrary.GetSearchFiltersResponse, error) {
	resp := &enginelibrary.GetSearchFiltersResponse{
		SearchFilters: &enginelibrary.SearchFilterOptions{},
	}

	if req.LibraryId != &demoLibrary {
		return resp, nil
	}

	if demoTrackMetadata.Artist != nil {
		resp.SearchFilters.Artists = append(resp.SearchFilters.Artists, &enginelibrary.SearchFilterValue{
			Value: demoTrackMetadata.Artist,
		})
	}

	if demoTrackMetadata.Album != nil {
		resp.SearchFilters.Albums = append(resp.SearchFilters.Albums, &enginelibrary.SearchFilterValue{
			Value: demoTrackMetadata.Album,
		})
	}

	if demoTrackMetadata.Bpm != nil {
		s := fmt.Sprint(*demoTrackMetadata.Bpm)
		resp.SearchFilters.Bpms = append(resp.SearchFilters.Bpms, &enginelibrary.SearchFilterValue{
			Value: &s,
		})
	}

	if demoTrackMetadata.Genre != nil {
		resp.SearchFilters.Genres = append(resp.SearchFilters.Genres, &enginelibrary.SearchFilterValue{
			Value: demoTrackMetadata.Genre,
		})
	}

	if demoTrackMetadata.Key != nil {
		resp.SearchFilters.Keys = append(resp.SearchFilters.Keys, &enginelibrary.SearchFilterValue{
			Value: demoTrackMetadata.Key,
		})
	}

	return resp, nil
}

func generateDemoTrackMetadata(trackID string) *enginelibrary.TrackMetadata {
	metadata := enginelibrary.TrackMetadata{
		Title:         demoTrackMetadata.Title,
		Bpm:           demoTrackMetadata.Bpm,
		DateAdded:     demoTrackMetadata.DateAdded,
		Artist:        demoTrackMetadata.Artist,
		Album:         demoTrackMetadata.Album,
		Key:           demoTrackMetadata.Key,
		Rating:        demoTrackMetadata.Rating,
		Year:          demoTrackMetadata.Year,
		Genre:         demoTrackMetadata.Genre,
		Comment:       demoTrackMetadata.Comment,
		Label:         demoTrackMetadata.Label,
		LengthSeconds: demoTrackMetadata.LengthSeconds,
		Composer:      demoTrackMetadata.Composer,
		Remixer:       demoTrackMetadata.Remixer,
		Id:            &trackID,
	}
	metadata.Id = &trackID
	return &metadata
}

// GetTrack implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetTrack(ctx context.Context, req *enginelibrary.GetTrackRequest) (*enginelibrary.GetTrackResponse, error) {
	if req.GetLibraryId() != demoLibrary && req.LibraryId != nil {
		return nil, status.Error(codes.NotFound, "library not found")
	}

	for _, trackID := range demoTrackIDs {
		if trackID == req.GetTrackId() {
			return &enginelibrary.GetTrackResponse{
				Blob: &enginelibrary.TrackBlob{
					Type: &enginelibrary.TrackBlob_Url{
						Url: &enginelibrary.TrackBlobUrl{
							Url:      &demoTrackURL,
							FileSize: &demoTrackLength,
						},
					},
				},
				Metadata:        generateDemoTrackMetadata(trackID),
				PerformanceData: nil, // TODO
			}, nil
		}
	}

	return nil, status.Error(codes.NotFound, "track not found")
}

// GetTracks implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) GetTracks(ctx context.Context, req *enginelibrary.GetTracksRequest) (*enginelibrary.GetTracksResponse, error) {
	switch req.GetLibraryId() {
	case "", demoLibrary:
		resp := &enginelibrary.GetTracksResponse{
			Tracks: []*enginelibrary.ListTrack{},
		}
		for _, trackID := range demoTrackIDs {
			resp.Tracks = append(resp.Tracks, &enginelibrary.ListTrack{
				Metadata:       generateDemoTrackMetadata(trackID),
				PreviewArtwork: demoTrackArtwork,
			})
		}
		return &enginelibrary.GetTracksResponse{
			Tracks: []*enginelibrary.ListTrack{},
		}, nil
	default:
		return &enginelibrary.GetTracksResponse{
			Tracks: []*enginelibrary.ListTrack{},
		}, nil
	}
}

// PutEvents implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) PutEvents(ctx context.Context, req *enginelibrary.PutEventsRequest) (*enginelibrary.PutEventsResponse, error) {
	return &enginelibrary.PutEventsResponse{}, nil
}

// SearchTracks implements enginelibrary.EngineLibraryServiceServer.
func (e *EngineLibraryServiceServer) SearchTracks(ctx context.Context, req *enginelibrary.SearchTracksRequest) (*enginelibrary.SearchTracksResponse, error) {
	resp := &enginelibrary.SearchTracksResponse{
		Tracks: []*enginelibrary.ListTrack{},
	}
trackLoop:
	for _, trackID := range demoTrackIDs {
		metadata := generateDemoTrackMetadata(trackID)
		for _, filter := range req.Filters {
			switch filter.GetField() {
			case enginelibrary.SearchFilterField_SEARCH_FILTER_FIELD_ALBUM:
				if metadata.Album == nil {
					continue trackLoop
				}
				for _, value := range filter.GetValue() {
					if !fuzzy.MatchFold(value, *metadata.Album) {
						continue trackLoop
					}
				}
			case enginelibrary.SearchFilterField_SEARCH_FILTER_FIELD_ARTIST:
				if metadata.Artist == nil {
					continue trackLoop
				}
				for _, value := range filter.GetValue() {
					if !fuzzy.MatchFold(value, *metadata.Artist) {
						continue trackLoop
					}
				}
			case enginelibrary.SearchFilterField_SEARCH_FILTER_FIELD_BPM:
				if metadata.Bpm == nil {
					continue trackLoop
				}
				for _, value := range filter.GetValue() {
					f, err := strconv.ParseFloat(fmt.Sprint(value), 64)
					if err != nil {
						continue trackLoop
					}
					if math.Abs(f-*metadata.Bpm) > 3 {
						continue trackLoop
					}
				}
			case enginelibrary.SearchFilterField_SEARCH_FILTER_FIELD_GENRE:
				if metadata.Genre == nil {
					continue trackLoop
				}
				for _, value := range filter.GetValue() {
					if !fuzzy.MatchFold(value, *metadata.Genre) {
						continue trackLoop
					}
				}
			case enginelibrary.SearchFilterField_SEARCH_FILTER_FIELD_KEY:
				if metadata.Key == nil {
					continue trackLoop
				}
				for _, value := range filter.GetValue() {
					if !fuzzy.MatchFold(value, *metadata.Key) {
						continue trackLoop
					}
				}
			default:
			}
		}
		resp.Tracks = append(resp.Tracks, &enginelibrary.ListTrack{
			Metadata:       metadata,
			PreviewArtwork: demoTrackArtwork,
		})
	}
	return resp, nil
}
