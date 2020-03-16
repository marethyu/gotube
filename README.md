# GoTube

![](https://img.shields.io/badge/version-v1.0-blue)

## Overview

This repository contains a single-file implementation of YouTube video downloader written in Go. It does not require any third-party packages, only built-in packages from the standard library. The code is compact and easily-readable.

## Installation

You can fetch one of the executable binaries (preferably latest one) from Releases. Unfortunately the binaries only support Windows x86-64 platform this time. However you can help me to make and upload more binaries for different platforms.

### From Source

If you don't have Go installed yet, download [here](https://golang.org/dl/). Clone this repository and go to the repository directory.
Then run the following command:
```markdown
go build gotube.go
```

## How It Works?

When GoTube receives a YouTube video URL as an input, it will extract a video id portion of the URL. Then it sends an HTTP GET request to a server through URL ```https://www.youtube.com/get_video_info``` with the video id as a query string. The GET request will return an encoded data containing YouTube video data which will be decoded by GoTube. It will then extract the video's title, file format and its download URL from the data. After getting all necessary data, GoTube makes a video file with specified title and file format in User's OS and sends another HTTP GET request to the download URL. The data received through the GET request will be copied to the created video file.

### Steps:

1. GoTube begins by extracting a video id from a YouTube URL (ex. if the URL is ```https://www.youtube.com/watch?v=9aUlHLaVmkc```, GoTube will extract ```9aUlHLaVmkc``` part).
2. GoTube sends an HTTP GET request to a server through ```https://www.youtube.com/get_video_info``` with the video id as a query string (ex. using id from the previous example, ```https://www.youtube.com/get_video_info?video_id=9aUlHLaVmkc```).
3. GoTube retrieves the video's title, file format and its download URL from an encoded data received through the GET request.
4. GoTube creates a video file in User's OS and make another GET request through the download URL.
5. GoTube copies contents of data received through the GET request to the created video file.

## Example Commands

```markdown
gotube https://www.youtube.com/watch?v=vUOHGL4Iv34 <download directory>
```

With verbose option:
```markdown
gotube https://www.youtube.com/watch?v=wJMkvlTAzHc <download directory> -v
```

```markdown
gotube https://www.youtube.com/watch?v=Hh_HyNfyOKs <download directory>
```

Replace \<download directory\> your download directory (ex. C:\users\marethyu\documents\poo)

## TODO
 - Option to download audio only
 - Download whole playlist
 - Progress bar?
