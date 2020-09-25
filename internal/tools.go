package internal

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"

	"github.com/dailymotion/asteroid/config"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// sort to show list with asc order
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

func RetrievePubKey(sshKey string) (string, error) {
	var keyName string

	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve the current user")
	}
	keyName = usr.HomeDir + "/.ssh/" + sshKey

	return keyName, nil
}

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

func CheckFlagValid(key string, address string, cmd string) error {
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

func showConfig(wireguardConfig string) {
	fmt.Println(wireguardConfig)
}

func generateWGConfigFile(peerKey *string, peerCIDR *string, conf *config.Config) (string, error) {
	privateKey := peerKey
	address := peerCIDR
	DNS := "9.9.9.9"
	endpoint := conf.WireguardIP
	allowedIPs := "0.0.0.0/0"

	wireguardClientConfig := fmt.Sprintf(`[Interface]
PrivateKey = %v
Address = %v
DNS = %v

[Peer]
PubblicKey = TODO
AllowedIPs = %v
EndPoint = %v
`, *privateKey, *address, DNS, allowedIPs, endpoint)

	return wireguardClientConfig, nil
}

func writeWGConfToFile(wireguardConf string, conf *config.Config) error {
	confToByte := []byte(wireguardConf)

	err := ioutil.WriteFile(conf.ClientConfigFile, confToByte, 0644)
	if err != nil {
		return err
	}

	f, err := os.Create(conf.ClientConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(wireguardConf)
	if err != nil {
		return err
	}

	log.Printf(
		"The wireguard config for the client has been created in this folder with the name: %v\n",
		conf.ClientConfigFile)

	return nil
}

func RetrieveWGConfig(generateFileFlag bool, generateOutputFlag bool, peerKeyFlag string, peerCIDRFlag string) error {
	conf, err  := config.ReadConfigFile()
	if err != nil {
		return err
	}

	wireguardConf, err := generateWGConfigFile(&peerKeyFlag, &peerCIDRFlag, &conf)
	if err != nil {
		return err
	}

	if generateFileFlag {
		err := writeWGConfToFile(wireguardConf, &conf)
		if err != nil {
			return err
		}
	}

	if generateOutputFlag {
		showConfig(wireguardConf)
	}
	return nil
}
