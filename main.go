package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	Handlers "github.com/igorcrevar/simple-eth-api-client/handlers"
)

func main() {
	var clientUrl string
	var privateKey string

	args := os.Args[1:]
	switch ln := len(args); {
	case ln == 1:
		clientUrl = args[0]
	case ln > 1:
		clientUrl = args[0]
		privateKey = args[1]
	default:
		clientUrl = "http://localhost:7545" // if you want to use ganache @see https://www.trufflesuite.com/ganache
	}

	// Create a client instance to connect to our provider
	client, err := ethclient.Dial(clientUrl)

	if err != nil {
		fmt.Println(err)
	}

	Handlers.Register(client, privateKey)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
