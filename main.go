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
	fmt.Println("    -address \"string\"   New peer address (to use with add)")
	fmt.Println("    -key \"string\"       New peer key (to use with add and delete)")
	fmt.Println("    -generateFile       Generate a file with all the client configurations")
	fmt.Println("    -generateStdout     Generate an output with all the client configurations")
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
	generateFile := addFlag.Bool("generateFile", false, "Generate a config file")
	generateOutput := addFlag.Bool("generateStdout", false, "Generate a config output")
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
			log.Printf("\nError parsing flag: %v\n", err)
			os.Exit(1)
		}

		// Checking if arguments have been given for the command
		if len(os.Args) <= 5 {
			log.Printf("Missing Arguments\n")
			addFlag.Usage()
			os.Exit(2)
		}

		// Check if arguments are empty or haven't all necessary requirements
		//err = internal.CheckFlagValid(*peerKeyFlag, *peerCIDRFlag, "add")
		//if err != nil {
		//	log.Fatalf("\nerror with arguments: %v\n", err)
		//}

		// Connect to the server and retrieve the conn object
		//conn, err := network.ConnectAndRetrieve(*peerCIDRFlag, "add")
		//if err != nil {
		//	log.Fatalf("\nerror: %v\n", err)
		//}

		// Add new Peer to the server
		//if err := peer.AddNewPeer(conn, *peerKeyFlag, *peerCIDRFlag); err != nil {
		//	log.Fatalf("error: %v\n", err)
		//}

		fmt.Println("Peer added !")
		// We retrieve all the peer vpn ip to show the new added peer
		//listPeers, err := network.RetrieveIPs(conn)
		//if err != nil {
		//	fmt.Println("error: ", err)
		//	os.Exit(1)
		//}

		//fmt.Printf("\n\nPeers informations\n-------------------\n")
		//network.ShowListIPs(listPeers)
		fmt.Println("generateOutput: ", *generateOutput)
		if *generateFile || *generateOutput {
			err := internal.RetrieveWGConfig(*generateFile, *generateOutput, *peerKeyFlag, *peerCIDRFlag)
			if err != nil {
				log.Println(err)
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

		conn, err := network.ConnectAndRetrieve("", "view")
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
		}

		listPeers, err := network.RetrieveIPs(conn)
		if err != nil {
			log.Fatalf("\nerror: %v\n", err)
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
			log.Fatalf("\nerror: %v\n", err)
		}
		if err = peer.DeletePeer(conn, *peerDeleteKeyFlag); err != nil {
			log.Fatalf("error: %v\n", err)
		}

		fmt.Printf("\nPeer %v deleted !\n", *peerDeleteKeyFlag)
		listPeers, err := network.RetrieveIPs(conn)
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
