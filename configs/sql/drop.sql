DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS playlists CASCADE;
DROP TABLE IF EXISTS playlist_tracks CASCADE;
DROP VIEW IF EXISTS tracks_in_playlist CASCADE;
DROP VIEW IF EXISTS user_albums CASCADE;
DROP VIEW IF EXISTS user_artists CASCADE;





-- 	Id       uint   `json:"id"`
-- 	Password string `json:"password"`
-- 	Name     string `json:"name"`
-- 	Login    string `json:"login"`
-- 	Sex      string `json:"sex"`
-- 	Image    string `json:"image"`
-- 	Email    string `json:"email"`