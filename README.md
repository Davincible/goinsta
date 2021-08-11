### Fork

This repository has been forked from [ahmdrz/goinsta](https://github.com/ahmdrz/goinsta). As the maintainer of this repositry has been absend the last few months, and the code in the repository was based on a 2 year old instagram app version, since which a lot has changed, I have taken the courtesy to bulid upon his great framework and update the code to be compatible with apk v195.0.0.31.123 (July 6, 2021). Currently, most endpoints have been moved over to the new version, however I have not yet writen tests for all. If you run into problems please open an issue. Documentation is also still on the TODO, as there are a few breaking changes, but nothing a glance at the code can't solve.

#### Golang + Instagram Private API
<p align="center"><img width=100% src="https://raw.githubusercontent.com/Davincible/goinsta/v1/resources/goinsta-image.png"></p>

> Unofficial Instagram API for Golang

[![Build Status](https://travis-ci.org/Davincible/goinsta.svg?branch=master)](https://travis-ci.org/Davincible/goinsta) [![GoDoc](https://godoc.org/github.com/Davincible/goinsta?status.svg)](https://godoc.org/github.com/Davincible/goinsta) [![Go Report Card](https://goreportcard.com/badge/github.com/Davincible/goinsta)](https://goreportcard.com/report/github.com/Davincible/goinsta) [![Gitter chat](https://badges.gitter.im/goinsta/community.png)](https://gitter.im/goinsta/community)

### Features

* **HTTP2 by default. Goinsta uses HTTP2 client enhancing performance.**
* **Object independency. Can handle multiple instagram accounts.**
* **Like Instagram mobile application**. Goinsta is very similar to Instagram official application.
* **Simple**. Goinsta is made by lazy programmers!
* **Backup methods**. You can use `Export` and `Import` functions.
* **Security**. Your password is only required to login. After login your password is deleted.
* **No External Dependencies**. GoInsta will not use any Go packages outside of the standard library.

### Package installation 

`go get -u -v github.com/Davincible/goinsta@latest`

### Example

```go
package main

import (
	"fmt"

	"github.com/Davincible/goinsta"
)

func main() {  
  insta := goinsta.New("USERNAME", "PASSWORD")
  
  err := insta.Login()
  if err != nil {
          panic(err)
  }

  // Export your configuration
  // after exporting you can use Import function instead of New function.
  // insta, err := goinsta.Import("~/.goinsta")
  // it's useful when you want use goinsta repeatedly.
  insta.Export("~/.goinsta")

  ...
}
```

### Legal

This code is in no way affiliated with, authorized, maintained, sponsored or endorsed by Instagram or any of its affiliates or subsidiaries. This is an independent and unofficial API. Use at your own risk.

### Versioning

Goinsta used gopkg.in as versioning control. Stable new API is the version v3.0. You can get it using:

```bash
$ go get -u -v github.com/Davincible/goinsta
```

Or 

If you have `GO111MODULE=on`

```
$ go get -u github.com/Davincible/goinsta
```

