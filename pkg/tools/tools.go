package tools

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"os/user"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// SortedListPeer Sorts to show list with asc order
func SortedListPeer(listPeer []map[string]string) {
	realIPs := make([]net.IP, 0, len(listPeer))
	realKeys := make([]string, 0, len(listPeer))

	for _, v := range listPeer {
		for y, z := range v {
			realIPs = append(realIPs, net.ParseIP(z))
			realKeys = append(realKeys, y)
		}
	}

	sort.Slice(realIPs, func(i, j int) bool {
		return bytes.Compare(realIPs[i], realIPs[j]) < 0
	})

	for k, v := range listPeer {
		for y := range v {
			listPeer[k][y] = realIPs[k].String()
		}
	}

}

// RetrievePubKey Going to the .ssh folder in user home to retrieve the ssh key to connect to Wireguard server
func RetrievePubKey(sshKey string) (string, error) {
	var keyName string

	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve the current user")
	}
	keyName = usr.HomeDir + "/.ssh/" + sshKey

	return keyName, nil
}

// ReadPubKey Reads ssh key previously obtained
func ReadPubKey(sshKeyPath string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the public key file")
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the private key")
	}
	return ssh.PublicKeys(key), nil
}

// CheckFlagValid Checks that nothing is missing or that a flag as all the requirements
func CheckFlagValid(wireguard WGConfig, cmd string) error {
	var key string
	var address string

	if cmd == "add"{
		key = wireguard.PeerKey
		address = wireguard.PeerCIDR
	} else {
		key = wireguard.PeerDeleteKey
		address = wireguard.PeerCIDR
	}

	fmt.Printf("Key: %v\naddresse: %v\n", key, address)

	switch cmd {
	case "add":
		if key == "" {
			return errors.New("key is empty")
		}
		if address == "" {
			return errors.New("address is empty")
		}
		if strings.Contains(address, "10.0.0") && !strings.Contains(address, "/32") {
			return errors.New("You forgot the Netmask in -address\nShould be: 10.0.0.x/xx")
		}
	case "delete":
		if key == "" {
			return errors.New("key is empty")
		}
	default:
		return errors.Errorf("unexpected command: %s", cmd)
	}

	return nil
}

// CreateEmoji Creates the planet emoji present in the Asteroid help info
func CreateEmoji() string {
	var emojiFinal string

	emoji := [][]int{
		// Emoticons decimal ID of a planet
		{127759, 127760},
	}

	for _, value := range emoji {
		for x := value[0]; x < value[1]; x++ {
			// Unescape the string (HTML Entity -> String).
			str := html.UnescapeString("&#" + strconv.Itoa(x) + ";")
			emojiFinal = str
		}
	}
	return emojiFinal
}

