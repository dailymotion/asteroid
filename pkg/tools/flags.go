package tools

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func flagUsage() {
	emojiPlanet := CreateEmoji()
	fmt.Println("NAME:\n   Asteroid CLI " + emojiPlanet + " - An app to manage peers for Wireguard VPN")
	fmt.Println("AUTHOR:")
	fmt.Println("   Ben Cyril")
	fmt.Printf("USAGE:\n   %s command [OPTIONS] [ARGUMENTS ...]\n", os.Args[0])
	fmt.Println("COMMANDS:")
	fmt.Println("   view     View the peers present on the VPN")
	fmt.Println("   add      Add a new peer on the VPN")
	fmt.Println("   delete   Delete a peer on the VPN")
	fmt.Println("OPTIONS:")
	fmt.Println("   -address \"string\"   New peer address (to use with add)")
	fmt.Println("   -key \"string\"       New peer key (to use with add and delete)")
	fmt.Println("   -name \"string\"       Allows to gives a username for the key (to use with add)")
	fmt.Println("   -generateFile       Generate a file with all the client configurations (to use with add)")
	fmt.Println("   -generateStdout     Generate an output with all the client configurations (to use with add)")
}

func addUsage() {
	fmt.Printf("USAGE:\n   %s [OPTIONS] [ARGUMENTS ...]\n", "add")
	fmt.Println("OPTIONS:")
	fmt.Println("   -address \"string\"    New peer address")
	fmt.Println("   -key \"string\"        New peer key")
	fmt.Println("   -name \"string\"       Allows to gives a username for the key")

}

func deleteUsage() {
	fmt.Printf("USAGE:\n   %s [OPTIONS] [ARGUMENTS ...]\n", "delete")
	fmt.Println("OPTIONS:")
	fmt.Println("   -key \"string\"        Peer key to delete")
}

// HandleFlags init flags and checks which was given
func HandleFlags(args []string) (*string, *string, *string, *string, *bool, *bool) {
	addFlag := flag.NewFlagSet("add", flag.ExitOnError)
	deleteFlag := flag.NewFlagSet("delete", flag.ExitOnError)
	// Flags subcommands
	peerDeleteKeyFlag := deleteFlag.String("key", "", "Peer key to delete")
	peerKeyFlag := addFlag.String("key", "", "New peer key")
	peerCIDRFlag := addFlag.String("address", "", "New peer address")
	peerNameFlag :=  addFlag.String("name", "", "New peer's name")
	generateFile := addFlag.Bool("generateFile", false, "Generate a config file")
	generateOutput := addFlag.Bool("generateStdout", false, "Generate a config output")
	// We change flag usage output to better show what's the app doing
	flag.Usage = flagUsage

	// Verify that flags has been provided
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if args[1] == "add" {
		addFlag.Usage = addUsage
		err := addFlag.Parse(os.Args[2:])
		if err != nil {
			log.Printf("\nError parsing flag: %v\n", err)
			os.Exit(1)
		}
	} else if args[1] == "delete" {
		deleteFlag.Usage = deleteUsage
		err := deleteFlag.Parse(os.Args[2:])
		if err != nil {
			log.Fatal("Issue with Parse Flag: ", err)
		}
	}

	return peerDeleteKeyFlag, peerKeyFlag, peerCIDRFlag, peerNameFlag, generateFile, generateOutput
}
