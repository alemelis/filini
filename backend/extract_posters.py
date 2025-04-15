from pathlib import Path

import moviepy


def extract_poster(webm_path):
    clip = moviepy.VideoFileClip(webm_path)

    # Save the first frame (at t=0) as an image
    poster_path = str(webm_path).replace(".webm", ".webp").replace("/webm/", "/poster/")
    clip.save_frame(poster_path, t=0)

    # Close the clip to release resources
    clip.close()


poster_path = Path("storage", "poster", "fantozzi")
poster_path.mkdir(parents=True, exist_ok=True)

for webm in Path("storage", "webm", "fantozzi").glob("*.webm"):
    try:
        extract_poster(webm)
    except Exception as e:
        print(f"Error processing {webm}: {e}")
