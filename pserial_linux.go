package pserial

import (
	"golang.org/x/sys/unix"
	"modernc.org/libc/termios"
)

// Linux specific constants.
const (
	tcLflag        = termios.ECHOCTL
	tcIflag        = termios.IUCLC | termios.PARMRK
	tcCMSPAR       = unix.CMSPAR
	tcIUCLC        = unix.IUCLC
	ioctlTcgetattr = unix.TCGETS
	ioctlTcsetattr = unix.TCSETS
	ioctlTcflsh    = unix.TCFLSH
	ioctlTioccbrk  = unix.TIOCCBRK
	ioctlTiocsbrk  = unix.TIOCSBRK
)

func setTermSettingsBaudrate(speed int, settings *unix.Termios) error {
	cflagBaud := linuxBaud(speed)
	if cflagBaud|baudMask != baudMask {
		return errInvalidBaudRate
	}
	// revert old baudrate
	settings.Cflag &^= baudMask

	// set new baudrate
	settings.Cflag |= cflagBaud
	settings.Ispeed = uint32(speed)
	settings.Ospeed = uint32(speed)
	return nil
}

const baudMask = unix.B50 | unix.B75 | unix.B110 |
	unix.B134 | unix.B150 | unix.B200 | unix.B300 |
	unix.B600 | unix.B1200 | unix.B1800 | unix.B2400 |
	unix.B4800 | unix.B9600 | unix.B19200 | unix.B38400 |
	unix.B57600 | unix.B115200 | unix.B230400 | unix.B460800 |
	unix.B500000 | unix.B576000 | unix.B921600 | unix.B1000000 |
	unix.B1152000 | unix.B1500000 | unix.B2000000 | unix.B2500000 |
	unix.B3000000 | unix.B3500000 | unix.B4000000

func linuxBaud(baud int) (CFLAG uint32) {
	switch baud {
	default:
		CFLAG = 0xffff_ffff
	case 50:
		CFLAG = unix.B50
	case 75:
		CFLAG = unix.B75
	case 110:
		CFLAG = unix.B110
	case 134:
		CFLAG = unix.B134
	case 150:
		CFLAG = unix.B150
	case 200:
		CFLAG = unix.B200
	case 300:
		CFLAG = unix.B300
	case 600:
		CFLAG = unix.B600
	case 1200:
		CFLAG = unix.B1200
	case 1800:
		CFLAG = unix.B1800
	case 2400:
		CFLAG = unix.B2400
	case 4800:
		CFLAG = unix.B4800
	case 9600:
		CFLAG = unix.B9600
	case 19200:
		CFLAG = unix.B19200
	case 38400:
		CFLAG = unix.B38400
	case 57600:
		CFLAG = unix.B57600
	case 115200:
		CFLAG = unix.B115200
	case 230400:
		CFLAG = unix.B230400
	case 460800:
		CFLAG = unix.B460800
	case 500000:
		CFLAG = unix.B500000
	case 576000:
		CFLAG = unix.B576000
	case 921600:
		CFLAG = unix.B921600
	case 1000000:
		CFLAG = unix.B1000000
	case 1152000:
		CFLAG = unix.B1152000
	case 1500000:
		CFLAG = unix.B1500000
	case 2000000:
		CFLAG = unix.B2000000
	case 2500000:
		CFLAG = unix.B2500000
	case 3000000:
		CFLAG = unix.B3000000
	case 3500000:
		CFLAG = unix.B3500000
	case 4000000:
		CFLAG = unix.B4000000
	}
	return CFLAG
}

// bauds is a list of all the baud rates supported by linux- could
// be used to find a close match for a given baud rate in future.
var bauds = [...]uint32{
	50,
	75,
	110,
	134,
	150,
	200,
	300,
	600,
	1200,
	1800,
	2400,
	4800,
	9600,
	19200,
	38400,
	57600,
	115200,
	230400,
	460800,
	500000,
	576000,
	921600,
	1000000,
	1152000,
	1500000,
	2000000,
	2500000,
	3000000,
	3500000,
	4000000,
}
