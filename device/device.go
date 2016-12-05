package device

import (
	"github.com/MooreGuy/waterapp/network"
	"github.com/gocql/gocql"
	"log"
	"math/big"
	"time"
)

const (
	ValidVersion = 1
)

type Sensor interface {
	Read() int
	UUID() gocql.UUID
}

type FakeSensor struct {
	uuid gocql.UUID
}

func NewFakeSensor() Sensor {
	return FakeSensor{gocql.TimeUUID()}
}

func (s FakeSensor) UUID() gocql.UUID {
	return s.uuid
}

func (s FakeSensor) Read() int {
	return 1
}

func ReadSensors(messageQueue chan network.Message) {
	fakeSensor := NewFakeSensor()
	for {
		mes := map[string]interface{}{
			"data":     fakeSensor.Read(),
			"sensorid": fakeSensor.UUID(),
		}

		log.Println("Sending fake sensor data.")
		messageQueue <- mes
		time.Sleep(time.Second * 1)
	}
}

func HandleDeviceSignal() (commandQueue chan Command) {
	commandQueue = make(chan Command, 100)
	go func(commandQueue chan Command) {
		for {
			command := <-commandQueue
			devices := GetFakeDevices()
			devices.sendCommand(command)
		}
	}(commandQueue)

	return
}

type Device struct {
	uuid gocql.UUID
}

func (d Device) UUID() gocql.UUID {
	return d.uuid
}

type Command struct {
	Name   string
	Data   int
	Target gocql.UUID
}

type DeviceCollection map[gocql.UUID]Device

func (col DeviceCollection) sendCommand(c Command) {
	_, ok := col[c.Target]
	if !ok {
		log.Println("Couldn't find device")
		log.Println(c.Target)
	}

	log.Println("This is where the device should be given the command.", c.Data)
}

func GetFakeDevices() DeviceCollection {
	uuid := gocql.TimeUUID()
	return DeviceCollection{uuid: Device{uuid}}
}

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
		if read, err := device.ReadRegister(VersionRegister); err == nil {
			i := new(big.Int)
			version := int32(i.SetBytes(read).Int64())
			if version == ValidVersion {
				functioningDevices = append(functioningDevices, device)
			}
		}
	}

	return functioningDevices
}
