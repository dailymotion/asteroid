package peer

import (
	"github.com/dailymotion/asteroid/pkg/network"
	"github.com/dailymotion/asteroid/pkg/tools"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// AddNewPeer With ssh connection and peer information, Adds the peer to the Wireguard server
func AddNewPeer(conn *ssh.Client, wireguard tools.WGConfig) error {
	peerKey := wireguard.PeerKey
	clientCIDR := wireguard.PeerCIDR

	command := "sudo wg set wg0 peer " + peerKey + " allowed-ips " + clientCIDR

	_, err := network.RunCommand(conn, command)
	return errors.Wrapf(err, "failed to run the command: %s", command)
}

// DeletePeer Delete the peer from the WG server
func DeletePeer(conn *ssh.Client, peerKey string) error {
	command := "sudo wg set wg0 peer " + peerKey + " remove"

	_, err := network.RunCommand(conn, command)
	return err
}
