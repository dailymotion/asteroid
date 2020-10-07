package tools

import (
	"errors"
	"fmt"
	"github.com/dailymotion/asteroid/pkg/config"
	"os"
)

func showConfig(wireguardConfig string) {
	fmt.Printf("\nClient config output\n" +
		"---------------\n" +
		"%v\n", wireguardConfig)
}

// WGConfig Struct with all related info to create peer on WG server
type WGConfig struct {
	PeerKey 		  string
	PeerCIDR 		  string
	PeerDeleteKey 	  string
	GenerateFile	  bool
	GenerateOutput	  bool
	conf 			  *config.Config
	ServerPubKey      string
}

// InitWG Init a Wireguard object to keep all related variables in one place
func InitWG(args []string ) (WGConfig, error) {
	conf, err := config.ReadConfigFile()
	if err != nil {
		return WGConfig{}, err
	}

	peerDeleteKey, peerKey, peerCIDR, generateFile, generateOutput := HandleFlags(args)
	serverPubKey := ""

	return WGConfig{
		PeerCIDR: *peerCIDR,
		PeerKey: *peerKey,
		PeerDeleteKey: *peerDeleteKey,
		GenerateFile:	  *generateFile,
		GenerateOutput:	  *generateOutput,
		conf: &conf,
		ServerPubKey: serverPubKey}, nil
}

// Generating client configuration with all the variables given
func (wg WGConfig) generateWGConfigFile() (string, error) {
	privateKey := "[Your Wireguard private key]"
	WGPort := wg.conf.WG.WGPort
	address := wg.PeerCIDR
	DNS := wg.conf.ClientConfig.DNS
	endpoint := wg.conf.WG.WireguardIP + WGPort
	allowedIPs := wg.conf.ClientConfig.AllowedIPs

	wireguardClientConfig := fmt.Sprintf(`[Interface]
PrivateKey = %v
Address = %v
DNS = %v

[Peer]
PublicKey = %v
AllowedIPs = %v
EndPoint = %v
`, privateKey, address, DNS, wg.ServerPubKey, allowedIPs, endpoint)

	return wireguardClientConfig, nil
}

// Writing the client configuration in a file whose name is given in the asteroid.yaml file
func (wg WGConfig) writeWGConfToFile(wireguardConf string) error {
	// Check if the file exist
	_, err := os.Stat(wg.conf.ClientConfig.Name)
	if err != nil {
		// If file doesn't exist we create it
		f, err := os.Create(wg.conf.ClientConfig.Name)
		if err != nil {
			return err
		}
		defer f.Close()

		// Writing the configuration to the file
		_, err = f.WriteString(wireguardConf)
		if err != nil {
			return err
		}
	} else {
		// If file exist we open it
		f, err := os.OpenFile(wg.conf.ClientConfig.Name, os.O_WRONLY, 0644)
		if err != nil {
			str := fmt.Sprintf("File: %v in OPEN", wg.conf.ClientConfig.Name)
			return errors.New(str)
		}
		// Writing the config to the file
		_, err = fmt.Fprintln(f, wireguardConf)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return err
		}
	}

	fmt.Printf(
		"Config file created in this folder with the name: %v\n" +
			"-----------------------------------------------------------\n", wg.conf.ClientConfig.Name)

	return nil
}

// RetrieveWGConfig Will generate, show or write to a file the client configuration
func (wg WGConfig) RetrieveWGConfig() error {

	wireguardConf, err := wg.generateWGConfigFile()
	if err != nil {
		return err
	}

	if wg.GenerateFile {
		err := wg.writeWGConfToFile(wireguardConf)
		if err != nil {
			return err
		}
	}

	if wg.GenerateOutput {
		showConfig(wireguardConf)
	}
	return nil
}
