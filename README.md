# irc in go

>IRC! Internet relay chat! It's how hackers communicate when they don't wanna be overheard.

bunnyirc is an irc client written in Go, inspired by 9front's ircrc. It accepts commands from stdin and outputs to stdout. Unlike its big brother, it can handle basic CTCP. Support for MIRC colours may or may not be coming.

Output takes the form of (barely) formatted messages from the server. Not all messages are handled, so expect to see a raw message once in a while. Despite this simplicity, I think it's a nice and perfectly servicable client.


## Usage of bunnyirc:

  -n string

    	Nickname (defaults to your login name)

  -p int

    	Port to use (default 6667)

  -s string

    	Server to connect to (default "chat.freenode.net")

  -u string

    	Username (defaults to your login name)

  -z	
  
      	Use TLS

## Commands

- /N - Send a NOTICE
- /m - Send a PRIVMSG
- /me - Send a CTCP ACTION
- /j - Join a channel
- /n - Change nick
- /c - Send a CTCP request
- /t - Set the target
- /r - Send a raw irc command
- // - Send a message beginning with a /

If you don't give a command, the input will be sent to the current target as a PRIVMSG.
