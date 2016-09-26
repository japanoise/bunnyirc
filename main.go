package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/sorcix/irc.v2"
	"log"
	"os"
	"os/user"
	"strings"
)

var target string

func printmsg(msg *irc.Message) {
	switch msg.Command {
	case "JOIN":
		fmt.Printf("%s has joined %s\n", msg.Prefix.Name, msg.Params[0])
	case "PRIVMSG":
		fmt.Printf("%s/%s: %s\n", msg.Prefix.Name, msg.Params[0], msg.Params[1])
	case "MODE":
		fmt.Printf("%s sets mode %s\n", msg.Prefix.Name, strings.Join(msg.Params[0:], " "))
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
		fmt.Println("RAW:", msg.String())
	}
}

func Parse(text string) (string, bool) {
	if text[0] == '/' {
		words := strings.Split(text, " ")
		if words[0] == "/t" && len(words) > 1 {
			target = words[1]
			return "", false
		}
		if words[0] == "/r" && len(words) > 1 {
			return strings.Join(words[1:], " "), true
		}
		if words[0] == "/j" && len(words) > 1 {
			return fmt.Sprintf("JOIN %s", words[1]), true
		}
		if words[0] == "/m" && len(words) > 2 {
			return fmt.Sprintf("PRIVMSG %s :%s", words[1], strings.Join(words[2:], " ")), true
		}
		if words[0] == "/n" && len(words) > 2 {
			return fmt.Sprintf("NOTICE %s :%s", words[1], strings.Join(words[2:], " ")), true
		}
	}
	return strings.Replace(fmt.Sprintf("PRIVMSG %s :%s", target, text), "\n", "", -1), target != ""
}

func Command(client *Client, text string) {
	send, dosend := Parse(text)
	if dosend {
		msg := irc.ParseMessage(send)
		if msg != nil {
			client.Send(msg)
		}
	}
}

func readloop(client *Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		Command(client, text)
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
	client.Auth()
	go printloop(client)
	readloop(client)
	client.Close()
}
