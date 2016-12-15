package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
)

var target string

func main() {
	current, _ := user.Current()
	nick := flag.String("n", current.Username, "Nickname")
	pass := flag.String("P", "", "Connection Password")
	user := flag.String("u", current.Username, "Username")
	server := flag.String("s", "chat.freenode.net", "Server to connect to")
	port := flag.Int("p", 6667, "Port to use")
	usetls := flag.Bool("z", false, "Use TLS")
	noverify := flag.Bool("v", false, "Skip TLS connection verification")
	scrollback := flag.Int("S", 10000, "Number of messages to keep in scrollback")
	flag.Parse()
	client, err := New(TlsCon{*usetls, *noverify},
		fmt.Sprint(*server, ":", *port), *nick, *user)
	if err != nil {
		log.Fatalln("Could not connect to IRC server; ", err.Error())
	}
	if *pass == "" {
		client.Auth()
	} else {
		client.Authpass(*pass)
	}
	initscreen()
	go outputloop(client, *scrollback)
	inputloop(client)
}
