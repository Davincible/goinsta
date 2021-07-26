package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Davincible/goinsta/v2"
)

var (
	filepath = flag.String("filepath", "video.mp4", "Video file path")
)

func main() {
	flag.Parse()
	log.Println("filepath", *filepath)
	insta := goinsta.New(
		os.Getenv("INSTAGRAM_USERNAME"),
		os.Getenv("INSTAGRAM_PASSWORD"),
	)
	if err := insta.Login(); err != nil {
		log.Fatal(err)
	}

	defer insta.Logout()

	log.Println("Download random photo")
	var client http.Client
	request, err := http.NewRequest("GET", "https://picsum.photos/800/800", nil)
	if err != nil {
		log.Fatal(err)
	}
	thumbnail, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer thumbnail.Body.Close()

	log.Println("Read video file")
	file, err := os.Open(*filepath)
	if err != nil {
		log.Fatal(err)
	}

	postedVideo, err := insta.UploadVideo(bufio.NewReader(file), "awesomeVID", "awesome! :)", thumbnail.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Success upload video %s", postedVideo.ID)
}
