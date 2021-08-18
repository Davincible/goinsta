package main

import (
	"flag"
	"log"

	"github.com/Davincible/goinsta"
)

var session = flag.String("session", "/tmp/session", "instagram session")

func main() {
	flag.Parse()
	insta, err := goinsta.Import(*session)
	if err != nil {
		log.Fatal(err)
	}
	followerUser := make([]goinsta.User, 0)
	followingUser := make([]goinsta.User, 0)
	followers := insta.Account.Followers()
	following := insta.Account.Following()

	for followers.Next() {
		followerUser = append(followerUser, followers.Users...)
	}

	for following.Next() {
		followingUser = append(followingUser, following.Users...)
	}

	log.Printf("Followers: %d Following: %d", len(followerUser), len(followingUser))
}
