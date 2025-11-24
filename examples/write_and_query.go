package main

import (
	"fmt"

	"github.com/0z0sk0/gohislip"
)

func main() {
	device := gohislip.NewHislipResource()
	err := device.Connect("TCPIP0::127.0.0.1::hislip0,4881::INSTR")
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}

	err = device.Write("SYST:PRES")
	if err != nil {
		fmt.Println("Write failed:", err)
		return
	}

	response, err := device.Query("SERV:PORT:COUN?")
	if err != nil {
		fmt.Println("Query failed:", err)
		return
	}

	fmt.Print("Response: " + response)
}
