package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/sorcix/irc.v2"
	"io"
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
		if words[0] == "/q" {
			if len(words) > 1 {
				return fmt.Sprintf("QUIT :%s", strings.Join(words[1:], " ")), true
			} else {
				return "QUIT :Bunnyirc", true
			}
		} else if words[0] == "/t" && len(words) > 1 {
			target = words[1]
			return "", false
		} else if words[0] == "/r" && len(words) > 1 {
			return strings.Join(words[1:], " "), true
		} else if words[0] == "/j" && len(words) > 1 {
			return fmt.Sprintf("JOIN %s", words[1]), true
		} else if words[0] == "/m" && len(words) > 2 {
			return fmt.Sprintf("PRIVMSG %s :%s", words[1], strings.Join(words[2:], " ")), true
		} else if words[0] == "/N" && len(words) > 2 {
			return fmt.Sprintf("NOTICE %s :%s", words[1], strings.Join(words[2:], " ")), true
		} else if words[0] == "/n" && len(words) > 1 {
			return fmt.Sprintf("NICK %s", words[1]), true
		} else if text[1] == '/' {
			return fmt.Sprintf("PRIVMSG %s :%s", target, strings.Replace(text, "/", "", 1)), true
		} else {
			fmt.Printf("Unknown command %s (%d args)\n", words[0], len(words) - 1)
			return "", false
		}
	}
	return fmt.Sprintf("PRIVMSG %s :%s", target, text), target != ""
}

func Command(client *Client, text string) bool {
	send, dosend := Parse(strings.Replace(text, "\n", "", -1))
	if dosend {
		msg := irc.ParseMessage(send)
		if msg != nil {
			client.Send(msg)
			return msg.Command == "QUIT"
		}
	}
	return false
}

func readloop(client *Client) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		if Command(client, text) {
			return
		}
	}
}

func printloop(client *Client) {
	for {
		msg, err := client.Receive()
		if err != nil {
			fmt.Println("Output loop closing:", err)
			return
		}
		printmsg(msg)
		if msg.Command == "ERROR" {
			return
		}
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
}
