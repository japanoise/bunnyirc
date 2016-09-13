package main

import (
	"gopkg.in/sorcix/irc.v2"
	"fmt"
	"flag"
	"os"
	"os/user"
	"bufio"
	"crypto/tls"
	"log"
)

func printmsg(msg *irc.Message){
	switch msg.Command {
	case "PRIVMSG":
		fmt.Printf("%s/%s: %s\n",msg.Prefix.Name,msg.Params[0],msg.Params[1])
	case "NOTICE":
		fmt.Printf("Notice from %s to %s: %s\n",msg.Prefix.Name,msg.Params[0],msg.Params[1])
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
		if msg.Command == "NOTICE" {
			conn.Encode(irc.ParseMessage(fmt.Sprintf("NICK %s",nick)))
			conn.Encode(irc.ParseMessage(fmt.Sprintf("USER %s * * :%s",user,user)))
			return
		}
	}
}

func main() {
	usetls := flag.Bool("z",false,"Use TLS")
	server := flag.String("s","chat.freenode.net","Server to connect to")
	port := flag.Int("p",6667,"Port to use")
	current, _ := user.Current()
	nick := flag.String("n",current.Username,"Nickname")
	user := flag.String("u",current.Username,"Username")
	flag.Parse()
	details := fmt.Sprint(*server,":",*port)
	var conn *irc.Conn
	var err error
	if *usetls {
		tconn, err := tls.Dial("tcp", details, &tls.Config{})
		if err != nil {
			log.Fatalln("Could not connect to IRC server; ", err.Error())
		}
		conn = irc.NewConn(tconn)
	} else {
		conn, err = irc.Dial(details)
		if err != nil {
			log.Fatalln("Could not connect to IRC server; ", err.Error())
		}
	}
	auth(conn,*nick,*user)
	go printloop(conn)
	readloop(conn)
	conn.Close()
}
