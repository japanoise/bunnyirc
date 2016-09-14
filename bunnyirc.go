package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/sorcix/irc.v2"
)

type Client struct {
	Tls     bool
	Details string
	Nick    string
	User    string
	Conn    *irc.Conn
}

func New(usetls bool, details, nick, user string) (*Client, error) {
	var ret Client
	var conn *irc.Conn
	var tconn *tls.Conn
	var err error
	if usetls {
		tconn, err = tls.Dial("tcp", details, &tls.Config{})
		conn = irc.NewConn(tconn)
	} else {
		conn, err = irc.Dial(details)
	}
	ret = Client{usetls, details, nick, user, conn}
	return &ret, err
}

func (c Client) Send(msg *irc.Message) {
	c.Conn.Encode(msg)
}

func (c Client) Receive() *irc.Message {
	msg, _ := c.Conn.Decode()
	if msg.Command == "PING" {
		pong := fmt.Sprintf("PONG :%s", msg.Params[0])
		c.Conn.Encode(irc.ParseMessage(pong))
	}
	return msg
}

func (c Client) Close() {
	c.Close()
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
