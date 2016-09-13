package main

import (
	"gopkg.in/sorcix/irc.v2"
	"fmt"
	"os"
	"bufio"
	"crypto/tls"
	"log"
)

func printmsg(msg *irc.Message){
	switch msg.Command {
	case "PRIVMSG":
		fmt.Printf("%s/%s: %s\n",msg.Prefix.Name,msg.Params[0],msg.Params[1])
	case "QUIT":
		fmt.Printf("%s has quit (%s)\n",msg.Prefix.Name,msg.Params[0])
	default:
		fmt.Println(msg.String())
	}
}

func readloop(conn *irc.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		msg := irc.ParseMessage(text)
		if msg == nil {
			fmt.Println("Badly formatted message.")
		} else {
			conn.Encode(msg)
		}
	}
}

func printloop(conn *irc.Conn) {
	for {
		msg, _ := conn.Decode()
		printmsg(msg)
		if msg.Command == "PING" {
			pong := fmt.Sprintf("PONG :%s",msg.Params[0])
			fmt.Println(pong)
			conn.Encode(irc.ParseMessage(pong))
		}
	}
}

func auth(conn *irc.Conn, nick, user string) {
	for {
		msg, _ := conn.Decode()
		printmsg(msg)
		if msg.Params[0] == "AUTH" {
			conn.Encode(irc.ParseMessage(fmt.Sprintf("NICK %s",nick)))
			conn.Encode(irc.ParseMessage(fmt.Sprintf("USER %s * * :%s",user,user)))
			return
		}
	}
}

func main() {
	usetls := false
	server := "irc.rizon.net"
	port := "6660"
	nick := "Tewimeleon"
	user := "Tewi"
	details := fmt.Sprintf("%s:%s",server,port)
	var conn *irc.Conn
	var err error
	if usetls {
		tconn, err := tls.Dial("tcp", details, &tls.Config{})
		if err != nil {
			log.Fatalln("Could not connect to IRC server")
		}
		conn = irc.NewConn(tconn)
	} else {
		conn, err = irc.Dial(details)
		if err != nil {
			log.Fatalln("Could not connect to IRC server")
		}
	}
	auth(conn,nick,user)
	go printloop(conn)
	readloop(conn)
	conn.Close()
}
