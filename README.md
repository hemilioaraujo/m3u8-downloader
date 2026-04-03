# M3U8 Downloader

A fast, concurrent Go application to download video streams from `.m3u8` playlists and save them directly as `.mp4` files. 

## Dependencies

To run this application, ensure you have the following installed on your system:

1. **Go (Golang)**: Version 1.25.5 or higher is required (as specified in `go.mod`).
2. **FFmpeg**: The core downloading and conversion is handled by FFmpeg. It must be installed and accessible via your system's `PATH`. 
   - *Ubuntu/Debian:* `sudo apt install ffmpeg`
   - *macOS:* `brew install ffmpeg`
   - *Windows:* Download from the [official site](https://ffmpeg.org/download.html) or install via `winget install ffmpeg`.

## How to Use

1. **Prepare your input list:** Create a JSON file (e.g., `videos.json`) containing an array of objects. Each object must have a `name` for the output file and the `url` to the `.m3u8` stream.

   Example `videos.json`:
   ```json
   [
     {
       "name": "my-cool-video",
       "url": "https://example.com/path/to/playlist.m3u8"
     },
     {
       "name": "second_video.mp4",
       "url": "https://example.com/path/to/another_playlist.m3u8"
     }
   ]
   ```

2. **Run the script:** Use the `go run` command and pass your JSON file as an argument.

   ```bash
   go run main.go videos.json
   ```

   *Alternatively, you can build a binary and run it:*
   ```bash
   go build -o m3u8-downloader main.go
   ./m3u8-downloader videos.json
   ```

3. **Check the output:** The application will automatically create a `downloads` directory (if it doesn't exist) in the same folder where the command was run, and save all the downloaded `.mp4` files there.

## Features

- **Concurrent Downloading:** Uses goroutines to download multiple `.m3u8` streams globally at the same time, speeding up large batch tasks.
- **Automatic `.mp4` Extention handling:** The tool automatically appends `.mp4` to file names if missing or replaces `.m3u8` extensions with `.mp4`.
- **Lossless Remuxing:** Instructs FFmpeg to directly copy streams (`-c copy`) into the MP4 container without re-encoding, preserving quality and keeping downloads lightweight.
