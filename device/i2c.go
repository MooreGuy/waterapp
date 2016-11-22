// Package i2c provides low level control over the linux i2c bus.
//
// Before usage you should load the i2c-dev kenel module
//
//      sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// linux systems contain several buses.
package device

import (
	"fmt"
	"github.com/gocql/gocql"
	"os"
	"syscall"
)

const (
	i2c_SLAVE           = 0x0703
	VersionRegister     = 0x05
	UUIDRegisterStart   = 0x06
	NumberUUIDRegisters = 16
)

// I2C represents a connection to an i2c device.
type I2C struct {
	rc         *os.File
	identifier gocql.UUID
}

// New opens a connection to an i2c device.
func New(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), i2c_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	return &I2C{f}, nil
}

// Write sends buf to the remote i2c device. The interpretation of
// the message is implementation dependant.
func (i2c *I2C) Write(buf []byte) (int, error) {
	return i2c.rc.Write(buf)
}

func (i2c *I2C) WriteByte(b byte) (int, error) {
	var buf [1]byte
	buf[0] = b
	return i2c.rc.Write(buf[:])
}

func (i2c *I2C) Read(p []byte) (int, error) {
	return i2c.rc.Read(p)
}

// Writes what register should be read from, waits 10 miliseconds and then
// reads from the i2c device.
func (device *I2C) ReadRegister(readRegister byte) ([]byte, error) {
	device.writeByte(readRegister)
	time.Sleep(time.Millisecond * 10)
	readBuffer := make([]byte, 2, 2)
	read, err := device.Read(readBuffer)
	if err != nil {
		return 0, err
	} else if read != 2 {
		return o, errors.NewError("Didn't read 2 bytes")
	}

	return readBuffer, nil
}

// Gets the stored UUID from the I2C device. This identifier matches up with
// the uuid stored in the database.
func (device *I2C) UUID() (gocql.UUID, error) {
	readBuffer := [16]byte{}
	var i byte
	for i = 0; i < NumUUIDRegisters; i++ {
		readBuffer[i], err = device.ReadRegister(i + UIIDRegisterStart)
		if err != nil {
			return readBuffer, err
		}

	}

	return readBuffer, nil
}

func (device *I2C) WriteUUID(uuid gocql.UUID) error {
	var i byte
	for i = 0; i < NumUUIDRegisters; i++ {
		written, err := device.WriteByte(uuid[i])
		if err != nil || written != 1 {
			return errors.NewError("Couldn't write UUID")
		}
	}

	return nil
}

func (i2c *I2C) Close() error {
	return i2c.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
