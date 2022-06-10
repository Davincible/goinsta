### Fork

This repository has been forked from [ahmdrz/goinsta](https://github.com/ahmdrz/goinsta). 
As the maintainer of this repositry has been absend the last few months, and 
the code in the repository was based on a 2 year old instagram app version, 
since which a lot has changed, I have taken the courtesy to build upon his 
great framework and update the code to be compatible with apk v195.0.0.31.123 
(July 6, 2021). After migrating the endpoints and adding new ones, there are 
are few breaking changes. You can check the full walkthrough documentation in
the [wiki](https://github.com/Davincible/goinsta/wiki/1.-Getting-Started), 
and looking at the code to further understand how it works is encouraged.

#### Golang + Instagram Private API
![goinsta logo](https://raw.githubusercontent.com/Davincible/goinsta/v1/resources/goinsta-image.png)

> Unofficial Instagram API for Golang

[![Build Status](https://travis-ci.org/Davincible/goinsta.svg?branch=master)](https://travis-ci.org/Davincible/goinsta) [![GoDoc](https://godoc.org/github.com/Davincible/goinsta?status.svg)](https://godoc.org/github.com/Davincible/goinsta) [![Go Report Card](https://goreportcard.com/badge/github.com/Davincible/goinsta)](https://goreportcard.com/report/github.com/Davincible/goinsta) [![Gitter chat](https://badges.gitter.im/goinsta/community.png)](https://gitter.im/goinsta/community)

### Features

* **HTTP2 by default. Goinsta uses HTTP2 client enhancing performance.**
* **Object independency. Can handle multiple instagram accounts.**
* **Like Instagram mobile application**. Goinsta is very similar to Instagram official application.
* **Simple**. Goinsta is made by lazy programmers!
* **Backup methods**. You can use Export`and Import`functions.
* **Security**. Your password is only required to login. After login your password is deleted.
* ~~**No External Dependencies**. GoInsta will not use any Go packages outside of the standard library.~~ goinsta now uses [chromedp](https://github.com/chromedp/chromedp) as headless browser driver to solve challanges and checkpoints.

### Package installation 

`go get -u github.com/Davincible/goinsta/v3@latest`

### Example

```go
package main

import (
	"fmt"

	"github.com/Davincible/goinsta/v3"
)

func main() {  
  insta := goinsta.New("USERNAME", "PASSWORD")
  
  // Only call Login the first time you login. Next time import your config
  err := insta.Login()
  if err != nil {
          panic(err)
  }

  // Export your configuration
  // after exporting you can use Import function instead of New function.
  // insta, err := goinsta.Import("~/.goinsta")
  // it's useful when you want use goinsta repeatedly.
  // Export is deffered because every run insta should be exported at the end of the run
  //   as the header cookies change constantly.
  defer insta.Export("~/.goinsta")

  ...
}
```

For the full documentation, check the [wiki](https://github.com/Davincible/goinsta/v3/wiki/1.-Getting-Started), or run `go doc -all`.

### Legal

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by Instagram or any of its affiliates or subsidiaries. This is an independent and unofficial API. Use at your own risk.

### Versioning

Goinsta used gopkg.in as versioning control. Stable new API is the version v3.0. You can get it using:

```bash
$ go get -u -v github.com/Davincible/goinsta/v3
```

Or 

If you have `GO111MODULE=on`

```
$ go get -u github.com/Davincible/goinsta/v3
```

