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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Davincible/goinsta"
)

func main() {
	// Open File
	path := "../../tests/.env"
	file, err := os.Open(path)
	checkErr(err)
	defer file.Close()

	// Read file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	checkErr(err)

	newBuf := new(bytes.Buffer)

	// Read Lines
	for buf.Len() != 0 {
		l, err := buf.ReadBytes(byte('\n'))
		checkErr(err)

		// Process Lines
		line := string(l)
		if strings.HasPrefix(line, "INSTAGRAM_ACT_") {
			// Extract Creds
			split := strings.Split(line, "=")
			name := split[0][14:]
			c := split[1][1:]
			creds := strings.Split(strings.Split(c, "\"")[0], ":")

			// Login
			fmt.Println("Processing", creds[0])
			insta := goinsta.New(creds[0], creds[1])
			err := insta.Login()
			checkErr(err)

			// Export Config
			enc, err := insta.ExportAsBase64String()
			checkErr(err)

			// Write Config
			_, err = newBuf.WriteString(line)
			checkErr(err)
			_, err = newBuf.WriteString(fmt.Sprintf("INSTAGRAM_BASE64_%s=\"%s\"\n\n", name, enc))
			checkErr(err)
		} else if !strings.HasPrefix(line, "INSTAGRAM") {
			_, err = newBuf.WriteString(line)
			checkErr(err)
		}
	}

	// Print Config to Stdout
	fmt.Printf("\n%s\n", string(newBuf.Bytes()))

	// Write File
	err = ioutil.WriteFile(path, newBuf.Bytes(), 0o644)
	checkErr(err)
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
