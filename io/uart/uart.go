package uart

// An abstraction layer for handling a
// serial interface on Linux

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// UART represents a serial port instance
type UART struct {
	f *os.File
}

const (
	TCDRAIN  = 0x5409
	TCFLSH   = 0x540B
	TCIFLUSH = 0
)

func Open(name string, baud int) (*UART, error) {
	var bauds = map[int]uint32{
		1200:   unix.B1200,
		2400:   unix.B2400,
		4800:   unix.B4800,
		9600:   unix.B9600,
		19200:  unix.B19200,
		38400:  unix.B38400,
		57600:  unix.B57600,
		115200: unix.B115200,
	}

	fd, err := os.OpenFile(name, os.O_RDWR|unix.O_NOCTTY, 0600)

	if err != nil {
		return nil, err
	}

	var termios unix.Termios

	if _, _, err := unix.Syscall(unix.SYS_IOCTL, uintptr(fd.Fd()), uintptr(unix.TCGETS), uintptr(unsafe.Pointer(&termios))); err != 0 {
		return nil, err
	}

	// Disable hardware protocols. Linux dosen't support DTR/DSR
	termios.Cflag &= unix.CRTSCTS

	// No software flow control
	termios.Iflag &^= unix.IXON | unix.IXOFF | unix.IXANY | unix.INLCR | unix.ICRNL | unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.IGNCR

	// Raw input mode
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ICANON | unix.ECHO | unix.ECHOE | unix.ISIG

	// Enable receiver,Ignore Modem Control lines
	termios.Cflag |= unix.CREAD | unix.CLOCAL

	// Set timeouts
	termios.Cc[unix.VMIN] = 0
	termios.Cc[unix.VTIME] = 5

	// Set N,8,1
	termios.Cflag &^= unix.PARENB
	termios.Cflag &^= unix.CSTOP
	termios.Cflag &^= unix.CSIZE
	termios.Cflag |= unix.CS8

	// Set baud rate
	rate, ok := bauds[baud]

	if !ok {
		return nil, fmt.Errorf("unrecognized baud rate")
	}

	termios.Ispeed = rate
	termios.Ospeed = rate

	if _, _, err := syscall.Syscall(unix.SYS_IOCTL, uintptr(fd.Fd()), unix.TCSETS, uintptr(unsafe.Pointer(&termios))); err != 0 {
		return nil, err
	}

	syscall.Syscall(unix.SYS_IOCTL, uintptr(fd.Fd()), unix.TCIOFLUSH, uintptr(unsafe.Pointer(nil)))

	return &UART{f: fd}, err
}

func (uart *UART) Close() error {
	return uart.f.Close()
}

func (uart *UART) WriteBytes(data []byte) error {
	nn, err := uart.f.Write(data)

	if nn != len(data) {
		err = errors.New("error to transmit data")
	}

	// Wait until the trasnmit end
	syscall.Syscall(syscall.SYS_IOCTL, uart.f.Fd(), TCDRAIN, 0)

	return err
}

func (uart *UART) WriteByte(data byte) error {
	return uart.WriteBytes([]byte{data})
}

func (uart *UART) WriteString(s string) error {
	return uart.WriteBytes([]byte(s))
}

func (uart *UART) ReadBytes(data []byte) (int, error) {
	return uart.f.Read(data)
}

func (uart *UART) Flush() {
	fd := uart.f.Fd()
	syscall.Syscall(syscall.SYS_IOCTL, fd, TCFLSH, uintptr(TCIFLUSH))
}
