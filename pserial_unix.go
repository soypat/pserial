//go:build linux || darwin || freebsd || openbsd

package pserial

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
	"modernc.org/libc/termios"
)

func nativeOpen(name string, config Mode) (*Serial, error) {
	const noblock = unix.O_NONBLOCK | unix.O_NDELAY
	fp, err := os.OpenFile(name, unix.O_RDWR|unix.O_NOCTTY|noblock, 0)
	if err != nil {
		switch err {
		case unix.EBUSY:
			return nil, errors.New("pserial: port busy")
		case unix.EACCES:
			return nil, errors.New("pserial: access denied")
		}
		return nil, err
	}
	fd := int(fp.Fd())
	err = reconfigurePort(fd, config)
	if err != nil {
		fp.Close()
		return nil, err
	}

	err = unix.SetNonblock(fd, true)
	if err != nil {
		fp.Close()
		return nil, err
	}
	err = exclusiveLock(fd)
	if err != nil {
		fp.Close()
		return nil, err
	}
	return &Serial{fd: fd}, nil

}

func reconfigurePort(fd int, config Mode) error {
	attr, err := unix.IoctlGetTermios(fd, ioctlTcgetattr)
	if err != nil {
		return err
	}

	termiosSetRaw(attr)

	err = termiosSetMode(attr, config)
	if err != nil {
		return err
	}

	err = unix.IoctlSetTermios(fd, ioctlTcsetattr, attr)
	if err != nil {
		return err
	}
	// MacOSX requires a special set baudrate if the baudrate is not one of unix baudrates here.
	return nil
}

func termiosSetRaw(attr *unix.Termios) {
	const (
		cflagOR uint32 = unix.CREAD | unix.CLOCAL

		lflagNAND uint32 = termios.ICANON | termios.ECHO | termios.ECHOE |
			termios.ECHOK | termios.ECHONL | termios.ISIG | termios.IEXTEN |
			tcLflag

		oflagNAND uint32 = termios.OPOST | termios.ONLCR | termios.OCRNL

		// bugst uses: brkint, istrp, ignpar, inpck, ixany, ixoff, ixon
		iflagNAND uint32 = termios.INLCR | termios.IGNCR | termios.ICRNL |
			termios.IGNBRK | tcIflag
	)
	attr.Cflag |= cflagOR
	attr.Lflag &^= lflagNAND
	attr.Oflag &^= oflagNAND
	attr.Iflag &^= iflagNAND
}

func exclusiveLock(fd int) error   { return unix.IoctlSetInt(fd, unix.TIOCEXCL, 0) }
func exclusiveUnlock(fd int) error { return unix.IoctlSetInt(fd, unix.TIOCNXCL, 0) }

func termiosSetMode(attr *unix.Termios, mode Mode) error {
	err := setTermSettingsBaudrate(mode.Baud, attr)
	if err != nil {
		return err
	}
	switch mode.ByteSize {
	case 0, 8:
		attr.Cflag |= termios.CS8
	case 7:
		attr.Cflag |= termios.CS7
	case 6:
		attr.Cflag |= termios.CS6
	case 5:
		attr.Cflag |= termios.CS5
	default:
		return errInvalidCharacterSize
	}

	switch mode.StopBits {
	case Stop1:
		attr.Cflag &^= termios.CSTOPB
	case Stop2:
		attr.Cflag |= termios.CSTOPB
	case Stop1_5:
		fallthrough // No unix support for 1.5 stop bits.
	default:
		return errInvalidStopBits
	}

	switch mode.Parity {
	case ParityNone:
		attr.Cflag &^= termios.PARENB | termios.PARODD | tcCMSPAR
	case ParityEven:
		attr.Cflag &^= termios.PARODD | tcCMSPAR
		attr.Cflag |= termios.PARENB
	case ParityOdd:
		attr.Cflag &^= tcCMSPAR
		attr.Cflag |= termios.PARODD | termios.PARENB
	default:
		return errInvalidParity
	}

	// Flow control setup disabled.
	attr.Iflag &^= termios.IXON | termios.IXOFF | termios.IXANY

	// RTS/CTS flow control setup disabled.
	attr.Cflag &^= termios.CRTSCTS

	// Setup VMIN and VTIME.
	attr.Cc[unix.VMIN] = 0
	attr.Cc[unix.VTIME] = 0
	return nil
}
