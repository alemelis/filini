import psycopg2
import os
from math import floor, ceil

# Connect to PostgreSQL and fetch timestamps
conn = psycopg2.connect(
    dbname="filini",
    user="postgres",
    password="supersecret",
    host="localhost",
    port="5432",
)

cur = conn.cursor()
cur.execute("SELECT start_time, end_time FROM subtitles ORDER BY start_time")
timestamps = cur.fetchall()
cur.close()
conn.close()

# Write timestamps to a file
times = []
for start, end in timestamps:
    times.append(floor(start))
    times.append(floor(end))
    print(times[-2:])

# Split the video using ffmpeg
os.system(
    f"ffmpeg -i /Users/ale/repos/filini/backend/storage/Fantozzi-1.avi -f segment -segment_times {','.join(map(str, times))} -reset_timestamps 1 -write_empty_segments 1 -c:v libx264 -map 0 /Users/ale/repos/filini/backend/storage/segments/output_%03d.mp4"
)

# # Convert segments to GIFs
# for i in range(len(timestamps)):
#     os.system(
#         f"ffmpeg -i /Users/ale/repos/filini/backend/storage/segments/output_{i:03d}.mp4 -vf 'fps=10,scale=320:-1:flags=lanczos' /Users/ale/repos/filini/backend/storage/gifs/output_{i:03d}.gif"
#     )
