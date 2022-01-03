package roborio

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// ErrorNoConnection represents a state where a roborio is unavailable at the specified addresses.
var ErrorNoConnection = errors.New("unable to connect to roborio")

type Roborio struct {
	config   *ssh.ClientConfig
	client   *ssh.Client
	team     int
	addrs    []string
	currAddr int
}

type Option func(*Roborio)

func Addresses(addrs ...string) Option {
	return func(r *Roborio) {
		r.addrs = addrs
	}
}

func Team(team int) Option {
	return func(r *Roborio) {
		r.team = team
	}
}

func NewRoborio(user, pass string, options ...Option) (*Roborio, error) {
	rio := &Roborio{}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	rio.config = config

	for _, option := range options {
		option(rio)
	}

	if rio.team == 0 {
		var err error
		rio.team, err = GetTeamNumber()
		if err != nil {
			return nil, err
		}
	}

	if rio.addrs == nil {
		rio.addrs = GetAddresses(rio.team)
	}

	return rio, nil
}

// Connect attempts to open a connection to the roborio using the list of addresses.
// If Connect fails to open a connection to the roborio it will cycle to the next address
// in the list on the next call to Connect. If there are no more addresses it returns ErrorNoConnection.
func (r *Roborio) Connect() error {
	if r.currAddr > len(r.addrs) {
		return ErrorNoConnection
	}
	client, err := ssh.Dial("tcp", r.addrs[r.currAddr]+":22", r.config)
	if err != nil {
		r.currAddr++
		return fmt.Errorf("unable to connect to roborio at address %s %w", r.addrs[r.currAddr-1], err)
	}

	r.client = client
	return nil
}

func (r *Roborio) Disconnect() error {
	return r.client.Close()
}

func (r *Roborio) Exec(command string) ([]byte, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return []byte{}, err
	}
	out, err := session.CombinedOutput(command)
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

func (r *Roborio) ListDir() ([]byte, error) {
	return r.Exec("ls")
}

func (r *Roborio) Copy(origin, dest string) ([]byte, error) {
	return r.Exec(fmt.Sprintf("cp %s %s", origin, dest))
}

func (r *Roborio) Touch(name string) ([]byte, error) {
	return r.Exec("touch " + name)
}

func (r *Roborio) Remove(pattern string) ([]byte, error) {
	return r.Exec("rm " + pattern)
}
