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
    played_wheel TINYINT (1) DEFAULT 0,
    allow_friend_requests TINYINT (1) DEFAULT 1,
    friends TEXT DEFAULT "",
    allow_emails TINYINT (1) DEFAULT 0,
    email_verified TINYINT (1) DEFAULT 1,
    new_gifts TINYINT UNSIGNED DEFAULT 0,
    username VARCHAR(50) NOT NULL,
    avatar TEXT DEFAULT "",
    -- TODO: DELETE
    gender TINYINT (1) UNSIGNED DEFAULT 2,
    top_color TINYINT UNSIGNED DEFAULT 2,
    skin_color TINYINT UNSIGNED DEFAULT 2,
    hair_color TINYINT UNSIGNED DEFAULT 3,
    legs_color TINYINT UNSIGNED DEFAULT 4,
    -- TODO: delete
    is_banned TINYINT (1) DEFAULT 0 NOT NULL,
    mastery TEXT DEFAULT '1201+1#1202+0#1203+2#1204+2#1205+0#1206+0#1207+0#1208+0#1209+0#1210+0#1211+0#1212+0#1213+0#1214+0#1215+0#1216+0#1217+0#1218+0#1219+0#1220+0#1221+0#1222+0#1223+0#1224+0#1225+0#1226+0#1227+0#1228+0#1229+0#1230+0#1231+0#1232+0#1233+0#1234+0#1235+0#1236+0#1237+0#1238+0#1239+0#1240+0#1241+0#1242+0#1243+0#1244+0#1245+0#1246+0#1247+0#1248+0#1249+0#1250+0#1251+0#1252+0#1253+0#1254+0#1255+0',
    achievement TEXT DEFAULT '2001+0#2002+0#2003+0#2004+0#2005+0#2006+0#2007+0#2008+0#2009+0#2010+0#2011+0#2012+0#2013+0#2014+0#2015+0#2016+0#2017+0#2018+0#2019+0#2020+0#2021+0#2022+0#2023+0#2024+0#2025+0#2026+0#2027+0#2028+0#2029+0#2030+0',
    last_login TEXT DEFAULT "",
    daily_login DATETIME DEFAULT CURRENT_TIMESTAMP,
    gifts TEXT DEFAULT ""
);

CREATE TABLE cafe (
    id INT auto_increment PRIMARY KEY,
    player_id INT NOT NULL,
    rating INT DEFAULT 50,
    luxury INT DEFAULT 0,
    expansion_id INT DEFAULT 0,
    tiles TEXT DEFAULT '7+101+101+101+101+101+101+101+101+4+4+4+4+4+4+4+101+4+4+4+4+4+4+4+101+4+4+4+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4',
    objects LONGTEXT collate utf8mb4_bin DEFAULT '[{"pos": [3, 0], "id": 901, "rotation": 0}, {"pos": [5, 0], "id": 901, "rotation": 0}, {"pos": [5, 1], "id": 601, "rotation": 3, "dish_id": -1, "dish_status": 0}, {"pos": [0, 2], "id": 201, "rotation": 0}, {"pos": [5, 2], "id": 401, "rotation": 0}, {"pos": [7, 2], "id": 601, "rotation": 3, "dish_id": -1, "dish_status": 0}, {"pos": [7, 3], "id": 401, "rotation": 0}, {"pos": [1, 4], "id": 351, "rotation": 0}, {"pos": [5, 4], "id": 401, "rotation": 0}, {"pos": [1, 5], "id": 252, "rotation": 0, "dish_id": -1}, {"pos": [3, 5], "id": 301, "rotation": 0, "dish_id": -1, "dish_amount": 0}, {"pos": [5, 5], "id": 601, "rotation": 1, "dish_id": -1, "dish_status": 0}, {"pos": [7, 5], "id": 401, "rotation": 0}, {"pos": [1, 6], "id": 252, "rotation": 0, "dish_id": -1}, {"pos": [3, 6], "id": 301, "rotation": 0, "dish_id": -1, "dish_amount": 0}, {"pos": [7, 6], "id": 601, "rotation": 1, "dish_id": -1, "dish_status": 0}, {"pos": [1, 7], "id": 252, "rotation": 0, "dish_id": -1}]' CHECK (json_valid (objects)),
    owner_name VARCHAR(50) NOT NULL,
    fridge_inv text DEFAULT '1314+3#1327+3',
    furniture_inv text DEFAULT '602+5#603+10',
    waiters LONGTEXT DEFAULT '[{"waiter_id": 0, "name": "James", "gender": 2, "priority": 50, "avatar": "James+2+1002$0#1022$2#1042$6#1052$0#1062$0#1082$0"}]',
    CONSTRAINT cafe_ibfk_1 FOREIGN KEY (player_id) REFERENCES player (id)
);

--- Data
--- Dummy1
INSERT INTO
    player (id, email, password, xp, cash, username)
VALUES
    (1, 'a@a.a', '1234', 99999999, 120999, 'bence');

INSERT INTO
    cafe (id, player_id, owner_name)
VALUES
    (1, 1, 'bence');

--- Dummy2
INSERT INTO
    player (id, email, password, xp, cash, username)
VALUES
    (2, 'b@b.b', '1234', 99999999, 120999, 'bence2');

INSERT INTO
    cafe (id, player_id, owner_name)
VALUES
    (2, 2, 'bence2');
