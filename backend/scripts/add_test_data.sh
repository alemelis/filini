#!/bin/bash

curl -X POST http://localhost:8080/upload \
     -F "id=1" \
     -F "title=Fantozzi" \
     -F "video_path=storage/Fantozzi-1.avi" \
     -F "srt_path=storage/Fantozzi-1.srt"
