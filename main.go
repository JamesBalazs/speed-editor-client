package main

import (
	"fmt"
	"log"

	"github.com/JamesBalazs/speed-editor-rebind/auth"
	"github.com/sstallion/go-hid"
)

const VID = 0x1edb
const PID = 0xda0e

func main() {
	// Initialize the hid package.
	if err := hid.Init(); err != nil {
		log.Fatal(err)
	}

	// Open the device using the VID and PID.
	d, err := hid.OpenFirst(VID, PID)
	if err != nil {
		log.Fatal(err)
	}
	defer hid.Exit()

	// Read the Manufacturer String.
	s, err := d.GetMfrStr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Manufacturer String: %s\n", s)

	// Read the Product String.
	s, err = d.GetProductStr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Product String: %s\n", s)

	// Read the Serial Number String.
	s, err = d.GetSerialNbr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Serial Number String: %s\n", s)

	auth.Authenticate(d)

	for {
		data := make([]byte, 16)
		i, err := d.Read(data)

		fmt.Printf("len: %d data: %v\n", i, data)

		if err != nil {
			log.Fatal(err)
		}
	}
}
