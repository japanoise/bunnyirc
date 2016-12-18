package main

import (
	"container/list"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"gopkg.in/sorcix/irc.v2"
	"strings"
	"unicode/utf8"
)

var buffer *list.List
var limit int

func formatmessage(msg *irc.Message) string {
	prefix := msg.Prefix
	if prefix == nil {
		prefix = &irc.Prefix{"unknown", "unknown", "unknown"}
	}
	switch msg.Command {
	case "JOIN":
		return fmt.Sprintf("%s has joined %s\n", prefix.Name, msg.Params[0])
	case "PRIVMSG":
		return fmt.Sprintf("(%s) %s: %s\n", msg.Params[0], prefix.Name, msg.Params[1])
	case "MODE":
		return fmt.Sprintf("%s sets mode %s\n", prefix.Name, strings.Join(msg.Params[0:], " "))
	case "NOTICE":
		return fmt.Sprintf("Notice from %s to %s: %s\n", prefix.Name, msg.Params[0], msg.Params[1])
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
		return fmt.Sprintln(msg.Params[1])
	case "QUIT":
		return fmt.Sprintf("%s has quit (%s)\n", prefix.Name, msg.Params[0])
	case "CTCP":
		if strings.HasPrefix(msg.Params[1], "ACTION") {
			return fmt.Sprintf("%s: * %s %s\n", msg.Params[0], prefix.Name, msg.Params[1][7:])
		} else {
			return fmt.Sprintf("CTCP request from %s to %s: %s\n", prefix.Name, msg.Params[0], msg.Params[1])
		}
	case "CTCPREPLY":
		return fmt.Sprintf("CTCP reply from %s to %s: %s\n", prefix.Name, msg.Params[0], msg.Params[1])
	default:
		return fmt.Sprint("RAW:", msg.String())
	}
}

/* This function should really be in termbox... */
func drawString(x, y int, str string) {
	i := 0
	for _, runeValue := range str {
		putCh(x+i, y, runeValue)
		i += runewidth.RuneWidth(runeValue)
	}
}

func printstring(str string, anchor, width int) int {
	retval := anchor - (len(str) / width)
	retval--
	y := retval
	clearLine(y, width)
	i := 0
	for _, runeValue := range str {
		if i%width == 0 && i >= width {
			y++
			clearLine(y, width)
		}
		putCh(i%width, y, runeValue)
		i += runewidth.RuneWidth(runeValue)
	}
	return retval
}

func updatescreen() {
	width, height := termbox.Size()
	anchor := height - 2
	i := buffer.Front()
	for anchor > 0 && i != nil {
		anchor = printstring(i.Value.(string), anchor, width)
		i = i.Next()
	}
	termbox.Flush()
}

func clearLine(y, width int) {
	for i := 0; i < width; i++ {
		eraseCh(i, y)
	}
}

func putCh(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
}

func eraseCh(x, y int) {
	putCh(x, y, ' ')
}

func GetString() string {
	width, height := termbox.Size()
	clearLine(height-1, width)
	retval := ""
	cursor := 0
	tlen := len(target) + 3
	drawString(0, height-1, fmt.Sprint(target, " >"))
	for {
		drawString(tlen, height-1, retval)
		termbox.SetCursor(tlen+cursor, height-1)
		termbox.Flush()
		ev := termbox.PollEvent()
		if ev.Ch == 0 {
			switch ev.Key {
			case termbox.KeySpace:
				retval += " "
				cursor++
			case termbox.KeyEnter:
				if len(retval) > 0 {
					termbox.HideCursor()
					return retval
				}
			case termbox.KeyDelete:
				fallthrough
			case termbox.KeyBackspace:
				if cursor > 0 {
					r, rs :=
						utf8.DecodeLastRuneInString(retval)
					retval = retval[0 : len(retval)-rs]
					eraseCh(cursor+(tlen)-runewidth.RuneWidth(r), height-1)
					cursor -= runewidth.RuneWidth(r)
				}
			}
		} else if ev.Ch > 31 {
			retval += string(ev.Ch)
			cursor += runewidth.RuneWidth(ev.Ch)
		}
	}
}

func sendtobuffer(str string) {
	buffer.PushFront(StripMircFormatting(str))
	if buffer.Len() > limit {
		buffer.Remove(buffer.Back())
	}
}

func printmsg(msg *irc.Message) {
	sendtobuffer(formatmessage(msg))
}

func initscreen() {
	termbox.Init()
}

func outputloop(client *Client, scrollback int) {
	limit = scrollback
	buffer = list.New()
	for {
		msg, err := client.Receive()
		if err != nil {
			fmt.Println("Output loop closing:", err)
			return
		}
		if msg != nil {
			if msg.Command != "PING" {
				printmsg(msg)
			}
			if msg.Command == "ERROR" {
				return
			}
		}
		updatescreen()
	}
}

func inputloop(client *Client) {
	defer termbox.Close()
	for {
		text := GetString()
		sendtobuffer(fmt.Sprint("(", target, ") ", text))
		updatescreen()
		if Command(client, text) {
			return
		}
	}
}
