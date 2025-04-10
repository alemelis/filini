![img](https://c.tenor.com/ylMNO_634MQAAAAM/tenor.gif)

## backend

```
cd backend
docker compose up
```

```
cd backend
source .venv/bin/activate
python split_videos.py --help
Usage: split_videos.py [OPTIONS]

Options:
  --video_path FILE     Path to the video file  [required]
  --subtitle_path FILE  Path to the subtitle file  [required]
  --series TEXT         Series name for the video
  --help                Show this message and exit.
```

```
cd backend
go run cmd/main.go
```

## frontend

```
cd frontend
. ~/.nvm/nvm.sh
nvm use 18
npm run dev
```

- go to [http://localhost:5173/](http://localhost:5173/)
