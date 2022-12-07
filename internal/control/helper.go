package control

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// createUuid
//
//	Will generated a unique uuid for the node upon node creation
func createUuid() string {
	return uuid.New().String()
}

// determineIpAddress
//
//	Will get the ip address that makes this instance recognizable
//	Will determine if this node is inside of a k8s service and will get the service name if applicable
//	Will allow for the ip address to be configurable as well upon service initialization
func determineIpAddress() string {

	discoverLocalIp := func() string {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)

		fmt.Println(localAddr.IP.String())
		return localAddr.IP.String()
	}

	checkIfK8s := func() (bool, string) {
		readfile, err := os.Open("/etc/resolv.conf")
		if err != nil {
			return false, ""
		}
		searchRegexp := regexp.MustCompile(`^\s*search\s*(([^\s]+\s*)*)$`)
		fileScanner := bufio.NewScanner(readfile)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {

			match := searchRegexp.FindSubmatch([]byte(fileScanner.Text()))
			if match != nil && strings.Contains(string(match[1]), "svc.cluster.local") {
				return true, string(match[1])
			}
		}
		return false, ""
	}

	if isInK8s, svcValue := checkIfK8s(); isInK8s {
		fmt.Println(svcValue)
		// Let's do something with the service value here
		// Need to get the $HOSTNAME env variable here, then create the connection to it
	}

	return discoverLocalIp()
}
