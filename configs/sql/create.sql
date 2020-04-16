CREATE TABLE artists
(
    ID    BIGSERIAL PRIMARY KEY,
    name  VARCHAR(50) NOT NULL,
    image VARCHAR(100) DEFAULT '/static/img/default.png',
    genre VARCHAR(30)
);


CREATE TABLE albums
(
    ID          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    image       VARCHAR(100) DEFAULT '/static/img/default.png',
    release     DATE         NOT NULL,
    artist_name VARCHAR(50)  NOT NULL,
    artist_ID   BIGSERIAL    NOT NULL,
    FOREIGN KEY (artist_ID) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE tracks
(
    ID        BIGSERIAL PRIMARY KEY,
    name      VARCHAR(100) NOT NULL,
    duration  INTEGER      NOT NULL,
    image     VARCHAR DEFAULT '/static/img/track/default.png',
    link      VARCHAR      NOT NULL,
    artist_id BIGSERIAL    NOT NULL,
    FOREIGN KEY (artist_ID) REFERENCES artists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE album_tracks
(
    track_id BIGSERIAL   NOT NULL,
    album_id BIGSERIAL   NOT NULL,
    index    SMALLSERIAL NOT NULL,
    FOREIGN KEY (track_id) REFERENCES tracks (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (album_id) REFERENCES albums (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (track_id, album_id)
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

insert into users (login, password, name, email, sex)
VALUES ('tesgjfkldt', '$2a$04$FDsuiBXJdLkyENStXt7ituaN46L.SCxPbrV2ULNRK2Na1vL1dLeX2', 'TestName', 'test@mail.ru',
        'yes');

CREATE OR REPLACE FUNCTION after_user_insert_func() RETURNS TRIGGER AS
$after_user_insert$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO user_stat VALUES (NEW.ID, 0, 0, 0, 0);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$after_user_insert$ LANGUAGE plpgsql;

CREATE TRIGGER after_user_insert
    AFTER INSERT
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE after_user_insert_func();

-- todo добавить триггер на добавление трека в плейлист (index increase)

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


CREATE TABLE playlists
(
    ID      BIGSERIAL PRIMARY KEY,
    name    VARCHAR(50) NOT NULL,
    image   VARCHAR(100) DEFAULT '/static/img/default.png',
    user_ID BIGSERIAL   NOT NULL,
    FOREIGN KEY (user_ID) REFERENCES users (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE playlist_tracks
(
    playlist_ID BIGSERIAL   NOT NULL,
    track_ID    BIGSERIAL   NOT NULL,
    index       SMALLSERIAL NOT NULL,
    image       VARCHAR DEFAULT '/static/img/track/default.png',
    FOREIGN KEY (playlist_ID) REFERENCES playlists (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_ID) REFERENCES tracks (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    PRIMARY KEY (playlist_ID, track_ID)
);


CREATE TABLE user_stat
(
    user_id   BIGINT NOT NULL PRIMARY KEY,
    tracks    INT    NOT NULL,
    albums    INT    NOT NULL,
    playlists INT    NOT NULL,
    artists   INT    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);


CREATE TABLE artist_stat
(
    artist_id   BIGINT NOT NULL PRIMARY KEY,
    tracks      INT    NOT NULL,
    albums      INT    NOT NULL,
    subscribers INT    NOT NULL,
    FOREIGN KEY (artist_id) REFERENCES artists (id)
);


CREATE OR REPLACE VIEW full_track_info AS
SELECT t.ID    as track_id,
       a.ID    as artist_id,
       t.name  as track_name,
       a.name     artist_name,
       t.duration,
       t.link,
       t.image as track_image
FROM tracks t,
     artists a
WHERE a.ID = t.artist_id;


CREATE VIEW tracks_in_playlist AS
SELECT p.ID       as playlist_id,
       p.name     as playlist_name,
       p.image    as playlist_image,
       t.track_id as track_id,
       t.track_name,
       t.duration,
       t.artist_name,
       t.link,
       pt.index   as index,
       pt.image   as track_image
FROM playlists p,
     playlist_tracks pt,
     full_track_info t
WHERE p.ID = pt.playlist_ID
  AND t.track_id = pt.track_ID;


CREATE VIEW user_albums AS
SELECT u.id           as user_ID,
       al.id          as album_id,
       al.name        as album_name,
       al.image       as album_image,
       al.artist_name as artist_name,
       al.artist_ID   as artist_id
FROM users u,
     albums al,
     liked_albums liked
WHERE liked.user_id = u.id
  AND liked.album_id = al.id;


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


CREATE VIEW full_album_info AS
SELECT al.id    as album_id,
       al.name  as album_name,
       al.image as album_image,
       ar.ID    as artist_id,
       ar.name  as artist_name,
       ar.genre as artist_genre,
       ar.image as artist_image
FROM albums as al
         JOIN artists as ar ON al.artist_ID = ar.ID;


CREATE VIEW tracks_in_album AS
SELECT a.ID    as album_id,
       t.track_id,
       t.artist_name,
       t.track_name,
       t.duration,
       t.link,
       at.index,
       a.image as track_image
FROM album_tracks as at,
     albums as a,
     full_track_info as t
WHERE at.track_id = t.track_id
  AND at.album_id = a.ID;

-- Если при вставки пишет, что id повторяется, значит траблы с последовательностью, ее надо обновить:
SELECT setval(pg_get_serial_sequence('artists', 'id'), coalesce(max(id) + 1, 1), false)
FROM artists;


CREATE VIEW artist_tracks AS
SELECT a.ID as atrist_Id,
       t.ID as track_ID
FROM artists a,
     tracks t
WHERE a.ID = t.artist_id;

INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Cвитер', 244, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Cвитер.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Ах мамочка', 227, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Ах мамочка.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Бабы-стервы', 183, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Бабы-стервы.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Белые розы', 252, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Белые розы.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Водка без пива (zaycev.net)', 150, '/static/img/track/verka.jpg',
        '/static/audio/Верка Сердючка - Водка без пива.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Горiлка', 172, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Горiлка.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Даже если вам немного за 30', 204, '/static/img/track/verka.jpg',
        '/static/audio/Верка Сердючка - Даже если вам немного за 30.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('День рождения', 168, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - День рождения.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Жениха Хотела', 203, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Жениха Хотела.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Он бы подошел, я бы отвернулась ', 223, '/static/img/track/verka.jpg',
        '/static/audio/Верка Сердючка - Он бы подошел, я бы отвернулась .mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('С Днём Рождения', 171, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - С Днём Рождения.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Свадебная', 181, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Свадебная.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Сигареточка', 175, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Сигареточка.mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Эх, свадьба, свадьба.. (zaycev.net)', 177, '/static/img/track/verka.jpg',
        '/static/audio/Верка Сердючка - Эх, свадьба, свадьба...mp3', 4);
INSERT INTO tracks (name, duration, image, link, artist_id)
VALUES ('Я Не Поняла', 227, '/static/img/track/verka.jpg', '/static/audio/Верка Сердючка - Я Не Поняла.mp3', 4);

