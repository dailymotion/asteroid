package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dailymotion/asteroid/pkg/db"
	"github.com/dailymotion/asteroid/pkg/network"
	"github.com/dailymotion/asteroid/pkg/peer"
	"github.com/dailymotion/asteroid/pkg/tools"
)

func main() {
	var err error

	// Init Wireguard to keep all the values in one place
	wireguard, err := tools.InitWG(os.Args)
	if err != nil {
		log.Printf("\nError initializing Wireguard: %v\n", err)
		os.Exit(1)
	}

	// Init DB connection
	DBConn, DBConf, err := tools.InitDB()
	if err != nil {
		log.Printf("\nError initializing DB: %v\n", err)
		os.Exit(1)
	}

	// Check which command has been given
	switch os.Args[1] {
	case "add":

		// Checking if Arguments are correct
		err := tools.CheckArguments(os.Args, "add")
		if err != nil {
			log.Fatalln(err)
		}

		// Check if arguments are empty or haven't all necessary requirements
		err = tools.CheckFlagValid(wireguard, "add")
		if err != nil {
			log.Fatalf("\nerror with arguments: %v\n", err)
		}

		// Connect to the server and get the connection
		conn, err := network.ConnectAndRetrieve(&wireguard, "add")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		// We retrieve all the peer vpn ip to show the new added peer
		tmpListPeers, serverPubKey, err := network.RetrieveIPs(conn)
		if err != nil {
			log.Fatal(err)
		}

		// Retrieving peers from server and checking if the one given already exist on it
		err = tools.RetrieveAndCheckForDouble(DBConn, DBConf, &wireguard, tmpListPeers, serverPubKey)
		if err != nil {
			log.Fatal(err)
		}

		// Add new Peer to the Wireguard server
		if err := peer.AddNewPeer(conn, wireguard); err != nil {
			log.Fatalf("error: %v\n", err)
		} else {
			tools.PrintResult("add", wireguard.PeerKey)
		}

		// Retrieving and match peer with key and cidr into a list of map
		listPeers, err := network.RetrieveAndMatchPeer(conn, DBConn)

		// Showing all the Peers in a nice ASCII table
		fmt.Printf("\n\nPeers informations\n-------------------\n")
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

		// Checking if Arguments are correct
		err := tools.CheckArguments(os.Args, "view")
		if err != nil {
			log.Fatalln(err)
		}

		//Connect to the server and get the connection
		conn, err := network.ConnectAndRetrieve(&wireguard, "view")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		// Retrieving and match peer with key and cidr into a list of map
		listPeers, err := network.RetrieveAndMatchPeer(conn, DBConn)

		fmt.Printf("\n\nPeers informations\n-------------------\n")
		network.ShowListIPs(listPeers)

	case "delete":
		// Checking if Arguments are correct
		err := tools.CheckArguments(os.Args, "delete")
		if err != nil {
			log.Fatalln(err)
		}

		// Check if arguments are empty or haven't all necessary requirements
		err = tools.CheckFlagValid(wireguard, "delete")
		if err != nil {
			log.Fatalf("\nerror with arguments: %v\n", err)
		}

		// Connect to the server and get the connection
		conn, err := network.ConnectAndRetrieve(&wireguard, "delete")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		// Retrieving and match peer with key and cidr into a list of map
		listPeers, err := network.RetrieveAndMatchPeer(conn, DBConn)

		// Checks if peer is present on the server
		ok := tools.CheckIfPresent(listPeers, wireguard.PeerDeleteKey)
		if ok {
			// We use the connection to delete the peer on the Wireguard server
			if err = peer.DeletePeer(conn, wireguard.PeerDeleteKey); err != nil {
				log.Fatalf("error: %v\n", err)
			} else {
				tools.PrintResult("delete", wireguard.PeerDeleteKey)
			}
		} else {
			log.Fatal("key not found on the server")
		}

		// If DB is enabled we delete the peer on the database too
		if DBConf.DBEnabled {
			err = db.DeleteUserInDB(DBConn, wireguard.PeerDeleteKey)
			if err != nil {
				log.Fatalf("\nerror: %v\n", err)
			}
		}

		// Retrieve peers after deletion to make sure it has been deleted
		listPeers, err = network.RetrieveAndMatchPeer(conn, DBConn)

		// Showing all the Peers in a nice ASCII table
		fmt.Printf("\n\nPeers informations\n-------------------\n")
		network.ShowListIPs(listPeers)

	case "-h", "--help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid command.\n\n", os.Args[1])
		flag.Usage()
		os.Exit(2)
	}
}
