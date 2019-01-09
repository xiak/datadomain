package services

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

var (
	ErrMsgConnFailed 	= "Failed to create ssh connection"
	ErrMsgSessFailed 	= "Failed to create ssh session"
	ErrMsgReqPtyFailed 	= "Failed to request pty"
)

type Command struct {
	Protocol	string
	Host 		string
	Port 		int
	User 		string
	Password 	string
	Client      *ssh.Client
	Result      string
	Prefix      string
}

func NewCommand(host string, user string, password string) *Command {
	if host == "" {
		host = "127.0.0.1"
	}
	return &Command{
		Protocol: "tcp",
		Host: host,
		Port: 22,
		User: user,
		Password: password,
		Prefix: "[ Command]:",
	}
}

func (c *Command) sshInteractive(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	answers = make([]string, len(questions))
	for n, _ := range questions {
		answers[n] = c.Password
	}

	return answers, nil
}

func (c *Command) connect() (error) {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	fmt.Println("Connecting to remote server")
	client, err := ssh.Dial(c.Protocol, addr, &ssh.ClientConfig{
		User: c.User,
		// AuthMethods:
		// 1. Password authentication method: (Need set PasswordAuthentication yes in /etc/ssh/sshd_config)
		// Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		// 2. Keyboard interactive authentication method
		Auth: []ssh.AuthMethod{ssh.KeyboardInteractive(c.sshInteractive)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	c.Client = client
	fmt.Println("Connected to remote server:", addr)
	return err
}

func (c *Command) Run(cmd string, a ...interface{}) (string, error) {
	if c.Client == nil {
		if err := c.connect(); err != nil {
			return ErrMsgConnFailed, err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return ErrMsgSessFailed, err
	}
	defer session.Close()
	cmdString := fmt.Sprintf(cmd, a ...)
	fmt.Println(c.Prefix, cmdString)
	buf, err := session.CombinedOutput(cmdString)

	c.Result = string(buf)
	return c.Result, err
}