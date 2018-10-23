package ble

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

const timeout = 2 * time.Second
const CommandCharacteristic = "88881311deadbea71523785feab7e123"
const DataCharacteristic = "88881312deadbea71523785feab7e123"

type VehDeviceConnection struct {
	client      *ClientAdaptor
	commandResp []byte
	commandSync chan bool
	dataItem    []byte
	dataSync    chan bool
}

type VehInfo struct {
	CurrentTime     uint32 `json:"currentTime"`
	CurrentSequence uint16 `json:"currentSequence"`
	RecordCount     uint16 `json:"recordCount"`
}

type VehConfig struct {
	AdvertisingInterval uint16 `json:"advertisingInterval"`
	SampleTime          uint16 `json:"sampleTime"`
}

type VehEvents struct {
	Timestamp      uint32 `json:"timestamp"`
	SequenceNumber uint16 `json:"sequenceNumber"`
	Event          byte   `json:"event"`
	Length         byte   `json:"length"`
	Data           []byte `json:"data"`
}

func NewVehDeviceConnection(address string) (*VehDeviceConnection, error) {
	vdc := &VehDeviceConnection{
		client:      NewClientAdaptor(address),
		commandSync: make(chan bool, 1),
		dataSync:    make(chan bool, 1),
	}

	err := vdc.client.Connect()
	if err == nil {
		err = vdc.client.Subscribe(CommandCharacteristic, true, func(d []byte, e error) {
			if e != nil {
				log.Fatalf("Command Characteristic error %v", e)
			}
			vdc.commandResp = d
			vdc.commandSync <- true
		})

		if err == nil {
			err = vdc.client.Subscribe(DataCharacteristic, false, func(d []byte, e error) {
				if e != nil {
					log.Fatalf("Data Characteristic error %v", e)
				}
				vdc.dataItem = d
				vdc.dataSync <- true
			})
		}
	}
	return vdc, err
}

func (vdc *VehDeviceConnection) Finalize() error {
	return vdc.client.Finalize()
}

func (vdc *VehDeviceConnection) msgAndWait(msg []byte) ([]byte, error) {
	err := vdc.client.WriteCharacteristic(CommandCharacteristic, msg)
	if err != nil {
		return nil, err
	}
	select {
	case <-vdc.commandSync:
		return vdc.commandResp, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("Timeout after %d seconds", timeout)
	}
}

func (vdc *VehDeviceConnection) VehUnlock() error {
	successful := []byte{0x00, 0x00}
	data, err := vdc.msgAndWait([]byte{0x00, 0x41, 0x42, 0x43})
	if err != nil {
		return fmt.Errorf("Unlock Error %v", err)
	}

	if !bytes.Equal(data, successful) {
		return fmt.Errorf("Unlock Error %v returned", data)
	}
	return nil
}

func (vdc *VehDeviceConnection) vehGenericGet(cmd []byte) ([]byte, error) {
	data, err := vdc.msgAndWait(cmd)
	if err != nil {
		return nil, fmt.Errorf("Command %02X Error %v", cmd[0], err)
	}
	if data[0] != cmd[0] {
		return nil, fmt.Errorf("Command Error, response did not start with %02X response: was %v", cmd[0], data)
	}
	return data, nil
}

func (vdc *VehDeviceConnection) VehGetVersion() (string, error) {
	data, err := vdc.vehGenericGet([]byte{0x01})
	if err != nil {
		return "", fmt.Errorf("GetVersion Error %v", err)
	}
	return fmt.Sprintf("S/W: 0x%02X.%02X Protocol: 0x%02X H/W: 0x%02X", data[3], data[2], data[4], data[5]), nil
}

func (vdc *VehDeviceConnection) VehGetInfo() (VehInfo, error) {
	vi := VehInfo{}
	data, err := vdc.vehGenericGet([]byte{0x02})
	if err != nil {
		return vi, fmt.Errorf("VehGetInfo Error %v", err)
	}
	vi.CurrentTime = binary.LittleEndian.Uint32(data[2:6])
	vi.CurrentSequence = binary.LittleEndian.Uint16(data[6:8])
	vi.RecordCount = binary.LittleEndian.Uint16(data[8:10])
	return vi, nil
}

func (vdc *VehDeviceConnection) VehGetConfig() (VehConfig, error) {
	vc := VehConfig{}
	data, err := vdc.vehGenericGet([]byte{0x03})
	if err != nil {
		return vc, fmt.Errorf("VehGetConfig Error %v", err)
	}
	vc.AdvertisingInterval = binary.LittleEndian.Uint16(data[2:4])
	vc.SampleTime = binary.LittleEndian.Uint16(data[4:6])
	return vc, nil
}

func (vdc *VehDeviceConnection) VehSetConfig(c VehConfig) error {
	// []byte{0x04}
	return fmt.Errorf("Not implemented")
}

func (vdc *VehDeviceConnection) VehGetFriendlyName() (string, error) {
	data, err := vdc.vehGenericGet([]byte{0x05})
	if err != nil {
		return "", fmt.Errorf("VehGetConfig Error %v", err)
	}
	return string(data[2:]), nil
}

func (vdc *VehDeviceConnection) VehSetFriendlyName(n string) error {
	// []byte{0x06}
	return fmt.Errorf("Not implemented")
}

func (vdc *VehDeviceConnection) VehReadEvents(sequence uint16) ([]VehEvents, error) {
	data, err := vdc.vehGenericGet([]byte{0x07, 0x00, 0x00})
	if err != nil {
		return nil, fmt.Errorf("VehReadEvents Error %v", err)
	}
	num := binary.LittleEndian.Uint16(data[2:4])
	fmt.Println("VehReadEvents returned num\t", num)
	events := make([]VehEvents, num)
	itemCount := uint16(0)

	for num > itemCount {
		select {
		case <-vdc.dataSync:
			ve := VehEvents{
				Timestamp:      binary.LittleEndian.Uint32(vdc.dataItem[0:4]),
				SequenceNumber: binary.LittleEndian.Uint16(vdc.dataItem[4:6]),
				Event:          vdc.dataItem[6],
				Length:         vdc.dataItem[7],
				Data:           vdc.dataItem[8:],
			}

			copy(events[itemCount:], []VehEvents{ve})
			itemCount = itemCount + 1
		case <-time.After(timeout):
			return nil, fmt.Errorf("VehReadEvents timeout after %d seconds", timeout)
		}
	}
	return events, nil
}
