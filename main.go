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
 * v1.1		22Apr21 Help function + email alerts
 *
 */
import (
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/mail"
	"net/smtp"
	"os/exec"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
)

var helpWanted *bool = flag.Bool("?", false, "Show help")

const waitforseconds time.Duration = 20 // A suitable interval, not too long, not too short

type Camera struct {
	Name  string `xml:"name,attr"`
	URL   string `xml:"settings>substream"`
	Login string `xml:"settings>login"`
}

type Cameras struct {
	XMLName xml.Name `xml:"objects"`
	C       []Camera `xml:"cameras>camera"`
}

func main() {

	fmt.Printf("AgentStarter - custom written for Wiltshires garage, Liphook\nCopyright (c) 2021 Bob Stammers\n\n")
	flag.Parse()

	if *helpWanted {
		showhelp()
		return
	}

	sendmail("Starting")

	waitforagent() // Wait until the Agent server up and running

	sendmail("Launching browser")

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

func sendmail(msg string) { // msg is used for subject and body so keep it short

	from := mail.Address{Name: "Wiltshires/Bobby", Address: "robot@pikit.co.uk"}
	to := mail.Address{Name: "", Address: "stammers.bob@gmail.com"}
	subj := msg
	body := msg

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700")

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.livemail.co.uk:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "robot@pikit.co.uk", "vs4hHpJS6", host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	c.Quit()

}
func showhelp() {
	const agentXML = `C:\Program Files\Agent\Media\XML\objects.xml`
	const txt = `The cameras, TP_Link Tapo C100s, are controlled by Agent DVR service (https://www.ispyconnect.com/download.aspx) which is configured using its browser menu.

Each camera needs a 'Camera Account' configured using the Tapo phone app, userid is probably 'saphena' with the usual password but see below.

Agent actually stores its configuration in %v

This program, AgentStarter, needs no configuration. It waits until it can see that Agent is serving then fires
up Microsoft Edge running in kiosk mode to present the output feeds.

`
	fmt.Printf(txt, agentXML)

	b, err := ioutil.ReadFile(agentXML)
	if err != nil {
		return
	}

	// The file claims to be utf16 but is actually utf8
	b = []byte(strings.Replace(string(b), ` encoding="utf-16"`, ``, 1))
	//fmt.Printf("Unmarshalling ...\n")
	c := Cameras{}
	err = xml.Unmarshal(b, &c)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	for _, cc := range c.C {
		fmt.Printf("Camera '%v' = %v [%v]\n", cc.Name, cc.URL, cc.Login)
	}
	fmt.Println()
	sendmail("Testing")
}

func waitforagent() {

	patience := []string{"patience young grasshopper",
		"patience, my padawan",
		"patience little one, patience",
		"calm yourself",
		"remain calm, all will be well",
		"haven't you got something to be getting on with?",
		"any minute now",
		"hang on a bit",
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
