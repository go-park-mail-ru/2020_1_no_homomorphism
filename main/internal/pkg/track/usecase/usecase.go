package usecase

import (
	"fmt"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/models"
	"github.com/2020_1_no_homomorphism/no_homo_main/internal/pkg/track"
	"strconv"
)

type TrackUseCase struct {
	Repository track.Repository
}

func (uc TrackUseCase) GetTrackById(id string) (models.Track, error) {
	return uc.Repository.GetTrackById(id)
}

func (uc TrackUseCase) RateTrack(uID, tID string) error {
	return uc.Repository.RateTrack(uID, tID)
}
func (uc TrackUseCase) GetAllTracks() ([]models.Track, error)  {
	return uc.Repository.GetAllTracks()
}

func (uc TrackUseCase) GetBoundedTracksByArtistId(id string, start, end uint64, uID string) ([]models.Track, error) {
	dbTracks, err := uc.Repository.GetBoundedTracksByArtistId(id, start, end)
	if err != nil {
		return nil, err
	}
	if uID != "" {
		if err = uc.setLikes(dbTracks, uID); err != nil {
			return nil, err
		}
	}
	return dbTracks, nil
}

func (uc TrackUseCase) GetBoundedTracksByAlbumId(id string, start, end uint64, uID string) ([]models.Track, error) {
	dbTracks, err := uc.Repository.GetBoundedTracksByAlbumId(id, start, end)
	if err != nil {
		return nil, err
	}
	if uID != "" {
		if err = uc.setLikes(dbTracks, uID); err != nil {
			return nil, err
		}
	}
	fmt.Println(dbTracks)
	return dbTracks, nil
}

func (uc TrackUseCase) GetBoundedTracksByPlaylistId(plId string, start, end uint64, uID string) ([]models.Track, error) {
	dbTracks, err := uc.Repository.GetBoundedTracksByPlaylistId(plId, start, end)
	if err != nil {
		return nil, err
	}
	if uID != "" {
		if err = uc.setLikes(dbTracks, uID); err != nil {
			return nil, err
		}
	}
	return dbTracks, nil
}

func (uc TrackUseCase) Search(text string, count uint) ([]models.TrackSearch, error) {
	return uc.Repository.Search(text, count)
}

func (uc TrackUseCase) GetUserTracks(uID string) ([]models.Track, error) {
	return uc.Repository.GetUserTracks(uID)
}

func (uc TrackUseCase) setLikes(tracks []models.Track, uID string) error {
	trackIDs, err := uc.Repository.GetUserLikedTracksIDs(uID)
	if err != nil {
		return err
	}
	intersection := make(map[string]int)

	for _, elem := range trackIDs {
		intersection[strconv.FormatInt(elem, 10)] = 1
	}

	for i, elem := range tracks {
		if _, ok := intersection[elem.Id]; ok {
			tracks[i].IsLiked = true
		}
	}
	return nil
}

func (uc TrackUseCase) IsLikedByUser(uID string, tID string) (bool, error) {
	if uID == "" || tID == "" {
		return false, nil
	}
	return uc.Repository.IsLikedByUser(uID, tID)
}