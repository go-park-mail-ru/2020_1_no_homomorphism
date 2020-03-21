DROP TABLE IF EXISTS tracks CASCADE;
DROP TABLE IF EXISTS albums CASCADE;
DROP TABLE IF EXISTS artists CASCADE;
DROP TABLE IF EXISTS liked_artists CASCADE;
DROP TABLE IF EXISTS liked_albums CASCADE;
DROP VIEW IF EXISTS track_with_artist CASCADE;
DROP TABLE IF EXISTS artist_tracks CASCADE;
DROP TABLE IF EXISTS album_tracks CASCADE;
DROP VIEW IF EXISTS full_track_info CASCADE;

CREATE TABLE artists -- TODO Подумать над полями
(
    ID    BIGSERIAL PRIMARY KEY,
    name  VARCHAR(50) NOT NULL,
    image VARCHAR(100) DEFAULT '/static/img/default.png',
    genre VARCHAR(30)
);

CREATE TABLE albums
(
    ID        BIGSERIAL PRIMARY KEY,
    name      VARCHAR(50) NOT NULL,
    image     VARCHAR(100) DEFAULT '/static/img/default.png',
    artist_ID BIGSERIAL   NOT NULL,
    FOREIGN KEY (artist_ID) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

SELECT *
FROM albums a
WHERE a.artist_id = 1;

CREATE TABLE tracks
(
    ID        BIGSERIAL PRIMARY KEY,
    name      VARCHAR(50) NOT NULL,
    duration  INTEGER     NOT NULL,
    --image    VARCHAR(100) DEFAULT '/static/img/default.png',
    link      VARCHAR     NOT NULL,
    artist_id BIGSERIAL   NOT NULL,
    FOREIGN KEY (artist_ID) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE album_tracks
(
    track_id BIGSERIAL NOT NULL,
    album_id BIGSERIAL NOT NULL,
    FOREIGN KEY (track_id) REFERENCES tracks (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (album_id) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (track_id, album_id)
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
SELECT p.ID       as playlist_id,
       p.name     as album_name,
       p.image    as album_image,
       t.track_id as track_id,
       t.track_name,
       t.duration,
       t.artist_name,
       t.link,
       pt.index   as index
FROM playlists p,
     playlist_tracks pt,
     full_track_info t
WHERE p.ID = pt.playlist_ID
  AND t.track_id = pt.track_ID;

SELECT track_id, track_name, artist_name, duration, link FROM tracks_in_playlist WHERE playlist_id = 4;

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

-- Полная информация, необходимая для отображения трека

-- CREATE VIEW track_with_artist AS -- TODO Изменить, тк трек не привязан в альбому
-- SELECT a.ID       as artist_ID,
--        t.ID       as track_ID,
--        t.name     as track_name,
--        a.name     as artist_name,
--        t.duration as duration,
--        t.link     as link,
--        al.image   as image,
--        t.index    as a_index
-- FROM artists a,
--      albums al,
--      tracks t
-- WHERE a.ID = al.artist_ID
--   AND al.ID = t.album_ID;

-- Список треков артиста

CREATE VIEW artist_tracks AS
SELECT a.ID as atrist_Id,
       t.ID as track_ID
FROM artists a,
     tracks t
WHERE a.ID = t.artist_id;

INSERT INTO artists (name, genre)
VALUES ('Vasya MC', 'hip-hop');

INSERT INTO tracks (name, duration, link, artist_id)
VALUES ('Rap God', 134, 'http://hip.hop', 1);

INSERT INTO artists (name, genre)
VALUES ('Metalll Band', 'metall');

INSERT INTO tracks (name, duration, link, artist_id)
VALUES ('Death 2', 215, 'http://death.com/666.mp3', 2);

CREATE VIEW full_track_info AS
SELECT t.ID as track_id, a.ID as artist_id, t.name as track_name, a.name artist_name, t.duration, t.link
FROM tracks t,
     artists a
WHERE a.ID = t.artist_id;

SELECT *
FROM full_track_info
WHERE track_id = 1;


