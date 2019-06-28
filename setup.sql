CREATE TABLE IF NOT EXISTS games (
    category_id TEXT PRIMARY KEY NOT NULL,
    guild_id    TEXT NOT NULL,
    text_id     TEXT NOT NULL,
    voice_id    TEXT NOT NULL,
    dm_id       TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS maps (
    guild_id    TEXT NOT NULL,
    text_id     TEXT PRIMARY KEY NOT NULL,
    name        TEXT NOT NULL,
    message_id  TEXT NOT NULL
);

GRANT SELECT, UPDATE, INSERT, DELETE 
ON maps, games 
TO bots, postgres;