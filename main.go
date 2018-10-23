package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/RobHumphris/Ble-test/ble"
)

func main() {
	a, err := ble.NewVehDeviceConnection("ED:24:62:72:24:9B")

	err = a.VehUnlock()
	if err != nil {
		log.Fatal(err)
	}

	version, err := a.VehGetVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Version\t", version)

	info, err := a.VehGetInfo()
	if err != nil {
		log.Fatal(err)
	}
	inf, _ := json.Marshal(info)
	fmt.Println("Info\t", string(inf))

	config, err := a.VehGetConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg, _ := json.Marshal(config)
	fmt.Println("Config\t", string(cfg))

	friendlyName, err := a.VehGetFriendlyName()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name\t", friendlyName)

	eventsLog, err := a.VehReadEvents(info.CurrentSequence)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Event Log size\t", len(eventsLog))
	err = a.Finalize()
	if err != nil {
		log.Fatal(err)
	}
}
