package repository

import (
	"errors"

	"no_homomorphism/internal/pkg/models"
)


type TrackRepository struct {
	Repository []*models.Track
}

func NewTestRepo() *TrackRepository{
	return &TrackRepository{
		Repository: []*models.Track{
			&models.Track{
				Id:       12345,
				Name:     "Symphony 40",
				Artist:   "Mozart",
				Duration: 37,
				Image:    "https://steemitimages.com/DQmdM5W5dBi5Kg4zwpeMC5Ty2fZEKig1kQ1tXQJUxdQP7Ph/John.jpg",
				Link:     "https://d4.hotplayer.ru/download/4219e191e9078c11d1e7825344da42b8/287373405_456239023/12a221000968e-51fb57871d78-d722a1dfefb/HIS%20NAME%20IS%20-%20JOHN%20CENA%20%232.mp3",
			},
			&models.Track{
				Id:       12346,
				Name:     "New year dubstep minimix",
				Artist:   "DJ Epoxxin",
				Duration: 123,
				Image:    "static/img/new_empire_vol1.jpg",
				Link:     "https://s3-us-west-2.amazonaws.com/s.cdpn.io/9473/new_year_dubstep_minimix.ogg",
			},
			&models.Track{
				Id:       12347,
				Name:     "Thirteen Thirty Five",
				Artist:   "Dillon",
				Duration: 223,
				Image:    "static/img/vk.jpg",
				Link:     "http://beloweb.ru/audio/dillon_-_thirteen_thirtyfive_.mp3",
			},
			&models.Track{
				Id:       12348,
				Name:     "Пчела",
				Artist:   "Пчеловод",
				Duration: 170,
				Image:    "static/img/ok.png",
				Link:     "https://ns1.topzaycevs.ru/files/dl/rasa_-_Tii_pchela_ya_pchelovod.mp3",
			},
			&models.Track{
				Id:       12349,
				Name:     "Крокодил",
				Artist:   "Стас",
				Duration: 40,
				Image:    "static/img/rocket.svg",
				Link:     "http://cdn1.pesnigoo.ru/uploads/files/2018-10/jekstaz-krokodil_456242584.mp3",
			},
			&models.Track{
				Id:       12344,
				Name:     "Самый лучший эмо панк",
				Artist:   "Пошлая Молли",
				Duration: 208,
				Image:    "static/img/HU.jpeg",
				Link:     "https://ns1.topzaycevs.ru/files/dl/Poshlaya_Molli_-_Samiiy_luchshiy_emo_pank.mp3",
			},
		},
	}
}

func (tr *TrackRepository) GetTrackById(id uint) (*models.Track, error) {
	for _, r:=range tr.Repository {
		if r.Id == id {
			return r, nil
		}
	}
	return nil, errors.New("track with this id does not exist")
}
