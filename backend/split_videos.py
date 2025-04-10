from pathlib import Path
import srt
import moviepy
import click
import psycopg2
from dotenv import load_dotenv
import os
import hashlib
import random
import time


def migrate_schema(conn):
    with conn.cursor() as cur:
        try:
            with open("db/schema.sql", "r") as schema_file:
                schema_sql = schema_file.read()

            cur.execute(schema_sql)
            conn.commit()
            print("Database schema migrated successfully")
        except psycopg2.errors.DuplicateTable:
            conn.rollback()
            print("Schema already exists, continuing...")
        except Exception as e:
            conn.rollback()
            print(f"Error during schema migration: {e}")
            raise


def split_video(conn, v: Path, s: Path, video_id: int, video_name: str, series: str, debug: bool = False):
    video = moviepy.VideoFileClip(v)

    for subtitle in srt.parse(s.read_text()):
        if subtitle.index in [10, "10"] and debug:
            break

        clip = video.subclipped(str(subtitle.start), str(subtitle.end))
        txt_clip = moviepy.TextClip(
            font="storage/LEMONMILK-Regular.otf",
            text=subtitle.content + "\n",
            font_size=18,
            color="yellow",
            method="caption",
            vertical_align="bottom",
            size=(544, 304),
            stroke_color="black",
            stroke_width=2,
        )

        subtitle_id = write_subtitle(conn, video_id, subtitle.content)
        if not subtitle_id:
            continue

        txt_clip = txt_clip.with_position("center").with_duration(str(subtitle.end - subtitle.start))
        overlay = moviepy.CompositeVideoClip([clip.without_audio(), txt_clip])

        random_hash = hashlib.md5(f"{series}{time.time()}{random.random()}{subtitle.content}".encode()).hexdigest()[:10]
        webm_path = Path("storage", "webm", series, f"{random_hash}.webm")
        overlay.write_videofile(webm_path)

        write_webm(conn, video_id, subtitle_id, webm_path)


def write_video(conn, video_name, series) -> int:
    cur = conn.cursor()
    try:
        cur.execute(
            "SELECT id FROM videos WHERE title = %s AND series = %s",
            (video_name, series),
        )
        result = cur.fetchone()

        if result:
            video_id = result[0]
        else:
            cur.execute(
                "INSERT INTO videos (title, series) VALUES (%s, %s) RETURNING id",
                (video_name, series),
            )
            video_id = cur.fetchone()[0]
            conn.commit()
    except Exception as e:
        print(f"Error writing video: {e}")
        conn.rollback()
        video_id = None
    finally:
        cur.close()
    return video_id


def write_subtitle(conn, video_id, subtitle_text) -> int:
    cur = conn.cursor()
    try:
        cur.execute(
            "SELECT id FROM subtitles WHERE video_id = %s AND text = %s",
            (video_id, subtitle_text),
        )
        result = cur.fetchone()

        if result:
            subtitle_id = None
        else:
            cur.execute(
                "INSERT INTO subtitles (video_id, text) VALUES (%s, %s) RETURNING id",
                (video_id, subtitle_text),
            )
            subtitle_id = cur.fetchone()[0]
            conn.commit()
    except Exception as e:
        print(f"Error writing subtitle: {e}")
        conn.rollback()
        subtitle_id = None
    finally:
        cur.close()
    return subtitle_id


def write_webm(conn, video_id, subtitle_id, webm_path):
    cur = conn.cursor()
    try:
        cur.execute(
            "SELECT id FROM webms WHERE video_id = %s AND subtitle_id = %s",
            (video_id, subtitle_id),
        )
        result = cur.fetchone()
        if not result:
            cur.execute(
                "INSERT INTO webms (video_id, subtitle_id, file_path) VALUES (%s, %s, %s) RETURNING id",
                (video_id, subtitle_id, str(webm_path)),
            )
            cur.fetchone()[0]
            conn.commit()
    except Exception as e:
        print(f"Error writing webm: {e}")
        conn.rollback()
    finally:
        cur.close()


@click.command()
@click.option(
    "--video_path",
    type=click.Path(exists=True, file_okay=True, dir_okay=False, path_type=Path),
    required=True,
    help="Path to the video file",
)
@click.option(
    "--subtitle_path",
    type=click.Path(exists=True, file_okay=True, dir_okay=False, path_type=Path),
    required=True,
    help="Path to the subtitle file",
)
@click.option(
    "--series",
    type=str,
    default="",
    help="Series name for the video",
)
def main(video_path, subtitle_path, series):
    video_name = video_path.stem.replace(" ", "_")

    load_dotenv()
    conn = psycopg2.connect(
        database=os.getenv("DB_NAME"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        host=os.getenv("DB_HOST", "localhost"),
        port=os.getenv("DB_PORT", "5432"),
    )

    migrate_schema(conn)

    video_id = write_video(conn, video_name, series)

    if not (webm_path := Path("storage", "webm", series)).exists():
        webm_path.mkdir(parents=True, exist_ok=True)

    split_video(conn, video_path, subtitle_path, video_id, video_name, series)


if __name__ == "__main__":
    main()
