package models

import (
	"errors"
	"sync"
)

func NewTrackStorage(mu *sync.Mutex) *TrackStorage {
	storage := &TrackStorage{
		TrackStorage: make(map[uint]*Track),
		Mutex:        mu,
	}
	track := Track{
		Id:       12345,
		Name:     "Symphony 40",
		Artist:   "Mozart",
		Duration: 37,
		Image:    "https://steemitimages.com/DQmdM5W5dBi5Kg4zwpeMC5Ty2fZEKig1kQ1tXQJUxdQP7Ph/John.jpg",
		Link:     "https://d4.hotplayer.ru/download/4219e191e9078c11d1e7825344da42b8/287373405_456239023/12a221000968e-51fb57871d78-d722a1dfefb/HIS%20NAME%20IS%20-%20JOHN%20CENA%20%232.mp3",
	}
	storage.TrackStorage[12345] = &track
	return storage
}

type TrackStorage struct {
	TrackStorage map[uint]*Track
	Mutex        *sync.Mutex
}

type Track struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration uint   `json:"duration"`
	Image    string `json:"image"`
	Link     string `json:"link"`
}

func (us *TrackStorage) GetFullTrackInfo(id uint) (Track, error) {
	us.Mutex.Lock()
	defer us.Mutex.Unlock()
	if track, ok := us.TrackStorage[id]; !ok {
		return Track{}, errors.New("track with this id does not exists: " + string(id))
	} else {
		return *track, nil
	}
}
