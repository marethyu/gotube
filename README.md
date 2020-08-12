# GoTube

[![Build Status](https://api.travis-ci.com/marethyu/gotube.svg?branch=master)](https://travis-ci.com/marethyu/gotube)
![](https://img.shields.io/badge/version-v1.2-blue)

## Overview

This repository contains a single-file implementation of YouTube video downloader written in Go. It does not require any third-party packages, only built-in packages from the standard library. The code is compact and easily-readable.

## Installation

You can fetch one of the executable binaries (preferably latest one) from Releases. The binaries only support Windows (32 and 64 bit) and Linux (32 and 64 bit) this time.

You can also install GoTube directly from command line. Simply just type a command below in your terminal:
```markdown
go get github.com/marethyu/gotube
```

## Building From Source

Run the following commands:
```markdown
git clone https://github.com/marethyu/gotube.git
cd gotube
go get golang.org/x/sync/errgroup
go build gotube.go
```

## How It Works?

When GoTube receives a YouTube video URL as an input, it will extract a video id portion of the URL. Then it sends a HTTP GET request to a server through URL ```https://www.youtube.com/get_video_info``` with the video id as a query string. The GET request will return an encoded data containing YouTube video data which will be decoded by GoTube. It will then extract the video's title, file format and its download URL from the data. After getting all necessary data, GoTube makes a video file with specified title and file format in User's OS and sends another HTTP GET request to the download URL. The data received through the GET request will be copied to the created video file.

### Steps:

1. GoTube begins by extracting a video id from a YouTube URL (ex. if the URL is ```https://www.youtube.com/watch?v=9aUlHLaVmkc```, GoTube will extract ```9aUlHLaVmkc``` part).
2. GoTube sends an HTTP GET request to a server through ```https://www.youtube.com/get_video_info``` with the video id as a query string (ex. using id from the previous example, ```https://www.youtube.com/get_video_info?video_id=9aUlHLaVmkc```).
3. GoTube retrieves the video's title, file format and its download URL from an encoded data received through the GET request.
4. GoTube creates a video file in User's OS and make another GET request through the download URL.
5. GoTube copies contents of data received through the GET request to the created video file.

## Example Usage

Basic (install to the current directory):
```markdown
gotube https://www.youtube.com/watch?v=2fsMUqYix0c
```

Download multiple videos:
```markdown
gotube https://www.youtube.com/watch?v=2fsMUqYix0c https://www.youtube.com/watch?v=wJMkvlTAzHc
```

Install to the specified download directory (\<download directory\>):
```markdown
gotube -outdir=<download directory> https://www.youtube.com/watch?v=MXStYQSLd_M
```

With verbose option:
```markdown
gotube -v https://www.youtube.com/watch?v=wJMkvlTAzHc
```

Option to download audio (requires ffmpeg, get one [here](https://github.com/adaptlearning/adapt_authoring/wiki/Installing-FFmpeg) if you don't have one installed):
```markdown
gotube -a https://www.youtube.com/watch?v=Hh_HyNfyOKs
```

## TODO
 - Download whole playlist

## Contributing

I appreciate any kind of contributions. It can include code changes, new features and recommendations.

## License

It's licensed under [BSD-3-Clause](LICENSE).
