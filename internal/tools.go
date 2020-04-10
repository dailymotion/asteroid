package internal

import (
	"bytes"
	"errors"
	"html"
	"io/ioutil"
	"net"
	"os/user"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// sort to show list with asc order
func SortedListPeer(listPeer []map[string]string) {
	realIPs := make([]net.IP, 0, len(listPeer))
	realKeys := make([]string, 0, len(listPeer))

	for _, v := range listPeer {
		for y, z := range v {
			realIPs =  append(realIPs, net.ParseIP(z))
			realKeys =  append(realKeys, y)
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

func RetrievePubKey(sshKey string) (string, error) {
	var keyName string

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	keyName = usr.HomeDir + "/.ssh/" + sshKey

	return keyName, nil
}

func ReadPubKey(sshKeyPath string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func CheckFlagValid(key string, address string, cmd string) error {
	if cmd == "add" {
		if key == "" {
			return errors.New("key is empty")
		} else if address == "" {
			return errors.New("address is empty")
		} else if strings.Contains(address, "10.0.0") && !strings.Contains(address, "/32") {
			return errors.New("You forgot the Netmask in -address\nShould be: 10.0.0.x/xx")
		}
	} else if cmd == "delete" {
		if key == "" {
			return errors.New("key is empty")
		}
	}
	return nil
}

func CreateEmoji() string {
	var emojiFinal string

	emoji := [][]int{
		// Emoticons decimal ID
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