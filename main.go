package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/sorcix/irc.v2"
	"log"
	"os"
	"os/user"
)

func printmsg(msg *irc.Message) {
	switch msg.Command {
	case "PRIVMSG":
		fmt.Printf("%s/%s: %s\n", msg.Prefix.Name, msg.Params[0], msg.Params[1])
	case "MODE":
		fmt.Printf("%s sets mode %s on %s\n", msg.Prefix.Name, msg.Params[1], msg.Params[0])
	case "NOTICE":
		fmt.Printf("Notice from %s to %s: %s\n", msg.Prefix.Name, msg.Params[0], msg.Params[1])
	case "001":
		fallthrough
	case "002":
		fallthrough
	case "003":
		fallthrough
	case "372":
		fallthrough
	case "375":
		fallthrough
	case "376":
		fmt.Println(msg.Params[1])
	case "QUIT":
		fmt.Printf("%s has quit (%s)\n", msg.Prefix.Name, msg.Params[0])
	default:
		fmt.Println(msg.String())
	}
}

func readloop(client *Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		msg := irc.ParseMessage(text)
		if msg == nil {
			fmt.Println("Badly formatted message.")
		} else {
			client.Send(msg)
		}
	}
}

func printloop(client *Client) {
	for {
		msg := client.Receive()
		printmsg(msg)
	}
}

func main() {
	current, _ := user.Current()
	nick := flag.String("n", current.Username, "Nickname")
	user := flag.String("u", current.Username, "Username")
	server := flag.String("s", "chat.freenode.net", "Server to connect to")
	port := flag.Int("p", 6667, "Port to use")
	usetls := flag.Bool("z", false, "Use TLS")
	flag.Parse()
	client, err := New(*usetls, fmt.Sprint(*server, ":", *port), *nick, *user)
	if err != nil {
		log.Fatalln("Could not connect to IRC server; ", err.Error())
	}
	/*var target string*/
	client.Auth()
	go printloop(client)
	readloop(client)
	client.Close()
}
