package usecase

import (
	"fmt"
	"no_homomorphism/internal/pkg/models"
	"no_homomorphism/internal/pkg/playlist"
	"no_homomorphism/internal/pkg/track"
	"strconv"
)

type PlaylistUseCase struct {
	PlRepository    playlist.Repository
	TrackRepository track.Repository
}

func (uc PlaylistUseCase) GetUserPlaylists(id string) (models.UserPlaylists, error) {
	plId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return models.UserPlaylists{}, fmt.Errorf("failed to convert id: %e", err)
	}
	playlists, err := uc.PlRepository.GetUserPlaylists(plId)
	return models.UserPlaylists{
		Count:     len(playlists),
		Playlists: playlists,
	}, err
}

func (uc PlaylistUseCase) GetPlaylistWithTracks(id string) (models.PlaylistTracks, error) {
	plId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return models.PlaylistTracks{}, fmt.Errorf("failed to convert id: %s", err)
	}
	dbPlaylist, err := uc.PlRepository.GetPlaylistById(plId)
	if err != nil {
		return models.PlaylistTracks{}, fmt.Errorf("failed to get dbPlaylist: %s", err)
	}

	tracks, err := uc.TrackRepository.GetPlaylistTracks(plId)
	if err != nil {
		return models.PlaylistTracks{}, fmt.Errorf("failed to get dbPlaylist' tracks: %s", err)
	}

	return models.PlaylistTracks{
		Playlist: dbPlaylist,
		Count:    len(tracks),
		Tracks:   tracks,
	}, nil
}
