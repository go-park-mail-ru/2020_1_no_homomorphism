CREATE TABLE artists
(
    ID    BIGSERIAL PRIMARY KEY,
    name  VARCHAR(50) NOT NULL,
    image VARCHAR(100) DEFAULT '/static/img/default.png',
    genre VARCHAR(30)
);

CREATE OR REPLACE FUNCTION artists_trigger_func() RETURNS TRIGGER AS
$artists_trigger$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO artist_stat VALUES (NEW.ID, 0, 0, 0);
        RETURN NEW;
    END IF;
    IF (TG_OP = 'DELETE') THEN
        delete from artist_stat where artist_id = old.ID;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$artists_trigger$ LANGUAGE plpgsql;

CREATE TRIGGER artists_trigger
    AFTER INSERT or update or delete
    ON artists
    FOR EACH ROW
EXECUTE PROCEDURE artists_trigger_func();


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

CREATE OR REPLACE FUNCTION albums_trigger_func() RETURNS TRIGGER AS
$albums_trigger$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        update artist_stat set albums = albums + 1 where artist_id = new.artist_ID;
    END IF;
    IF (TG_OP = 'DELETE') THEN
        update artist_stat set albums = albums - 1 where artist_id = new.artist_ID;
    END IF;
    RETURN NULL;
END;
$albums_trigger$ LANGUAGE plpgsql;

CREATE TRIGGER albums_trigger
    AFTER INSERT or update or delete
    ON albums
    FOR EACH ROW
EXECUTE PROCEDURE albums_trigger_func();


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

CREATE OR REPLACE FUNCTION tracks_trigger_func() RETURNS TRIGGER AS
$tracks_trigger$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        update artist_stat set tracks = tracks + 1 where artist_id = new.artist_ID;
    END IF;
    IF (TG_OP = 'DELETE') THEN
        update artist_stat set tracks = tracks - 1 where artist_id = new.artist_ID;
    END IF;
    RETURN NULL;
END;
$tracks_trigger$ LANGUAGE plpgsql;

CREATE TRIGGER tracks_trigger
    AFTER INSERT or update or delete
    ON tracks
    FOR EACH ROW
EXECUTE PROCEDURE tracks_trigger_func();


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
    ID           BIGSERIAL PRIMARY KEY,
    login        VARCHAR(32)  NOT NULL UNIQUE,
    password     BYTEA        NOT NULL,
    name         VARCHAR(50)  NOT NULL,
    email        VARCHAR(320) NOT NULL UNIQUE,
    sex          VARCHAR(10)  NOT NULL,
    image        VARCHAR(100) DEFAULT '/static/img/avatar/default.png',
    liked_tracks integer[]    DEFAULT '{}'
);

CREATE OR REPLACE VIEW user_liked_tracks AS
with liked as (
    select unnest(liked_tracks) as liked_id
    from users
    where id = ?
)
select row_number() over () as index,
       t.ID                 as id,
       t.name               as name,
       a.name               as artist,
       t.duration           as duration,
       t.image              as image,
       artist_id            as artist_id,
       t.link               as link
from liked
         join tracks as t on liked.liked_id = t.ID
         join artists as a on t.artist_id = a.ID
order by index desc;

select *
from user_liked_tracks;

select unnest(liked_tracks) as liked_id
    from users
    where id = ?;

update users set liked_tracks = liked_tracks || '{545}' where id = 1002;
update users set liked_tracks = array_remove(liked_tracks, 45) where id = 1002;

CREATE OR REPLACE FUNCTION after_user_insert_func() RETURNS TRIGGER AS
$after_user_insert$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO user_stat VALUES (NEW.ID, 0, 0, 0, 0);
        RETURN NEW;
    END IF;
    IF (TG_OP = 'DELETE') THEN
        delete from user_stat where user_id = old.ID;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$after_user_insert$ LANGUAGE plpgsql;

CREATE TRIGGER after_user_insert
    AFTER INSERT or delete
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE after_user_insert_func();


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

CREATE OR REPLACE FUNCTION after_liked_artists_func() RETURNS TRIGGER AS
$after_liked_artists$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        update artist_stat set subscribers = subscribers + 1 where artist_id = new.artist_ID;
        update user_stat set artists = artists + 1 where user_ID = new.user_ID;
        RETURN NEW;
    END IF;
    IF (TG_OP = 'DELETE') THEN
        update artist_stat set subscribers = subscribers - 1 where artist_id = old.artist_ID;
        update user_stat set artists = artists - 1 where user_ID = old.user_ID;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$after_liked_artists$ LANGUAGE plpgsql;

CREATE TRIGGER after_liked_artists
    AFTER INSERT or DELETE
    ON liked_artists
    FOR EACH ROW
EXECUTE PROCEDURE after_liked_artists_func();


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
    private bool         default TRUE,
    FOREIGN KEY (user_ID) REFERENCES users (ID)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION before_playlist_insert_func() RETURNS TRIGGER AS
$before_playlist_insert$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        if NEW.image = '' then
            new.image = '/static/img/playlist/default.png';
        end if;
        update user_stat set playlists = playlists + 1 where user_ID = new.user_ID;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$before_playlist_insert$ LANGUAGE plpgsql;

CREATE TRIGGER before_playlist_insert
    BEFORE INSERT
    ON playlists
    FOR EACH ROW
EXECUTE PROCEDURE before_playlist_insert_func();

CREATE OR REPLACE FUNCTION after_playlist_delete_func() RETURNS TRIGGER AS
$after_playlist_insert$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        update user_stat set playlists = playlists - 1 where user_ID = old.user_ID;
    END IF;
    RETURN NULL;
END;
$after_playlist_insert$ LANGUAGE plpgsql;

CREATE TRIGGER after_playlist_delete
    after delete
    ON playlists
    FOR EACH ROW
EXECUTE PROCEDURE after_playlist_delete_func();


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

CREATE OR REPLACE FUNCTION before_playlist_track_insert_func() RETURNS TRIGGER AS
$before_playlist_track_insert$
DECLARE
    max_index smallint;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        max_index := (SELECT max(pl.index)
                      FROM playlist_tracks as pl
                      WHERE pl.playlist_ID = new.playlist_ID);
        IF max_index IS NULL then
            max_index = 0;
        end if;
        NEW.index := max_index + 1;

        if NEW.image = '' then
            new.image = '/static/img/track/default.png';
        end if;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$before_playlist_track_insert$ LANGUAGE plpgsql;

CREATE TRIGGER before_playlist_track_insert
    BEFORE INSERT
    ON playlist_tracks
    FOR EACH ROW
EXECUTE PROCEDURE before_playlist_track_insert_func();


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
       t.artist_id,
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


CREATE or replace VIEW user_artists AS
SELECT u.id    as user_id,
       a.ID    as artist_id,
       a.name  as name,
       a.image as image
FROM users u,
     artists a,
     liked_artists liked
WHERE u.id = liked.user_id
  AND a.id = liked.artist_id;

CREATE or replace VIEW sub_artists AS
SELECT liked.user_ID as user_id,
       a.ID          as artist_id,
       a.name        as name,
       a.image       as image
FROM artists a,
     liked_artists liked
WHERE a.id = liked.artist_id;

select *
from liked_artists
where user_id = 5
  and artist_id = 4;
insert into liked_artists (artist_id, user_id)
values (5, 5);

select *
from sub_artists
where user_id = 5;


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
       t.artist_id,
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
-- SELECT setval(pg_get_serial_sequence('artists', 'id'), coalesce(max(id) + 1, 1), false)
-- FROM artists;

CREATE VIEW artist_tracks AS
SELECT a.ID as atrist_Id,
       t.ID as track_ID
FROM artists a,
     tracks t
WHERE a.ID = t.artist_id;

update playlists
set private = not private
where id = 32;
select private
from playlists
where id = 32;
