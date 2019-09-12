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
