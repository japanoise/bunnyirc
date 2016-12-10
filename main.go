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
	fmt.Println(msg.String())
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
		} else if words[0] == "/c" && len(words) > 2 {
			return fmt.Sprintf("PRIVMSG %s :\x01%s\x01", words[1], strings.Join(words[2:], " ")), true
		} else if words[0] == "/N" && len(words) > 2 {
			return fmt.Sprintf("NOTICE %s :%s", words[1], strings.Join(words[2:], " ")), true
		} else if words[0] == "/n" && len(words) > 1 {
			return fmt.Sprintf("NICK %s", words[1]), true
		} else if words[0] == "/me" && len(words) > 1 {
			return fmt.Sprintf("PRIVMSG %s :\x01ACTION %s\x01", target, strings.Join(words[1:], " ")), true
		} else if text[1] == '/' {
			return fmt.Sprintf("PRIVMSG %s :%s", target, strings.Replace(text, "/", "", 1)), true
		} else {
			fmt.Printf("Unknown command %s (%d args)\n", words[0], len(words)-1)
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

func inputloop(client *Client) {
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

func outputloop(client *Client) {
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
	pass := flag.String("P", "", "Connection Password")
	user := flag.String("u", current.Username, "Username")
	server := flag.String("s", "chat.freenode.net", "Server to connect to")
	port := flag.Int("p", 6667, "Port to use")
	usetls := flag.Bool("z", false, "Use TLS")
	noverify := flag.Bool("v", false, "Skip TLS connection verification")
	flag.Parse()
	client, err := New(TlsCon{*usetls, *noverify},
		fmt.Sprint(*server, ":", *port), *nick, *user)
	if err != nil {
		log.Fatalln("Could not connect to IRC server; ", err.Error())
	}
	fmt.Println("Ok, let's auth!")
	if *pass == "" {
		client.Auth()
	} else {
		client.Authpass(*pass)
	}
	go outputloop(client)
	inputloop(client)
}
