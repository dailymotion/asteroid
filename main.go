package main

import (
	"flag"

	"fmt"
	"log"
	"os"

	"github.com/dailymotion/asteroid/internal"
	"github.com/dailymotion/asteroid/network"
	"github.com/dailymotion/asteroid/peer"
)

func flagUsage() {
	emojiPlanet := internal.CreateEmoji()
	fmt.Println("NAME:\n   Asteroid CLI " + emojiPlanet + " - An app to manage peers for Wireguard VPN")
	fmt.Println("AUTHOR:")
	fmt.Println("   Ben Cyril")
    fmt.Printf("USAGE:\n   %s command [OPTIONS] [ARGUMENTS ...]\n", os.Args[0])
	//fmt.Println("VERSION:")
	//fmt.Println("   1.0.0")
	fmt.Println("COMMANDS:")
	fmt.Println("   view     View the peers present on the VPN")
	fmt.Println("   add      Add a new peer on the VPN")
	fmt.Println("   delete   Delete a peer on the VPN")
	fmt.Println("OPTIONS:")
	fmt.Println("    -address \"string\"    New peer address (to use with add)")
	fmt.Println("    -key \"string\"        New peer key (to use with add and delete)")
}


func addUsage() {
	fmt.Printf("\nUSAGE:\n   %s [OPTIONS] [ARGUMENTS ...]\n", "add")
	fmt.Println("OPTIONS:")
	fmt.Println("   -address \"string\"    New peer address (to use with add)")
	fmt.Println("   -key \"string\"        New peer key (to use with add and delete)")

}

func deleteUsage() {
	fmt.Printf("\nUSAGE:\n   %s [OPTIONS] [ARGUMENTS ...]\n", "delete")
	fmt.Println("OPTIONS:")
	fmt.Println("   -key \"string\"        Peer key to delete")
}

func main() {
	// Command Master
	addFlag := flag.NewFlagSet("add", flag.ExitOnError)
	deleteFlag := flag.NewFlagSet("delete", flag.ExitOnError)
	// Subcommands
	peerDeleteKeyFlag := deleteFlag.String("key", "", "Peer key to delete")
	peerKeyFlag := addFlag.String("key", "", "New peer key")
	peerCIDRFlag := addFlag.String("address", "", "New peer address")
	// We change the output of Flag Usage to better show what's the app doing
	flag.Usage = flagUsage

	// Verify that flags has been provided
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Check which command has been given
	switch os.Args[1] {
	case "add":
		// Better output for the add flag
		addFlag.Usage = addUsage
		err := addFlag.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("\nError parsing flag: %v\n", err)
			os.Exit(1)
		}

		// Checking if arguments have been given for the command
		if len(os.Args) <= 5 {
			fmt.Printf("Missing Arguments\n")
			addFlag.Usage()
			os.Exit(2)
		}

		// Check if arguments are empty or haven't all necessary requirements
		err = internal.CheckFlagValid(*peerKeyFlag, *peerCIDRFlag, "add")
		if err != nil {
			fmt.Printf("\nError with arguments: %v\n", err)
			os.Exit(2)
		}

		// Connect to the server and retrieve the conn object
		conn, err := network.ConnectAndRetrieve(*peerCIDRFlag, "add")
		if err != nil {
			fmt.Println("\n/!\\ error:", err)
			os.Exit(1)
		}

		// Add new Peer to the server
		_, stdErr := peer.AddNewPeer(conn, *peerKeyFlag, *peerCIDRFlag)
		if stdErr != "" {
			log.Fatalf("stdErr: %v\n", stdErr)
		} else {
			fmt.Println("Peer added !")
			// We retrieve all the peer vpn ip to show the new added peer
			listPeers, err := network.RetrieveIPs(conn)
			if err != nil {
				fmt.Println("error: ", err)
				os.Exit(1)
			}

			fmt.Printf("\n\nPeers informations\n-------------------\n")
			network.ShowListIPs(listPeers)
		}
	case "view":
		flag.Parse()
		// We alert if arguments are given to the command
		if len(os.Args) > 2 {
			fmt.Printf("View doesn't take options\n\n")
			flag.Usage()
			os.Exit(2)
		}

		conn, err := network.ConnectAndRetrieve("", "view")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		listPeers, err := network.RetrieveIPs(conn)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}

		fmt.Printf("\n\nPeers informations\n-------------------\n")
		network.ShowListIPs(listPeers)

	case "delete":
		deleteFlag.Usage = deleteUsage
		err := deleteFlag.Parse(os.Args[2:])
		if err != nil {
			log.Fatal("Issue with Parse Flag: ", err)
		}
		if len(os.Args) < 3 {
			deleteFlag.Usage()
			os.Exit(2)
		}
		err = internal.CheckFlagValid(*peerDeleteKeyFlag, "", "delete")
		if err != nil {
			fmt.Printf("Error with arguments: %v\n", err)
			deleteFlag.Usage()
			os.Exit(2)
		}
		conn, err := network.ConnectAndRetrieve("", "delete")
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		_, stdErr := peer.DeletePeer(conn, *peerDeleteKeyFlag)
		if stdErr != "" {
			fmt.Printf("stdErr: %v\n", stdErr)
			os.Exit(1)
		} else {
			fmt.Printf("\nPeer %v deleted !\n", *peerDeleteKeyFlag)
			listPeers, err := network.RetrieveIPs(conn)
			if err != nil {
				fmt.Println("error: ", err)
				os.Exit(1)
			}

			network.ShowListIPs(listPeers)
		}
	case "-h", "--help":
		flag.Usage()
	default:
		fmt.Printf("%q is not valid command.\n\n", os.Args[1])
		flag.Usage()
		os.Exit(2)
	}
}