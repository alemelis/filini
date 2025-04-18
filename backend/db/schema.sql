CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    series TEXT NOT NULL
);

CREATE TABLE subtitles (
    id SERIAL PRIMARY KEY,
    video_id INT REFERENCES videos (id) ON DELETE CASCADE,
    text TEXT NOT NULL
);

CREATE TABLE webms (
    id SERIAL PRIMARY KEY,
    video_id INT REFERENCES videos (id) ON DELETE CASCADE,
    subtitle_id INT REFERENCES subtitles (id) ON DELETE CASCADE,
    file_path TEXT NOT NULL
);
