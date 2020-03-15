DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS playlists CASCADE;
DROP TABLE IF EXISTS playlist_tracks CASCADE;
DROP VIEW IF EXISTS tracks_in_playlist CASCADE;
DROP VIEW IF EXISTS user_albums CASCADE;
DROP VIEW IF EXISTS user_artists CASCADE;

CREATE TABLE users
(
    ID       BIGSERIAL PRIMARY KEY,
    login    VARCHAR(32)  NOT NULL UNIQUE,
    password BYTEA        NOT NULL,
    name     VARCHAR(50)  NOT NULL,
    email    VARCHAR(320) NOT NULL UNIQUE,
    sex      VARCHAR(10)  NOT NULL,
    image    VARCHAR(100) DEFAULT '/static/img/avatar/default.png'
);

CREATE TABLE playlists
(
    ID      BIGSERIAL PRIMARY KEY,
    name    VARCHAR(50) NOT NULL,
    image   VARCHAR(100) DEFAULT '/static/img/avatar/default.png',
    user_ID BIGSERIAL   NOT NULL,
    FOREIGN KEY (user_ID) REFERENCES users (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE playlist_tracks
(
    playlist_ID BIGSERIAL NOT NULL,
    track_ID    BIGSERIAL NOT NULL,
    index       SMALLINT  NOT NULL,
    FOREIGN KEY (playlist_ID) REFERENCES playlists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_ID) REFERENCES tracks (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (playlist_ID, track_ID)
);

-- Список треков для плейлиста пользователя

CREATE VIEW tracks_in_playlist AS
SELECT u.ID     as user_ID,
       p.ID     as playlist_ID,
       t.ID     as track_ID,
       pt.index as index
FROM users u,
     playlists p,
     playlist_tracks pt,
     tracks t
WHERE u.ID = p.user_ID
  AND p.ID = pt.playlist_ID
  AND t.ID = pt.track_ID;

-- Список лайкнутых альбомов

CREATE VIEW user_albums AS
SELECT u.id     as user_ID,
       a.ID     as artist_ID,
       a.name   as artist_name, -- TODO Мб нужна аватарочка артиста??
       al.id    as album_id,
       al.name  as album_name,
       al.image as album_image
FROM users u,
     artists a,
     albums al,
     liked_albums liked
WHERE liked.user_id = u.id
  AND liked.album_id = al.id
  AND a.ID = al.artist_id;

-- Список пописок

CREATE VIEW user_artists AS
SELECT u.id    as user_id,
       a.ID    as artist_id,
       a.name  as artist_name,
       a.image as artist_image
FROM users u,
     artists a,
     liked_artists liked
WHERE u.id = liked.user_id
  AND a.id = liked.artist_id;



-- 	Id       uint   `json:"id"`
-- 	Password string `json:"password"`
-- 	Name     string `json:"name"`
-- 	Login    string `json:"login"`
-- 	Sex      string `json:"sex"`
-- 	Image    string `json:"image"`
-- 	Email    string `json:"email"`