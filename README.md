# GoTube v1

## Brief

A simple command line tool implemented with Go for downloading YouTube videos.

## Build Instructions

If you don't have Go installed yet, download [here](https://golang.org/dl/). Clone this repository and go to the repository directory.
Then run the following command:
```markdown
go build gotube.go
```

## Example Commands

```markdown
gotube https://www.youtube.com/watch?v=vUOHGL4Iv34 <download directory>
```

```markdown
gotube https://www.youtube.com/watch?v=wJMkvlTAzHc <download directory> -v
```

```markdown
gotube https://www.youtube.com/watch?v=Hh_HyNfyOKs <download directory>
```

Replace \<download directory\> your download directory (ex. C:\users\marethyu\documents\poo)

## TODO
 - Use Cobra
 - Download whole playlist
 - Progress bar?
