package peer

import (
	"github.com/dailymotion/asteroid/network"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func AddNewPeer(conn *ssh.Client, peerKey string, clientCIDR string) error {
	command := "sudo wg set wg0 peer " + peerKey + " allowed-ips " + clientCIDR

	_, err := network.RunCommand(conn, command)
	return errors.Wrapf(err, "failed to run the command: %s", command)
}

func DeletePeer(conn *ssh.Client, peerKey string) error {
	command := "sudo wg set wg0 peer " + peerKey + " remove"

	_, err := network.RunCommand(conn, command)
	return errors.Wrapf(err, "failed to run the command: %s", command)
}
