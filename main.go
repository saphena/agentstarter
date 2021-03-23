package main

/*
 * AgentStart
 *
 * This acts as a startup manager for some software called 'Agent' which processes and presents
 * live video streams from webcams and similar devices.
 *
 * I expect to operate on a headless computer without human intervention so I need to bring up
 * a web browser in kiosk mode and deal with anything which would otherwise need human intervention.
 *
 * v1.0		23Mar21	Live release to Wiltshires
 *
 */
import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"time"

	"github.com/micmonay/keybd_event"
)

const waitforseconds time.Duration = 20 // A suitable interval, not too long, not too short

func main() {

	fmt.Printf("AgentStarter - custom written for Wiltshires garage, Liphook\nCopyright (c) 2021 Bob Stammers\n\n")

	waitforagent() // Wait until the Agent server up and running

	cmd := exec.Command("cmd", "/C", "start msedge --kiosk http://localhost:8090 --edge-kiosk-type=fullscreen")

	err := cmd.Start() // Fire up a web browser to present live video
	if err != nil {
		panic(err)
	}

	time.Sleep(waitforseconds * time.Second) // Give it some time to do its thing

	// Now need to close the "language" dialog box so ...

	kb, err := keybd_event.NewKeyBonding() // Interface to the virtual keyboard
	if err != nil {
		panic(err)
	}

	kb.SetKeys(keybd_event.VK_ESC) // Prepare to hit [Esc]

	err = kb.Launching() // Hit it
	if err != nil {
		panic(err)
	}

	fmt.Printf("My work here is done\n")

	time.Sleep(waitforseconds * time.Second) // Stay alive long enough for that [Esc] to do its thing

}

func waitforagent() {

	patience := []string{"patience young grasshopper",
		"patience, my padawan",
		"patience little one, patience",
		"calm yourself",
		"remain calm, all will be well",
		"haven't you got something to be getting on with?",
		"any minute now",
		"it'll take longer if you watch",
		"all things come to those who wait"}

	rand.Seed(time.Now().UnixNano()) // A good enough 'random' number

	fmt.Printf("Are we in a 'go' situation already?  ")
	for {
		ok := raw_connect("127.0.0.1", "8090") // Try connecting to the Agent server
		if ok {
			fmt.Println("ok")
			return
		}
		fmt.Println(patience[rand.Intn(len(patience))])
		fmt.Printf("Waiting for Agent startup ... ")
		time.Sleep(waitforseconds * time.Second)
	}
}

func raw_connect(host string, port string) bool {

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	fmt.Println()
	return false
}
