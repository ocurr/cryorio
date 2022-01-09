package roborio

import (
	"golang.org/x/crypto/ssh"
)

func DialSsh(user, pass, addr string) (Conn, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return &SshConn{client}, nil
}

type SshConn struct {
	client *ssh.Client
}

func (s *SshConn) Close() error {
	return s.client.Close()
}

func (s *SshConn) Exec(command string) ([]byte, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return []byte{}, err
	}
	out, err := session.CombinedOutput(command)
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}
