package main

import "adinunno.fr/twitter-to-telegram/src"

func main() {
	src.OpenDatabase()
	src.CreateTwitterClient()
	src.SetupBot()
}
