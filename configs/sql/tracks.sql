DROP TABLE IF EXISTS tracks CASCADE;
DROP TABLE IF EXISTS albums CASCADE;
DROP TABLE IF EXISTS artists CASCADE;
DROP TABLE IF EXISTS liked_artists CASCADE;
DROP TABLE IF EXISTS liked_albums CASCADE;
DROP VIEW IF EXISTS track_with_artist CASCADE;
DROP VIEW IF EXISTS artist_tracks CASCADE;

CREATE TABLE artists -- TODO Подумать над полями
(
    ID    BIGSERIAL PRIMARY KEY,
    name  VARCHAR(50) NOT NULL,
    image VARCHAR(100) DEFAULT '/static/img/default.png'
);

CREATE TABLE albums
(
    ID       BIGSERIAL PRIMARY KEY,
    name     VARCHAR(50) NOT NULL,
    image    VARCHAR(100) DEFAULT '/static/img/default.png',
    artist_ID BIGSERIAL   NOT NULL,
    FOREIGN KEY (artist_ID) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

SELECT * FROM albums a WHERE a.artist_id = 1;

CREATE TABLE tracks -- TODO Изменить на N-M с артистом + добавить доп сущность для связи альбома и трека
(
    ID       BIGSERIAL PRIMARY KEY,
    name     VARCHAR(50) NOT NULL,
    duration INTEGER     NOT NULL,
    --image    VARCHAR(100) DEFAULT '/static/img/default.png',
    album_ID  BIGSERIAL   NOT NULL,
    index SMALLINT     NOT NULL,
    link     VARCHAR     NOT NULL,
    FOREIGN KEY (album_ID) REFERENCES albums (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE liked_artists
(
    user_ID   BIGSERIAL NOT NULL,
    artist_ID BIGSERIAL NOT NULL,
    FOREIGN KEY (user_ID) REFERENCES users (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (artist_ID) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (user_ID, artist_ID)
);

CREATE TABLE liked_albums
(
    user_ID  BIGSERIAL NOT NULL,
    album_ID BIGSERIAL NOT NULL,
    FOREIGN KEY (user_ID) REFERENCES users (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (album_ID) REFERENCES albums (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (user_ID, album_ID)

);

-- Полная информация, необходимая для отображения трека

CREATE VIEW track_with_artist AS -- TODO Изменить, тк трек не привязан в альбому
SELECT a.ID   as artist_ID,
       t.ID   as track_ID,
       t.name as track_name,
       a.name as artist_name,
       t.duration as duration,
       t.link as link,
       al.image as image,
       t.index as a_index
FROM artists a,
     albums al,
     tracks t
WHERE a.ID = al.artist_ID
  AND al.ID = t.album_ID;

-- Список треков артиста

CREATE VIEW artist_tracks AS
    SELECT a.ID as atrist_Id,
           t.ID as track_ID
FROM artists a,
     albums al,
     tracks t
WHERE a.ID = al.artist_ID AND al.ID = t.album_ID;


-- 	Id       string   `json:"id"`
-- 	Name     string `json:"name"`
-- 	Artist   string `json:"artist"`
-- 	Duration uint   `json:"duration"`
-- 	Image    string `json:"image"`
-- 	Link     string `json:"link"`
