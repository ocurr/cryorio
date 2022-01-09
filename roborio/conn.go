package roborio

type Conn interface {
	Close() error
	Exec(string) ([]byte, error)
}

type DialFunc func(user string, pass string, addr string) (Conn, error)
