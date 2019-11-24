package main

import TToT "github.com/AliceDiNunno/TwitterToTelegram"

func main() {
	TToT.OpenDatabase()
	TToT.CreateTwitterClient()
	TToT.SetupBot()
}
