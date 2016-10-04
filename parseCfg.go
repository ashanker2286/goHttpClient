package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	SETUP_CFG_FILE = "testCfg.json"
)

//Struct definition for port obj in setup json file
type PortData struct {
	//Name         string
	//Speed        int
	//AdminState   string
	//Mtu          int
	//BreakOutMode string
	IntfRef    string
	AdminState string
	Speed      int
	Autoneg    string
}

//Struct definition for vlan obj in setup json file
type VlanData struct {
	VlanId        int
	IntfList      string
	UntagIntfList string
}

//Struct definition for ipv4intf obj in setup json file
type IPv4IntfData struct {
	IpAddr  string
	IntfRef string
}

type NextHopInfo struct {
	NextHopIp string
}

//Struct definition for route obj in setup json file
type RouteData struct {
	DestinationNw string
	NetworkMask   string
	Cost          int
	NextHop       []NextHopInfo
	Protocol      string
}

//Test configuration data
type testCfgData struct {
	Ports    []PortData
	Vlans    []VlanData
	IPv4Intf []IPv4IntfData
	Routes   []RouteData
}

//Global test setup data
var testCfg testCfgData

func parseTestCfg() error {
	//Read test data json file
	bytes, err := ioutil.ReadFile(SETUP_CFG_FILE)
	if err != nil {
		fmt.Println("Failed to read ", SETUP_CFG_FILE)
		return err
	}
	//Parse json byte stream
	err = json.Unmarshal(bytes, &testCfg)
	if err != nil {
		fmt.Println("Failed to parse json byte stream - ", err)
		return err
	}
	return nil
}
