CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    file_path TEXT NOT NULL
);

CREATE TABLE subtitles (
    id SERIAL PRIMARY KEY,
    video_id INT REFERENCES videos(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    start_time FLOAT NOT NULL,
    end_time FLOAT NOT NULL
);
