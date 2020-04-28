package network

import (
	//"errors"
	"fmt"
	"io"
	//"log"
	//"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dailymotion/asteroid/config"
	"github.com/dailymotion/asteroid/internal"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func RunCommand(conn *ssh.Client, cmd string) (io.Reader, io.Reader, error){
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	err = sess.Run(cmd)
	if err != nil {
		return nil, nil, err
	}
	return sessStdOut, sessStderr, nil
}

func ConnectAndRetrieve(IPAddress string, cmd string) (*ssh.Client, error) {
	//We read the config file to retrieve the connection arguments
	configWG, err := config.ReadConfigFile()
	if err != nil {
		return nil, err
	}

	// We retrieve the ssh key path
	sshKeyPath, err := internal.RetrievePubKey(configWG.SSHKeyName)
	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving PubKey")
	}

	key, err := internal.ReadPubKey(sshKeyPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading PubKey")
	}

	// We're creating the connection to the server
	conn, err := connectToServer(key, configWG)
	if err != nil {
		return nil, errors.Wrapf(err, "Are you connected to the VPN ? ")
	}

	if cmd == "add" {
		// We retrieve all the user vpn ip to use them for different checks
		listPeers, err := RetrieveIPs(conn)
		if err != nil {
			return nil, errors.Wrapf(err, "error retrieving VPN IPs")
		}

		ok := CheckForDouble(listPeers, IPAddress)
		if ok {
			return nil, errors.Wrapf(err,"IP already exist in the server")
		}
	}

	return conn, nil
}

func connectToServer(sshKey ssh.AuthMethod, config config.Config) (*ssh.Client, error) {
	// Build our new spinner (connection animation)
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	fmt.Printf("Connecting to server ")
	s.Start()

	// SSH config
	sshConfig := ssh.ClientConfig{}
	sshConfig.User = config.Username
	sshConfig.Timeout = 10 * time.Second
	sshConfig.Auth = []ssh.AuthMethod{
		sshKey,
	}

	// If user doesn't want HostKey verification
	if !config.HostKey {
		sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	connection, err := ssh.Dial("tcp", config.WireguardIP +  ":" + config.SSHPort, &sshConfig)
	if err != nil {
		return nil, err
	}
	s.Stop()
	return connection, nil
}