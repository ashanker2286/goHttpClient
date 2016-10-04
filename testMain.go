package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func configurePorts(ipAddr, port string) error {
	for _, portInfo := range testCfg.Ports {
		jsonStr, err := json.Marshal(portInfo)
		if err != nil {
			return err
		}
		url := "http://" + ipAddr + ":" + port + "/public/v1/config/Port"
		fmt.Println("URL:>", url)

		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		SendHttpCmd(req)
	}
	return nil
}

func configureVlan(ipAddr, port string) error {
	for _, vlanInfo := range testCfg.Vlans {
		jsonStr, err := json.Marshal(vlanInfo)
		if err != nil {
			return err
		}
		url := "http://" + ipAddr + ":" + port + "/public/v1/config/Vlan"
		fmt.Println("URL:>", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		SendHttpCmd(req)
	}
	return nil
}

func configureIPv4Intf(ipAddr, port string) error {
	for _, ipv4Info := range testCfg.IPv4Intf {
		jsonStr, err := json.Marshal(ipv4Info)
		if err != nil {
			return err
		}
		url := "http://" + ipAddr + ":" + port + "/public/v1/config/IPv4Intf"
		fmt.Println("URL:>", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		SendHttpCmd(req)
	}
	return nil
}

func testSetup(ipAddr string, port string) error {
	fmt.Println("Run Setup")
	err := configurePorts(ipAddr, port)
	if err != nil {
		return err
	}

	err = configureVlan(ipAddr, port)
	if err != nil {
		return err
	}

	err = configureIPv4Intf(ipAddr, port)
	if err != nil {
		return err
	}
	return nil
}

func addRoute(ipAddr, port string) error {
	fmt.Println("Add Route")
	for _, routeInfo := range testCfg.Routes {
		jsonStr, err := json.Marshal(routeInfo)
		if err != nil {
			return err
		}
		url := "http://" + ipAddr + ":" + port + "/public/v1/config/IPv4Route"
		fmt.Println("URL:>", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		SendHttpCmd(req)
	}
	return nil
}

func delRoute(ipAddr, port string) error {
	fmt.Println("Delete Route")
	for _, routeInfo := range testCfg.Routes {
		jsonStr, err := json.Marshal(routeInfo)
		if err != nil {
			return err
		}
		url := "http://" + ipAddr + ":" + port + "/public/v1/config/DeleteIPv4Route"
		fmt.Println("URL:>", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		SendHttpCmd(req)
	}
	return nil
}

func pollData(ipAddr, port string) error {
	fmt.Println("Poll Data")
	fmt.Println("Port Data -")
	fmt.Println("============")
	var jsonStr = []byte(nil)
	url := "http://" + ipAddr + ":" + port + "/public/v1/state/Ports"
	fmt.Println("URL:>", url)

	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	SendHttpCmd(req)

	fmt.Println("Routes -")
	fmt.Println("========")
	jsonStr = []byte(nil)
	url = "http://" + ipAddr + ":" + port + "/public/v1/state/IPv4Routes"
	fmt.Println("URL:>", url)

	req, _ = http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Accept", "application/json")
	SendHttpCmd(req)
	return nil
}

func main() {
	var err error
	parseTestCfg()
	ipAddr := flag.String("IP", "localhost", "Ip Address")
	port := flag.String("Port", "8080", "port number")
	fmt.Println(*ipAddr)
	fmt.Println(*port)
	var execMode string
	flag.StringVar(&execMode, "mode", "configure",
		`[setup|addroute|delroute|polldata]
	           a) setup :
	           - Apply port config
	           - Create Vlan
	           - Create IPv4Intf

	           b) addroute :
	           - Create static route

	           c) delroute :
	           - Delete static route

	           d) polldata :
	           - Read and display port stats
	           - Read and display routes`)

	flag.Parse()
	switch execMode {
	case "setup":
		err = testSetup(*ipAddr, *port)
		if err != nil {
			fmt.Println("Failed to apply test configuration - ", err.Error())
		}
	case "addroute":
		err = addRoute(*ipAddr, *port)
		if err != nil {
			fmt.Println("Failed to add route - ", err.Error)
		}
	case "delroute":
		err = delRoute(*ipAddr, *port)
		if err != nil {
			fmt.Println("Failed to delete route - ", err.Error())
		}
	case "polldata":
		err = pollData(*ipAddr, *port)
		if err != nil {
			fmt.Println("Failed to poll port stats/route data - ", err.Error())
		}
	default:
		fmt.Println("Invalid exec mode. Aborting !")
	}
	return
}

func SendHttpCmd(req *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
