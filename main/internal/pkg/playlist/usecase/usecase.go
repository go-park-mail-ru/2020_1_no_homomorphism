package usecase

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/playlist"
)

type PlaylistUseCase struct {
	PlRepository playlist.Repository
}

func (uc PlaylistUseCase) GetUserPlaylists(id string) ([]models.Playlist, error) {
	return uc.PlRepository.GetUserPlaylists(id)
}

func (uc PlaylistUseCase) GetPlaylistById(id string) (models.Playlist, error) {
	return uc.PlRepository.GetPlaylistById(id)
}

func (uc PlaylistUseCase) CreatePlaylist(name string, uID string) (plID string, err error) {
	return uc.PlRepository.CreatePlaylist(name, uID)
}

func (uc PlaylistUseCase) AddTrackToPlaylist(plTracks models.PlaylistTracks) error {
	return uc.PlRepository.AddTrackToPlaylist(plTracks)
}

func (uc PlaylistUseCase) GetUserPlaylistsIdByTrack(userID, trackID string) ([]string, error) {
	return uc.PlRepository.GetUserPlaylistsIdByTrack(userID, trackID)
}

func (uc PlaylistUseCase) DeletePlaylist(plID string) error {
	return uc.PlRepository.DeletePlaylist(plID)
}

func (uc PlaylistUseCase) DeleteTrackFromPlaylist(plID, trackID string) error {
	return uc.PlRepository.DeleteTrackFromPlaylist(plID, trackID)
}

func (uc PlaylistUseCase) ChangePrivacy(plID string) error {
	return uc.PlRepository.ChangePrivacy(plID)
}

func (uc PlaylistUseCase) AddSharedPlaylist(plID string, uID string) (string, error) {
	pl, err := uc.PlRepository.GetPlaylistById(plID)
	if err != nil {
		return "", fmt.Errorf("cant get playlist: %v", err)
	}
	newPl, err := uc.PlRepository.CreatePlaylist(pl.Name, uID)
	if err != nil {
		return "", fmt.Errorf("cant create playlist: %v", err)
	}
	tracks, err := uc.PlRepository.GetAllPlaylistTracks(plID)

	for _, elem := range tracks {
		plTracks := models.PlaylistTracks{
			PlaylistID: newPl,
			TrackID:    elem.TrackID,
			Index:      elem.Index,
			Image:      elem.Image,
		}
		err := uc.PlRepository.AddTrackToPlaylist(plTracks)
		if err != nil {
			return "", fmt.Errorf("failed to add track to playlist: %v", err)
		}
	}

	return newPl, nil
}

func (uc PlaylistUseCase) CheckAccessToPlaylist(userId string, playlistId string, isStrict bool) (bool, error) {
	pl, err := uc.PlRepository.GetPlaylistById(playlistId)
	if err != nil {
		return false, fmt.Errorf("cant get playlist: %v", err)
	}

	if isStrict {
		return pl.UserId == userId, nil
	}
	return pl.UserId == userId || !pl.Private, nil
}
