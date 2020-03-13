DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users
(
    ID       BIGSERIAL PRIMARY KEY,
    login    VARCHAR(32)  NOT NULL UNIQUE,
    password BYTEA        NOT NULL,
    name     VARCHAR(50)  NOT NULL,
    email    VARCHAR(320) NOT NULL UNIQUE,
    sex      VARCHAR(10)  NOT NULL,
    image    VARCHAR(100) DEFAULT '/static/default.png'
);


-- защитить бд фаерволом или паролем
-- prepared statement
-- max open connect
-- easy json, рефлексия

-- 	Id       uint   `json:"id"`
-- 	Password string `json:"password"`
-- 	Name     string `json:"name"`
-- 	Login    string `json:"login"`
-- 	Sex      string `json:"sex"`
-- 	Image    string `json:"image"`
-- 	Email    string `json:"email"`