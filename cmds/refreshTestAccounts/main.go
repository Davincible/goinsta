//
// Use this program to parse instagram accounts from goinsta/tests/.env, and
//   create base64 encoded configs that can be used for testing.
//
// To add accounts, add to tests/.env: INSTAGRAM_ACT_<act-name>="<user>:<pass>"
//
// Also make sure to add a pixabay api key under PIXABAY_API_KEY="<key>" as
//   this is needed for some upload tests.
//
package main

import (
	"github.com/Davincible/goinsta/v3"
)

func main() {
	// Open File
	path := "../../tests/.env"
	err := goinsta.EnvProvision(path)
	if err != nil {
		panic(err)
	}
}
