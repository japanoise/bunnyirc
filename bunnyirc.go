package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/sorcix/irc.v2"
	"regexp"
	"strings"
)

type Client struct {
	Tls     bool
	Details string
	Nick    string
	User    string
	Conn    *irc.Conn
}

type TlsCon struct {
	Usetls   bool
	NoVerify bool
}

func Parse(text string) (string, bool) {
	if text[0] == '/' && len(text) > 1 {
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

func New(tc TlsCon, details, nick, user string) (*Client, error) {
	var ret Client
	var conn *irc.Conn
	var tconn *tls.Conn
	var err error
	if tc.Usetls {
		tconn, err = tls.Dial("tcp", details,
			&tls.Config{InsecureSkipVerify: tc.NoVerify})
		conn = irc.NewConn(tconn)
	} else {
		conn, err = irc.Dial(details)
	}
	ret = Client{tc.Usetls, details, nick, user, conn}
	return &ret, err
}

func (c Client) Send(msg *irc.Message) {
	c.Conn.Encode(msg)
}

func (c Client) Receive() (*irc.Message, error) {
	msg, err := c.Conn.Decode()
	if msg == nil {
		return nil, nil
	} else if msg.Command == "PING" {
		pong := fmt.Sprintf("PONG :%s", msg.Params[0])
		c.Conn.Encode(irc.ParseMessage(pong))
	} else if msg.Command == "PRIVMSG" && msg.Params[1][0] == '\x01' {
		msg.Command = "CTCP"
		msg.Params[1] = strings.Replace(msg.Params[1], "\x01", "", -1)
		if msg.Params[1] == "VERSION" {
			reply := fmt.Sprintf("NOTICE %s :%s", msg.Prefix.Name, "\x01VERSION Bunnyirc (https://github.com/japanoise/bunnyirc)\x01")
			c.Conn.Encode(irc.ParseMessage(reply))
		} else if strings.HasPrefix(msg.Params[1], "PING") {
			reply := fmt.Sprintf("NOTICE %s :\x01%s\x01", msg.Prefix.Name, msg.Params[1])
			c.Conn.Encode(irc.ParseMessage(reply))
		}
	} else if msg.Command == "NOTICE" && msg.Params[1][0] == '\x01' {
		msg.Command = "CTCPREPLY"
		msg.Params[1] = strings.Replace(msg.Params[1], "\x01", "", -1)
	}
	return msg, err
}

func (c Client) Close() {
	c.Close()
}

func (c Client) Authpass(pass string) {
	c.Conn.Encode(&irc.Message{Command: "PASS", Params: []string{pass}})
	c.Conn.Encode(irc.ParseMessage(fmt.Sprintf("NICK %s", c.Nick)))
	c.Conn.Encode(irc.ParseMessage(fmt.Sprintf("USER %s * * :%s", c.User, c.User)))
}

func (c Client) Auth() {
	for {
		msg, _ := c.Conn.Decode()
		if msg.Command == "NOTICE" {
			c.Conn.Encode(irc.ParseMessage(fmt.Sprintf("NICK %s", c.Nick)))
			c.Conn.Encode(irc.ParseMessage(fmt.Sprintf("USER %s * * :%s", c.User, c.User)))
			return
		}
	}
}

func StripMircFormatting(msg string) string {
	mirc := regexp.MustCompile("\x03[0-9]?[0-9]?(,[0-9][0-9]?)?|\x02|\x1D|\x1F|\x16|\x0F")
	return mirc.ReplaceAllLiteralString(msg, "")
}
