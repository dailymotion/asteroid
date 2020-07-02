package network

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dailymotion/asteroid/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func ShowListIPs(listPeers []map[string]string) {
	var data [][]string
	var row []string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Peer", "Local IP"})

	for _, v := range listPeers {
		for x, y := range v {
			row := []string{x, y}
			data = append(data, row)
		}
		data = append(data, row)
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func CheckForDouble(listPeer []map[string]string, IPAddress string) bool {
	cleanIP := IPAddress[:len(IPAddress)-3]

	for _, v := range listPeer {
		for _, ip := range v {
			if ip == cleanIP {
				return true
			}
		}
	}
	return false
}

func ReaderToString(cmdReader io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(cmdReader)
	if err != nil {
		return "", err
	}
	s := buf.String()

	return s, nil
}

func RetrieveIPs(conn *ssh.Client) ([]map[string]string, error) {
	var listPeers []map[string]string
	key := ""

	// command to show all peers created on the server
	command := "sudo wg"
	stdOut, stdErr, err := RunCommand(conn, command)
	if err != nil {
		return listPeers, errors.Wrapf(err, "Error with runCommand.\nStdErr: %v", stdErr)
		//log.Fatalf("\nError with runCommand: %v\nStdErr: %v", err, stdErr)
	}

	s, err := ReaderToString(stdOut)
	if err != nil {
		return nil, err
	}

	// regex to retrieve only the lines with the IP and the peer address
	for _, line := range strings.Split(strings.TrimSuffix(s, "\n"), "\n") {
		peerIPs := make(map[string]string)

		// regex for the ip
		regexFindIP, err := regexp.Compile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
		if err != nil {
			return listPeers, errors.Wrapf(err, "Error compiling Regex to find IPs\n")
			//log.Fatalf("\nError with compiling Regex: %v\n", err)
		}
		// regex for the address
		//findPeerKey, err := regexp.Compile(`peer:\s(\W.+|\w.+)`)
		findPeerKey, err := regexp.Compile(`(peer:|public\ key:)\s(\W.+|\w.+)`)
		if err != nil {
			return listPeers, errors.Wrapf(err, "Error compiling Regex to find peer key\n")
			//log.Fatalf("\nError with compiling Regex: %v\n", err)
		}

		regex := regexFindIP.FindStringSubmatch(line)
		regex2 := findPeerKey.FindStringSubmatch(line)

		if len(regex2) > 0 || len(regex) > 0 {
			if len(regex2) > 0 {
				fmt.Println(regex2[2])
				key = regex2[2]
			}
			if len(regex) > 0 {
				for _, v := range regex {
					if strings.Contains(v, "10.0") || strings.Contains(v, "172.16"){
						peerIPs[key] = v
						key = ""
					}
				}
			}
		}

		if len(peerIPs) > 0 {
			listPeers = append(listPeers, peerIPs)
		}
	}
	internal.SortedListPeer(listPeers)
	return listPeers, nil
}