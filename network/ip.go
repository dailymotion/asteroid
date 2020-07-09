package network

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"

	//"github.com/dailymotion/asteroid/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

var (
	// regex for the ip
	regexFindIP = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	// regex for the address
	findPeerKey = regexp.MustCompile(`(peer:|public\ key:)\s(\W.+|\w.+)`)
	// 10.0.0.0
	cidrTwentyfourBit = "10.0"
	// 172.16.0.0
	cidrTwentyBit = "172.16"
	// 192.168.0.0
	cidrSixteenBit = "192.168"
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

func readerToString(cmdReader io.Reader) (string, error) {
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
	stdOut, err := RunCommand(conn, command)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to run the command: %s", command)
	}

	// regex to retrieve only the lines with the IP and the peer address
	for _, line := range strings.Split(strings.TrimSuffix(stdOut, "\n"), "\n") {
		peerIPs := make(map[string]string)
		ipAddress := regexFindIP.FindStringSubmatch(line)
		peerKey := findPeerKey.FindStringSubmatch(line)

		if len(peerKey) > 0 || len(peerKey) > 0 {
			if len(peerKey) > 0 {
				key = peerKey[2]
			}
			if len(ipAddress) > 0 {
				for _, v := range ipAddress {
					if strings.Contains(v, "10.0") || strings.Contains(v, "172.16") {
						peerIPs[key] = v
						key = ""
					}
				}
			}
		}
		if len(peerKey) > 0 {
			if peerKey[1] == "public key:" {
				//TODO to mention that this key belongs to the server itself
			} else {
				key = peerKey[2]
			}
		}

		if len(ipAddress) > 0 {
			for _, v := range ipAddress {
				if strings.Contains(v, "10.0") || strings.Contains(v, "172.16") || strings.Contains(v, "192.168"){
					peerIPs[key] = v
					key = ""
				}
			}
		}

		if len(peerIPs) > 0 {
			listPeers = append(listPeers, peerIPs)
		}
	}

	// Sort function issue to fix: key/ip mismatch
	//internal.SortedListPeer(listPeers)
	return listPeers, nil
}