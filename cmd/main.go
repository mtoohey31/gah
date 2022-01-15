package main

import (
	"fmt"

	"mtoohey.com/gah"
)

func main() {
	gah.Cmd{
		Name:        "gah",
		Author:      "Matthew Toohey <contact@mtoohey.com>",
		Version:     "v0.1.1",
		Description: "Go Argument Handler",
		Content: func(f struct {
			FireStation string `long:"fires" short:"5"`
			Test        int    `minVal:"5"`
			Yeet        bool   `short:"y" invert:""`
			Root        string `path:""`
		}, a struct {
			Paths   []string `path:"" min:"1" max:"2"`
			Strings []string `max:"0"`
			Test    int
		}) {
			println("in function")
			println(f.FireStation)
			println(f.Test)
			println(f.Yeet)
			println(f.Root)
			fmt.Println(a.Paths)
			fmt.Println(a.Strings)
		},
		// Content: []gah.Cmd{
		// 	{
		// 		Name:        "stuff",
		// 		Description: "interesting",
		// 		Content: func(_ struct{}, _ struct{}) {
		// 			println("hmm")
		// 		},
		// 	},
		// 	{
		// 		Name:        "bruh",
		// 		Description: "interesting",
		// 		Content: func(_ struct{}, _ struct{}) {
		// 			println("bruh moment")
		// 		},
		// 	},
		// },
	}.SimpleEval()
}
