package roborio

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorNoConnection represents a state where a roborio is unavailable at the specified addresses.
var ErrorNoConnection = errors.New("unable to connect to roborio")
var ErrorFileNotExist = errors.New("file does not exist")

type Roborio struct {
	conn     Conn
	dial     DialFunc
	user     string
	pass     string
	team     int
	addrs    []string
	currAddr int
	port     int
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

func Port(port int) Option {
	return func(r *Roborio) {
		r.port = port
	}
}

func NewRoborio(dial DialFunc, user, pass string, options ...Option) (*Roborio, error) {
	rio := &Roborio{port: 22, dial: dial}

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
	conn, err := r.dial(r.user, r.pass, fmt.Sprintf("%s:%d", r.addrs[r.currAddr], r.port))
	if err != nil {
		r.currAddr++
		return fmt.Errorf("unable to connect to roborio at address %s %w", r.addrs[r.currAddr-1], err)
	}

	r.conn = conn
	return nil
}

func (r *Roborio) Disconnect() error {
	return r.conn.Close()
}

func (r *Roborio) Exec(command string) ([]byte, error) {
	return r.conn.Exec(command)
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

func (r *Roborio) BackupFile(src, dest string) error {
	out, _ := r.ListDir()
	if !strings.Contains(string(out), src) {
		return ErrorFileNotExist
	}
	dest = fmt.Sprintf("%s.backup", dest)

	_, err := r.Copy(src, dest)
	if err != nil {
		return fmt.Errorf("unable to backup file: %w", err)
	}
	return nil
}
