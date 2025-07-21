--Base Table
CREATE TABLE player (
    id INT auto_increment PRIMARY KEY,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    cash BIGINT DEFAULT 2000,
    gold BIGINT DEFAULT 11,
    xp BIGINT DEFAULT 999999, -- 0,
    instant_cookings INT DEFAULT 0,
    open_jobs BIGINT UNSIGNED DEFAULT 0,
    coop_id INT DEFAULT -1,
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
    achievement LONG DEFAULT '2001+0#2002+0#2003+0#2004+0#2005+0#2006+0#2007+0#2008+0#2009+0#2010+0#2011+0#2012+0#2013+0#2014+0#2015+0#2016+0#2017+0#2018+0#2019+0#2020+0#2021+0#2022+0#2023+0#2024+0#2025+0#2026+0#2027+0#2028+0#2029+0#2030+0',
    last_login DATETIME DEFAULT CURRENT_TIMESTAMP,
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
    tiles LONGTEXT DEFAULT '7+101+101+101+101+101+101+101+101+4+4+4+4+4+4+4+101+4+4+4+4+4+4+4+101+4+4+4+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4+101+1+1+1+4+4+4+4',
    objects LONGTEXT DEFAULT '3+0+901+0#5+0+901+0#5+1+601+3+-1+0#0+2+201+0#5+2+401+0#7+2+601+3+-1+0#7+3+401+0#1+4+351+0#5+4+401+0#1+5+252+0+-1#3+5+301+0+-1+0#5+5+601+1+-1+0#7+5+401+0#1+6+252+0+-1#3+6+301+0+-1+0#7+6+601+1+-1+0#1+7+252+0+-1',
    owner_name VARCHAR(50) NOT NULL,
    fridge_inv TEXT DEFAULT '1314+3#1327+3',
    furniture_inv TEXT DEFAULT '602+5#603+10',
    waiters TEXT DEFAULT 'James+1002$0#1022$2#1042$6#1052$0#1062$0#1082$0+50'
);

CREATE TABLE coop (
    id INT auto_increment PRIMARY KEY,
    type INT NOT NULL, --
    dishes TEXT NOT NULL, -- dish_id+amount#
    host INT NOT NULL, --
    members TEXT NOT NULL,
    start DATETIME (3) DEFAULT CURRENT_TIMESTAMP,
    end DATETIME (3) NOT NULL
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
    (
        1,
        'a@a.a',
        '$2a$14$yCr6IBITIv.u9S7PcdLrHuyxS1lvPDWJ1OCF6SFBEZkJbBcRObfR.', -- pass: 1234
        99999999,
        120999,
        'bence',
        3
    );

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
    (
        2,
        'b@b.b',
        '$2a$14$yCr6IBITIv.u9S7PcdLrHuyxS1lvPDWJ1OCF6SFBEZkJbBcRObfR.', -- pass: 1234
        99999999,
        120999,
        'bence2',
        3
    );

INSERT INTO
    cafe (id, player_id, owner_name)
VALUES
    (2, 2, 'bence2');
