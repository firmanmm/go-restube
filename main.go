package main

import (
	"github.com/firmanmm/go-restube/app"
)

func main() {
	restubeInstance := app.NewRestube()
	restubeInstance.Run(":8080")
}
