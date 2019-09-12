package main

import (
	"log"

	ya "github.com/art-frela/yadisk"
)

func main() {
	disk := ya.NewYaDisk("https://cloud-api.yandex.net/v1/disk/resources/download", "AgAAAAAntFapAAXeLG-mZsrwakt3vehS76AoR2s")
	err := disk.SaveToDisk("Горы.jpg")
	if err != nil {
		log.Println("save error,", err)
	}
}
