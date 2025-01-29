--Base Table
CREATE TABLE player (
    id INT auto_increment PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    cash INT UNSIGNED DEFAULT 2000,
    gold INT UNSIGNED DEFAULT 11,
    xp INT UNSIGNED DEFAULT 999999, -- 0,
    instant_cookings INT UNSIGNED DEFAULT 0,
    open_jobs TINYINT UNSIGNED DEFAULT 0,
    played_wheel BOOL DEFAULT FALSE,
    allow_friend_requests BOOL DEFAULT TRUE,
    friends TEXT DEFAULT "",
    friends_with_gifts TEXT DEFAULT "",
    allow_emails BOOL DEFAULT FALSE,
    email_verified BOOL DEFAULT TRUE,
    username VARCHAR(50) NOT NULL,
    avatar TEXT DEFAULT "1002$2#1022$2#1042$3#1052$4#1062$0#1082$0",
    is_banned BOOL DEFAULT FALSE,
    mastery TEXT DEFAULT '1201+1#1202+0#1203+2#1204+2#1205+0#1206+0#1207+0#1208+0#1209+0#1210+0#1211+0#1212+0#1213+0#1214+0#1215+0#1216+0#1217+0#1218+0#1219+0#1220+0#1221+0#1222+0#1223+0#1224+0#1225+0#1226+0#1227+0#1228+0#1229+0#1230+0#1231+0#1232+0#1233+0#1234+0#1235+0#1236+0#1237+0#1238+0#1239+0#1240+0#1241+0#1242+0#1243+0#1244+0#1245+0#1246+0#1247+0#1248+0#1249+0#1250+0#1251+0#1252+0#1253+0#1254+0#1255+0',
    achievement TEXT DEFAULT '2001+0#2002+0#2003+0#2004+0#2005+0#2006+0#2007+0#2008+0#2009+0#2010+0#2011+0#2012+0#2013+0#2014+0#2015+0#2016+0#2017+0#2018+0#2019+0#2020+0#2021+0#2022+0#2023+0#2024+0#2025+0#2026+0#2027+0#2028+0#2029+0#2030+0',
    last_login TEXT DEFAULT "",
    daily_login DATETIME DEFAULT CURRENT_TIMESTAMP,
    gifts TEXT DEFAULT "",
    sendable_gifts TEXT DEFAULT "",
    gift_refresh_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    access_level TINYINT DEFAULT 0
);

CREATE TABLE cafe (
    id INT auto_increment PRIMARY KEY,
    player_id INT NOT NULL,
    rating INT DEFAULT 50,
    luxury INT DEFAULT 0,
    size INT DEFAULT 8,
    tiles TEXT DEFAULT '7+101+101+101+101+101+101+101+101+4+4+4+4+4+4+4+101+4+4+4+4+4+4+4+101+4+4+4+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4',
    objects LONGTEXT collate utf8mb4_bin DEFAULT '[{"pos": [3, 0], "id": 901, "rotation": 0}, {"pos": [5, 0], "id": 901, "rotation": 0}, {"pos": [5, 1], "id": 601, "rotation": 3, "dish_id": -1, "dish_status": 0}, {"pos": [0, 2], "id": 201, "rotation": 0}, {"pos": [5, 2], "id": 401, "rotation": 0}, {"pos": [7, 2], "id": 601, "rotation": 3, "dish_id": -1, "dish_status": 0}, {"pos": [7, 3], "id": 401, "rotation": 0}, {"pos": [1, 4], "id": 351, "rotation": 0}, {"pos": [5, 4], "id": 401, "rotation": 0}, {"pos": [1, 5], "id": 252, "rotation": 0, "dish_id": -1}, {"pos": [3, 5], "id": 301, "rotation": 0, "dish_id": -1, "dish_amount": 0}, {"pos": [5, 5], "id": 601, "rotation": 1, "dish_id": -1, "dish_status": 0}, {"pos": [7, 5], "id": 401, "rotation": 0}, {"pos": [1, 6], "id": 252, "rotation": 0, "dish_id": -1}, {"pos": [3, 6], "id": 301, "rotation": 0, "dish_id": -1, "dish_amount": 0}, {"pos": [7, 6], "id": 601, "rotation": 1, "dish_id": -1, "dish_status": 0}, {"pos": [1, 7], "id": 252, "rotation": 0, "dish_id": -1}]' CHECK (json_valid (objects)),
    owner_name VARCHAR(50) NOT NULL,
    fridge_inv text DEFAULT '1314+3#1327+3',
    furniture_inv text DEFAULT '602+5#603+10',
    waiters LONGTEXT DEFAULT 'James+1002$0#1022$2#1042$6#1052$0#1062$0#1082$0+50'
    -- CONSTRAINT cafe_ibfk_1 FOREIGN KEY (player_id) REFERENCES player (id)
);

--- Data
--- Dummy1
INSERT INTO
    player (
        id,
        email,
        password,
        xp,
        cash,
        username,
        access_level
    )
VALUES
    (1, 'a@a.a', '1234', 99999999, 120999, 'bence', 3);

INSERT INTO
    cafe (id, player_id, owner_name)
VALUES
    (1, 1, 'bence');

--- Dummy2
INSERT INTO
    player (
        id,
        email,
        password,
        xp,
        cash,
        username,
        access_level
    )
VALUES
    (2, 'b@b.b', '1234', 99999999, 120999, 'bence2', 3);

INSERT INTO
    cafe (id, player_id, owner_name)
VALUES
    (2, 2, 'bence2');
