# This Go script allows you to batch resize and compress videos using FFmpeg.

Prerequisites
FFmpeg: Ensure you have [FFmpeg](https://www.ffmpeg.org/) installed on your system.
Go: Ensure you have [Go](https://go.dev/) installed on your system.

**Usage**

```bash
go run main.go --input <<PATH>> --output <<PATH>>
```

--input: Path to the folder containing the videos you want to resize and compress.<br>
--output: Path to the folder where the resized and compressed videos will be saved.

**Example**

```bash
go run main.go --input /path/to/input/videos --output /path/to/output/videos
```

Output:

```bash
Finished in 20s seconds
File: example.mov
Input file size: 45.0 MiB
Output file size: 1.4 MiB
Saved: 43.7 MiB
```

This will resize and compress all the videos in the specified input folder and save the processed videos in the output folder.
