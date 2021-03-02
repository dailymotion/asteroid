package tools

import (
	"database/sql"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/dailymotion/asteroid/pkg/config"
	"github.com/dailymotion/asteroid/pkg/db"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func CheckArguments(args []string, flagType string) error {
	switch flags := flagType; flags {
	case "add":
		if len(args) <= 5 {
			return errors.New("missing Arguments\n")
		}
	case "view":
		// We alert if too much arguments are given to the command
		if len(os.Args) > 2 {
			flag.Usage()
			return errors.New("View doesn't take options\n\n")
		}
	case "delete":
		if len(os.Args) < 3 {
			return errors.New("issue with argument")
		}
	}
	return nil
}

func CheckIfPresent(peerList []map[string]string, deleteKey string) bool {

	for _, peers := range peerList {
		for key, value := range peers {
			if key == "key" {
				if value == deleteKey  {
					return true
				}
			}
		}
	}
	return false
}

func PrintResult(flagType string, key string) {
	switch flagType {
	case "add":
		// Message be like:
		//################
		//# Peer added ! #
		//################
		fmt.Printf("\n################\n# Peer added ! #\n################\n")
		fmt.Printf("Peer: %v has been added !\n", key)
	case "delete":
		// Message be like:
		//##################
		//# Peer deleted ! #
		//##################
		fmt.Printf("\n##################\n# Peer deleted ! #\n##################\n")
		fmt.Printf("Peer: %v has been deleted !\n", key)
	}
}

// AddToListPeer Will add the corresponding user to it's key in the list
func AddToListPeer(listPeers []map[string]string, DBConn *sql.DB) ([]map[string]string, error){
	var peersList []map[string]string
	for _, value := range listPeers {
		for key, cidr := range value {
			tmpList := make(map[string]string)
			tmpList["cidr"] = cidr
			tmpList["key"] = key
			tmpList["username"] = ""
			if DBConn != nil {
				DBUserList, err := db.ReadKeyFromDB(DBConn)
				if err != nil {
					return peersList, err
				}
				for _, v := range DBUserList {
					if tmpList["key"] == v.Key {
						tmpList["username"] = v.Username
					}
				}
			}
			peersList = append(peersList, tmpList)
		}
	}
	return peersList, nil
}

// InitDB will retrieve the config and start a connection with the DB
func InitDB() (*sql.DB,  config.ConfigDB, error) {
	var conn *sql.DB

	conf, err := config.RetrieveDBConfig()
	if err != nil {
		return nil, conf, err
	}

	if conf.DBEnabled {
		conn, err = db.ConnectToDB(conf)
		if err != nil {
			return nil, conf, err
		}
	}

	return conn, conf, nil
}

// RetrievePubKey Will go to the .ssh folder in the user $HOME, retrieve the ssh key and connect to the Wireguard server
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

	if cmd == "add" {
		key = wireguard.PeerKey
		address = wireguard.PeerCIDR
	} else {
		key = wireguard.PeerDeleteKey
		address = wireguard.PeerCIDR
	}

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

func CheckForDouble(peerList []map[string]string, wireguard *WGConfig) bool {
	for _, peers := range peerList {
		for key, cidr := range peers {
			if key == wireguard.PeerKey || cidr == wireguard.PeerCIDR {
				return true
			}
		}
	}
	return false
}

func RetrieveAndCheckForDouble(DBConn *sql.DB, DBConf config.ConfigDB,
	wireguard *WGConfig, tmpListPeers []map[string]string, serverPubKey string) error {

	wireguard.ServerPubKey = serverPubKey

	// If DB is enabled, add user into DB
	for _, peers := range tmpListPeers {
		for key, cidr := range peers {
			if key == wireguard.PeerKey || cidr == wireguard.PeerCIDR {
				return errors.New("IP/Key already exist on the server")
			}
		}
	}
	if DBConf.DBEnabled {
		err := db.InsertUserInDB(DBConn, wireguard.PeerKey, wireguard.PeerCIDR, wireguard.PeerName)
		if err != nil {
			return err
		}
	}
	return nil
}
