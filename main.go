package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/micmonay/keybd_event"
)

const waitforseconds time.Duration = 20

func main() {

	fmt.Printf("AgentStarter - custom written for Wiltshires garage, Liphook\nCopyright (c) 2021 Bob Stammers\n\n")

	waitforagent()

	cmd := exec.Command("cmd", "/C", "start msedge --kiosk http://localhost:8090 --edge-kiosk-type=fullscreen")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(waitforseconds * time.Second)
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_ESC)

	// Press the selected keys
	err = kb.Launching()
	if err != nil {
		panic(err)
	}

	fmt.Printf("My work here is done\n")
	time.Sleep(waitforseconds * time.Second)

}

func waitforagent() {

	for {
		ok := raw_connect("127.0.0.1", "8090")
		if ok {
			return
		}
		fmt.Printf("Waiting for Agent startup ... ")
		time.Sleep(waitforseconds * time.Second)

	}
}

func raw_connect(host string, port string) bool {

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		fmt.Println("patience little one, patience")
		return false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("ok")
		return true
	}
	fmt.Println()
	return false
}
