[![Build Status](https://secure.travis-ci.org/japanoise/bunnyirc.png)](http://travis-ci.org/japanoise/bunnyirc)

# irc in go

>IRC! Internet relay chat! It's how hackers communicate when they don't wanna be overheard.

bunnyirc is an irc client written in Go, inspired by 9front's ircrc. All output
goes to one place - the output window - implemented using Termbox.

## Features

- Connects to irc.
- Sends messages.
- Recieves messages.
- Some CTCP features (action, ping)

## Possible future features

- Nick colors. This wouldn't be too difficult to do. It improves the reading
  experience as you can visually differentiate people quicker.
- Timestamps. You could previously do this with pipes, but not any more!
  Shouldn't be too hard to implement, however.

## Anti-features

- Buffers by window. Would cause code complexity. Besides, you'd end up
  spending a lot of time switching between buffers. Having them all
  together is actually very pleasant, as you can see messages as they come
  from multiple buffers.
- DCC - this is potentially dangerous, and an abuse of the irc protocol.
- mIRC-style formatting - No.
- Multiple servers - This would clutter the output. A suggested solution is to
  use a screenrc:

~~~
# Allow real scrolling
termcapinfo xterm* ti@:te@

# No startup messages
startup_message off

# Nice status line
hardstatus off
hardstatus alwayslastline
hardstatus string '%{= kG}[ %{G}%H %{g}][%= %{= kw}%?%-Lw%?%{r}(%{W}%n*%f%t%?(%u)%?%{r})%{w}%?%+Lw%?%?%= %{g}][%{B} %m-%d %{W} %c %{g}]'

# Set your servers up here, e.g.
# screen bunnyirc -s 'chat.server.org'
# title "server.org"
~~~

## Usage of bunnyirc:

  -n string

    	Nickname (defaults to your login name)

  -p int

    	Port to use (default 6667)

  -P int

    	Connection password

  -s string

    	Server to connect to (default "chat.freenode.net")

  -u string

    	Username (defaults to your login name)

  -z

      	Use TLS

  -v

      	Skip TLS connection verification

## Commands

- /N - Send a NOTICE
- /m - Send a PRIVMSG
- /me - Send a CTCP ACTION
- /j - Join a channel
- /q - Send QUIT with optional reason
- /n - Change nick
- /c - Send a CTCP request
- /t - Set the target
- /r - Send a raw irc command
- // - Send a message beginning with a /

If you don't give a command, the input will be sent to the current target as a
PRIVMSG.
