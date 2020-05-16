package main

import "github.com/gavinc95/go-blog/db"

func main() {
	app := NewApp(":8010", &db.GenID{})
	app.Run()
}
