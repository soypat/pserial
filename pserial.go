package pserial

import (
	"errors"
	"sync"

	"go.bug.st/serial"
)

var (
	errInvalidCharacterSize = errors.New("pserial: invalid or unsupported character size")
	errInvalidStopBits      = errors.New("pserial: invalid or unsupported stop bits")
	errInvalidParity        = errors.New("pserial: invalid or unsupported parity")
	errInvalidBaudRate      = errors.New("pserial: invalid or unsupported baud rate")
)

type (
	Parity   uint8
	StopBits uint8
)

const (
	Stop1 = iota
	Stop2
	Stop1_5 // 1.5 stop bits. Only available on some platforms.
)

const (
	ParityNone = iota
	ParityOdd
	ParityEven
	ParityMark
)

type Serial struct {
	fd int
	mu sync.Mutex
}

type Mode struct {
	// Baud is amount of data and non-data bits sent over the wire per second.
	// Common bauds are 9600, 19200, 115200.
	Baud int
	//
	ByteSize uint8
	Parity   Parity
	StopBits StopBits
}

func Open(name string, config Mode) (*Serial, error) {
	if config.Baud <= 0 {
		return nil, errors.New("pserial: invalid baud rate")
	}
	serial.Open("", nil)
	return nativeOpen(name, config)
}

func (s *Serial) lock() {
	s.mu.Lock()
}
func (s *Serial) unlock() {
	s.mu.Unlock()
}
