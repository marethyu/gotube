/**
 * Copyright (c) 2020, Jimmy Yang <codingexpert123@gmail.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are
 * permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this list of
 * conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice, this list of
 * conditions and the following disclaimer in the documentation and/or other materials provided
 * with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors may be used
 * to endorse or promote products derived from this software without specific prior written
 * permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS
 * OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
 * AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER
 * OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
 * ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
 * OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 ******************************************************************************************
 * GoTube: Simple YouTube Video Downloader (v1.1)
 *
 * REQUIRES: ffmpeg
 * WARNING: The download process might be very slow and will destroy your computer if it happens. (LOL)
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var percent int

type WriteCounter struct {
	BytesDownloaded int64
	TotalBytes      int64
}

func DisplayStatus() {
	for percent < 100 {
		fmt.Printf("\rGoTube: Download progress: %%%d complete", percent)
	}

	fmt.Println("\rGoTube: Download progress: %100 complete")
}

func (pWc *WriteCounter) Write(b []byte) (n int, err error) {
	n = len(b)
	pWc.BytesDownloaded += int64(n)
	percent = int(math.Round(float64(pWc.BytesDownloaded) * 100.0 / float64(pWc.TotalBytes)))
	return
}

// Check if the given file/directory exists
func Exists(path string) (bool, error, os.FileInfo) {
	fi, err := os.Stat(path)

	if err == nil {
		return true, nil, fi
	}
	if os.IsNotExist(err) {
		return false, nil, fi
	}

	return true, err, fi
}

func GetVideoID(videoURL string) (string, error) {
	u, err := url.Parse(videoURL)
	return u.Query()["v"][0], err
}

// Go's version of PHP's parse_str
// Shamelessly stolen from https://github.com/syyongx/php2go/blob/master/php.go
func ParseStr(encodedString string, result map[string]interface{}) error {
	// build nested map.
	var build func(map[string]interface{}, []string, interface{}) error

	build = func(result map[string]interface{}, keys []string, value interface{}) error {
		length := len(keys)
		// trim ',"
		key := strings.Trim(keys[0], "'\"")
		if length == 1 {
			result[key] = value
			return nil
		}

		// The end is slice. like f[], f[a][]
		if keys[1] == "" && length == 2 {
			// todo nested slice
			if key == "" {
				return nil
			}
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{value}
				return nil
			}
			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
			}
			result[key] = append(children, value)
			return nil
		}

		// The end is slice + map. like f[][a]
		if keys[1] == "" && length > 2 && keys[2] != "" {
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{}
				val = result[key]
			}
			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
			}
			if l := len(children); l > 0 {
				if child, ok := children[l-1].(map[string]interface{}); ok {
					if _, ok := child[keys[2]]; !ok {
						_ = build(child, keys[2:], value)
						return nil
					}
				}
			}
			child := map[string]interface{}{}
			_ = build(child, keys[2:], value)
			result[key] = append(children, child)

			return nil
		}

		// map. like f[a], f[a][b]
		val, ok := result[key]
		if !ok {
			result[key] = map[string]interface{}{}
			val = result[key]
		}
		children, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected type 'map[string]interface{}' for key '%s', but got '%T'", key, val)
		}

		return build(children, keys[1:], value)
	}

	// split encodedString.
	parts := strings.Split(encodedString, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return err
		}
		for key[0] == ' ' {
			key = key[1:]
		}
		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return err
		}

		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}
			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}
		keys[0] = first

		// build nested map
		if err := build(result, keys, value); err != nil {
			return err
		}
	}

	return nil
}

func DownloadYTVideo(videoURL, outputDirectory string, verbose, audio bool) {
	isMatch, _ := regexp.MatchString("https://www\\.youtube\\.com/watch\\?v=[\\w-]+", videoURL) // TODO need better regex pattern

	if !isMatch {
		log.Fatal("GoTube: Invalid YouTube URL!")
	}

	doesExist, _, fi := Exists(outputDirectory)

	if !doesExist {
		log.Fatal("GoTube: The output directory doesn't exist!")
	}

	if !fi.Mode().IsDir() {
		log.Fatal("GoTube: The directory is a file!")
	}

	id, _ := GetVideoID(videoURL)
	u := "https://www.youtube.com/get_video_info?video_id=" + id

	if verbose {
		fmt.Printf("GoTube: Making a HTTP GET request thru %s...\n", u)
	}

	resp, _ := http.Get(u)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("GoTube: Bad status: %s (%s)", resp.Status, http.StatusText(resp.StatusCode)))
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)

	data := make(map[string]interface{})
	ParseStr(string(byteArray[:]), data)

	// We only need to retrieve video title, format and download url nothing else

	var videoData map[string]interface{}
	json.Unmarshal([]byte(data["player_response"].(string)), &videoData)

	videoDetails := videoData["videoDetails"].(map[string]interface{})
	streamingData := videoData["streamingData"].(map[string]interface{})
	formats := streamingData["formats"].([]interface{})

	// Let's try the first format...
	moreData := formats[0].(map[string]interface{})
	moreData["mime"] = moreData["mimeType"]
	s := moreData["mime"].(string)

	title := strings.Replace(strings.ToLower(videoDetails["title"].(string)), " ", "_", -1)
	format := s[strings.Index(s, "/")+1 : strings.Index(s, ";")]
	downloadURL := moreData["url"].(string)

	// Remove characters like ':' and '?' in the video title
	re := regexp.MustCompile("[^A-Za-z0-9.\\_\\-]")
	fileName := re.ReplaceAllString(title+"."+format, "")
	path := filepath.Join(outputDirectory, fileName)

	if verbose {
		fmt.Printf("GoTube: Creating a file %s...\n", path)
	}

	output, _ := os.Create(path)
	defer output.Close()

	client := &http.Client{}

	// Determine the video size in bytes
	resp, err := client.Head(downloadURL)

	if err != nil {
		panic(err)
	}

	videoSize, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	request, _ := http.NewRequest("GET", downloadURL, nil)
	request.Header.Set("Cache-Control", "public")
	request.Header.Set("Content-Description", "File Transfer")
	request.Header.Set("Content-Disposition", "attachment; filename="+fileName)
	request.Header.Set("Content-Type", "application/zip")
	request.Header.Set("Content-Transfer-Encoding", "binary")

	if verbose {
		fmt.Printf("GoTube: Making another HTTP GET Request thru %s...\n", downloadURL)
	}

	resp, _ = client.Do(request)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("GoTube: Bad status: %s (%s)", resp.Status, http.StatusText(resp.StatusCode)))
	}

	var body io.Reader

	if !verbose {
		body = resp.Body
	} else {
		go DisplayStatus()
		body = io.TeeReader(resp.Body, &WriteCounter{0, videoSize}) // Pipe stream
	}

	_, err = io.Copy(output, body)

	if err != nil {
		log.Fatal("GoTube: Unable to download the video! :(")
	} else if verbose {
		fmt.Println("GoTube: The video downloaded successfully! :))")
	}

	if audio {
		audioFile := filepath.Join(outputDirectory, strings.TrimRight(fileName, filepath.Ext(fileName))+".mp3")

		if verbose {
			fmt.Printf("GoTube: Creating a file %s...\n", audioFile)
		}

		ffmpeg, err := exec.LookPath("ffmpeg")

		if err != nil {
			log.Fatal("ffmpeg not found")
		}

		cmd := exec.Command(ffmpeg, "-i", path, "-vn", "-ar", "44100", "-ac", "1", "-b:a", "32k", "-f", "mp3", audioFile)

		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		err = cmd.Run()

		if err != nil {
			log.Fatal(err)
		} else if verbose {
			fmt.Println("GoTube: The video audio extracted successfully! :))")
		}
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: gotube -vidurl=<YT_VID_URL> [-outdir=<OUT_DIRECTORY>] [-v] [-a]\n")
	}

	var videoURL string
	var outputDirectory string
	var verbose bool
	var audio bool

	flag.StringVar(&videoURL, "vidurl", "", "URL of a YouTube video")
	flag.StringVar(&outputDirectory, "outdir", ".", "Directory where you want the video to be downloaded")
	flag.BoolVar(&verbose, "v", false, "If true, GoTube will display detailed download process")
	flag.BoolVar(&audio, "a", false, "If true, GoTube will download video's audio as well")

	flag.Parse()

	if videoURL == "" || outputDirectory == "" {
		flag.Usage()
		os.Exit(1)
	}

	DownloadYTVideo(videoURL, outputDirectory, verbose, audio)
}
