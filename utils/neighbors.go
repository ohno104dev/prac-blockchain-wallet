package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

func IsFoundNode(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	if _, err := net.DialTimeout("tcp", target, time.Second*1); err != nil {
		fmt.Printf("%s %v\n", target, err)
		return false
	}

	return true
}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(hostIp string, port uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", hostIp, port)

	m := PATTERN.FindStringSubmatch(hostIp)
	if m == nil {
		return nil
	}
	ipPrefix := m[1]
	hostIdent, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for guessPort := startPort; guessPort <= endPort; guessPort += 1 {
		for varHostIdent := startIp; varHostIdent <= endIp; varHostIdent += 1 {
			guessIp := fmt.Sprintf("%s%d", ipPrefix, hostIdent+int(varHostIdent))
			guessTarget := fmt.Sprintf("%s:%d", guessIp, guessPort)
			if guessTarget != address && IsFoundNode(guessIp, guessPort) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}

	return neighbors
}

func GetHost() string {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Println("ERROR:", err)
		os.Exit(1)
	}
	defer conn.Close()

	address := conn.LocalAddr().(*net.UDPAddr)
	ipStr := fmt.Sprintf("%v", address.IP)

	return ipStr
}
