package i2c

import (
	"os"
	"syscall"
)

type I2C struct {
	f *os.File
}

const (
	I2C_SLAVE = 0x0703
)

func Open(bus string, address int) (*I2C, error) {
	fd, err := os.OpenFile(bus, os.O_RDWR, 0600)

	if err != nil {
		return nil, err
	}

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd.Fd(), I2C_SLAVE, uintptr(address), 0, 0, 0)

	if errno != 0 {
		fd.Close()
		return nil, err
	}

	return &I2C{f: fd}, nil
}

func (i2c *I2C) Close() {
	i2c.f.Close()
}

func (i2c *I2C) Write(data []byte) error {
	_, err := i2c.f.Write(data)

	return err
}

func (i2c *I2C) Read(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := i2c.f.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf, nil
}
