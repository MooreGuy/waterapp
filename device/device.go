package device

import (
	"errors"
	//	"bytes"
	//	"encoding/binary"
)

const (
	ValidVersion = 1
)

// TODO: These should all be async, really gross that they're not.

// Gathers all of the devices connected to this controller.
// Right now devices are just I2C devices, because that is the only interface
// method supported currently, but i2C can be abstracted to device, where the
// communication method would be agnostic.
// Send over a byte to signal that the controller would like to see if the i2c
// device is something that it can talk to.
func FindDevices() []*I2C {
	devices := []*I2C{}

	for i := 0x0; i < 0x7F; i = i + 0x1 {
		i2cDevice, err := New(uint8(i), 1)
		if err != nil {
			continue
		}
		defer i2cDevice.Close()

		_, err = i2cDevice.WriteByte(0x01)
		if err == nil {
			devices = append(devices, i2cDevice)
		}

	}

	return devices
}

func FindFunctioningDevices(deviceList []*I2C) []*I2C {
	functioningDevices := []*I2C{}
	for _, device := range deviceList {
		if read, err := device.ReadDevice(VersionRegister); err == nil &&
			read == ValidVersion {
			functioningDevices = append(functioningDevices, device)
		}
	}

	return functioningDevices
}
