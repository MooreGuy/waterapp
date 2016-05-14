package device

import (
	"log"
)

// Gathers all of the devices connected to this controller.
// Right now devices are just I2C devices, because that is the only interface
// method supported currently, but i2C can be abstracted to device, where the
// communication method would be agnostic.
// Send over a byte to signal that the controller would like to see if the i2c
// device is something that it can talk to.
// TODO: Change this to a async method that returns a channel
func findDevices() []*I2C {
	devices := []*I2C{}

	for start := 0x20; start < 0x60; start = start + 0x1 {
		i2cDevice, err := New(0x26, 1)
		if err != nil {
			continue
		}
		defer i2cDevice.Close()

		_, err = i2cDevice.WriteByte(0x01)
		if err != nil {
			continue
		}

		var buf []byte = make([]byte, 1)
		_, err = i2cDevice.Read(buf)
		if err != nil {
			continue
		}

		reader := bytes.NewReader(buf)
		response, err := binary.ReadVarint(reader)
		if response != 0xFF {
			continue
		}

		devices = append(devices, i2cDevice)
	}

	return devices
}

func sendTestMessage(devices []*I2C) {
	for x := 0; x < len(devices); x++ {
		curDevice := devices[x]
		curDevice.WriteByte('A')
	}
}
