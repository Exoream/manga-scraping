package main

import "manga/route"

func main() {
	r := route.SetupRouter()
	r.Run(":8080")
}
