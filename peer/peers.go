package peer

import (
	"fmt"
	"log"

	"github.com/dailymotion/asteroid/network"
	"golang.org/x/crypto/ssh"
)

func AddNewPeer(conn *ssh.Client, peerKey string, clientCIDR string) (string, string) {
	command := "sudo wg set wg0 peer " + peerKey + " allowed-ips " + clientCIDR

	stdOut, stdErr, err := network.RunCommand(conn, command)
	if err != nil {
		log.Fatalf("\nError with runCommand: %v\nStdErr: %v", err, stdErr)
	}

	stringOut, err := network.ReaderToString(stdOut)
	if err != nil {
		fmt.Println("Issue with reader to string: ", err)
	}

	stringErr, err := network.ReaderToString(stdErr)
	if err != nil {
		fmt.Println("Issue with reader to string: ", err)
	}

	return stringOut, stringErr
}

func DeletePeer(conn *ssh.Client, peerKey string) (string, string) {
	command := "sudo wg set wg0 peer " + peerKey + " remove"

	stdOut, stdErr, err := network.RunCommand(conn, command)
	if err != nil {
		log.Fatalf("\nError with runCommand: %v\nStdErr: %v", err, stdErr)
	}
	stringOut, err := network.ReaderToString(stdOut)
	if err != nil {
		fmt.Println("Issue with reader to string: ", err)
	}

	stringErr, err := network.ReaderToString(stdErr)
	if err != nil {
		fmt.Println("Issue with reader to string: ", err)
	}

	return stringOut, stringErr
}