# yadisk - simple Go package for working with Yandex.disk

(homework for [Geekbrains Go course](https://geekbrains.ru/geek_university/golang), 2nd qrt "Essentials of Golang")

## Purpose

Interaction with resources at the cloud storage Yandex.disk by REST API.  
Implement only two simple functions:

- download private file with OAuth token
- download public file by share URL

## Install

`go get github.com/art-frela/yadisk`

## Example

```golang
package main

import (
	"log"

	ya "github.com/art-frela/yadisk"
)

func main() {
	publicDisk := ya.NewPublicYaDisk("https://cloud-api.yandex.net/v1")
	err := publicDisk.SaveToDisk("https://yadi.sk/i/5q-OcFRvNbE41A", true)
	if err != nil {
		log.Println("save error", err)
	}
}
```