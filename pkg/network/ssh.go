package network

import (
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dailymotion/asteroid/pkg/config"
	"github.com/dailymotion/asteroid/pkg/tools"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// RunCommand will use the ssh connection to run the bash command on the server
func RunCommand(conn *ssh.Client, cmd string) (string, error) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "failed to connect stdout pipi")
	}

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return "", errors.Wrap(err, "failed to connect stderr pipe")
	}

	if err = sess.Run(cmd); err != nil {
		return "", errors.Wrapf(err, "failed to run the command: %s", cmd)
	}

	stringOut, err := readerToString(sessStdOut)
	if err != nil {
		return "", errors.Wrap(err, "failed to read the stdout")
	}

	stringErr, err := readerToString(sessStderr)
	if err != nil {
		return "", errors.Wrap(err, "failed to read the stderr")
	}

	if len(stringErr) > 0 {
		return "", errors.Errorf("The command returned with the following error: %s", stringErr)
	}

	return stringOut, nil
}

// ConnectAndRetrieve connect to the WG server and retrieve all the peers informations
func ConnectAndRetrieve(wireguard tools.WGConfig, cmd string) (*ssh.Client, error) {
	//We read the config file to retrieve the connection arguments
	configWG, err := config.ReadConfigFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the config file")
	}

	// We retrieve the ssh key path
	sshKeyPath, err := tools.RetrievePubKey(configWG.WG.SSHKeyName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve the public key")
	}

	key, err := tools.ReadPubKey(sshKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read the public key")
	}

	// We're creating the connection to the server
	conn, err := connectToServer(key, configWG)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the server. Are you connected to the VPN ?")
	}

	if cmd == "add" {
		// We retrieve all the user vpn ip to use them for different checks
		listPeers, _, err := RetrieveIPs(conn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve VPN IPs")
		}

		if ok := CheckForDouble(listPeers, wireguard.PeerCIDR); ok {
			return nil, errors.New("IP already exists on the server")
		}
	}

	return conn, nil
}

func connectToServer(sshKey ssh.AuthMethod, config config.Config) (*ssh.Client, error) {
	// Build our new spinner (connection animation)
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	log.Printf("Connecting to server...")
	s.Start()

	// SSH config
	sshConfig := ssh.ClientConfig{}
	sshConfig.User = config.WG.Username
	sshConfig.Timeout = 10 * time.Second
	sshConfig.Auth = []ssh.AuthMethod{
		sshKey,
	}

	// If user doesn't want HostKey verification
	if !config.WG.HostKey {
		sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	connection, err := ssh.Dial("tcp", config.WG.WireguardIP+":"+config.WG.SSHPort, &sshConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ssh dial %s", config.WG.WireguardIP+":"+config.WG.SSHPort)
	}

	s.Stop()
	return connection, nil
}
