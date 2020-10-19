package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dailymotion/asteroid/pkg/network"
	"github.com/dailymotion/asteroid/pkg/peer"
	"github.com/dailymotion/asteroid/pkg/tools"
)

func main() {
	var err error

	// Init wireguard to keep all the value in one place
	wireguard, err := tools.InitWG(os.Args)
	if err != nil {
		log.Printf("\nError parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Check which command has been given
	switch os.Args[1] {
	case "add":

		// Checking if arguments have been given for the command
		if len(os.Args) <= 5 {
			log.Printf("Missing Arguments\n")
			//addFlag.Usage()
			os.Exit(2)
		}

		// Check if arguments are empty or haven't all necessary requirements
		err = tools.CheckFlagValid(wireguard, "add")
		if err != nil {
			log.Fatalf("\nerror with arguments: %v\n", err)
		}

		// Connect to the server and retrieve the conn object
		conn, err := network.ConnectAndRetrieve(wireguard, "add")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		// Add new Peer to the server
		if err := peer.AddNewPeer(conn, wireguard); err != nil {
			log.Fatalf("error: %v\n", err)
		} else {
		// Message be like:
		//################
		//# Peer added ! #
		//################
		fmt.Printf("\n################\n# Peer added ! #\n################\n\n")
		}

		// We retrieve all the peer vpn ip to show the new added peer
		listPeers, serverPubKey, err := network.RetrieveIPs(conn)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}

		// Adding server public key to the wireguard object
		wireguard.ServerPubKey = serverPubKey

		//fmt.Printf("\n\nPeers informations\n-------------------\n")
		network.ShowListIPs(listPeers)

		// We check that one of the flag is true
		if wireguard.GenerateFile || wireguard.GenerateOutput {
			err := wireguard.RetrieveWGConfig()
			if err != nil {
				log.Fatalf("\nerror: %v\n", err)
			}
		}
	case "view":
		flag.Parse()
		// We alert if arguments are given to the command
		if len(os.Args) > 2 {
			fmt.Printf("View doesn't take options\n\n")
			flag.Usage()
			os.Exit(2)
		}

		conn, err := network.ConnectAndRetrieve(wireguard, "view")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		listPeers, _, err := network.RetrieveIPs(conn)
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		fmt.Printf("\n\nPeers informations\n-------------------\n")
		network.ShowListIPs(listPeers)

	case "delete":
		if len(os.Args) < 3 {
			//deleteFlag.Usage()
			os.Exit(2)
		}

		err = tools.CheckFlagValid(wireguard, "delete")
		if err != nil {
			fmt.Printf("Error with arguments: %v\n", err)
			//deleteFlag.Usage()
			os.Exit(2)
		}
		conn, err := network.ConnectAndRetrieve(wireguard, "delete")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}
		if err = peer.DeletePeer(conn, wireguard.PeerDeleteKey); err != nil {
			log.Fatalf("error: %v\n", err)
		}

		fmt.Printf("\nPeer %v deleted !\n", wireguard.PeerDeleteKey)
		listPeers, _, err := network.RetrieveIPs(conn)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}

		network.ShowListIPs(listPeers)

	case "-h", "--help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid command.\n\n", os.Args[1])
		flag.Usage()
		os.Exit(2)
	}
}
