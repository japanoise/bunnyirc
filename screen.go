package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"gopkg.in/sorcix/irc.v2"
	"strings"
)

var buffer []string

func formatmessage(msg *irc.Message) string {
	switch msg.Command {
	case "JOIN":
		return fmt.Sprintf("%s has joined %s\n", msg.Prefix.Name, msg.Params[0])
	case "PRIVMSG":
		return fmt.Sprintf("%s ─→ %s: %s\n", msg.Prefix.Name, strings.TrimSpace(msg.Params[0]), msg.Params[1])
	case "MODE":
		return fmt.Sprintf("%s sets mode %s\n", msg.Prefix.Name, strings.Join(msg.Params[0:], " "))
	case "NOTICE":
		return fmt.Sprintf("Notice from %s to %s: %s\n", msg.Prefix.Name, msg.Params[0], msg.Params[1])
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
		return fmt.Sprintf("%s has quit (%s)\n", msg.Prefix.Name, msg.Params[0])
	case "CTCP":
		if strings.HasPrefix(msg.Params[1], "ACTION") {
			return fmt.Sprintf("%s: * %s %s\n", msg.Params[0], msg.Prefix.Name, msg.Params[1][7:])
		} else {
			return fmt.Sprintf("CTCP request from %s to %s: %s\n", msg.Prefix.Name, msg.Params[0], msg.Params[1])
		}
	case "CTCPREPLY":
		return fmt.Sprintf("CTCP reply from %s to %s: %s\n", msg.Prefix.Name, msg.Params[0], msg.Params[1])
	default:
		return fmt.Sprint("RAW:", msg.String())
	}
}

/* This function should really be in termbox... */
func drawString(x, y int, str string) {
	for i, runeValue := range str {
		putCh(x+i, y, runeValue)
	}
}

func printstring(str string, anchor, width int) int {
	retval := anchor - (len(str) / width)
	retval--
	y := retval
	clearLine(y, width)
	for i, runeValue := range str {
		if i%width == 0 && i >= width {
			y++
			clearLine(y, width)
		}
		putCh(i%width, y, runeValue)
	}
	return retval
}

func updatescreen() {
	width, height := termbox.Size()
	anchor := height - 2
	i := len(buffer) - 1
	for anchor > 0 && i >= 0 {
		anchor = printstring(buffer[i], anchor, width)
		i--
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
	drawString(0, height-1, ">")
	for {
		drawString(2, height-1, retval)
		termbox.SetCursor(2+cursor, height-1)
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
					retval = retval[0 : len(retval)-1]
					eraseCh(cursor+1, height-1)
					cursor--
				}
			}
		} else if ev.Ch > 31 {
			retval += string(ev.Ch)
			cursor++
		}
	}
}

func sendtobuffer(str string) {
	buffer = append(buffer, StripMircFormatting(str))
	if len(buffer) > cap(buffer)-2 {
		/* Reverse arrays are an ugly hack. This should be a stack. */
		buffer = buffer[1:]
	}
}

func printmsg(msg *irc.Message) {
	sendtobuffer(formatmessage(msg))
}

func initscreen() {
	termbox.Init()
}

func outputloop(client *Client, scrollback int) {
	buffer = make([]string, 0, scrollback)
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
		sendtobuffer(text)
		updatescreen()
		if Command(client, text) {
			return
		}
	}
}
