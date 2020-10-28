/*
Copyright 2020 IBM All Rights Reserved.

fabcar.go main program containing some chaincode interaction commands for testing written BL chaincode locally

edited for using command line arguments for chaincode interaction 

Edited by Malte Garmhausen

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("fabcar")

	result, err := contract.EvaluateTransaction("QueryAllBls")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	if len(os.Args) > 1 {

		var argument string = os.Args[1]
		
		if ( argument == "help" ) || ( argument == "H" ) || (argument == "Help" ) || ( argument == "h" ) {
			fmt.Println()
			fmt.Println("Use this program to interact with deployed chaincode")
			fmt.Println("Arguments are:")
			fmt.Println("    queryAll")
			fmt.Println("    query")
			fmt.Println("    create")
			fmt.Println("    changeOwner")

		}	

		if argument == "queryAll" {
		result, err := contract.EvaluateTransaction("queryAllBls")
		if err != nil {
			fmt.Printf("Failed to evaluate transaction: %s\n", err)
			os.Exit(1)
			}
			fmt.Println(string(result))

		} else if  argument == "create" {
			if len(os.Args) == 7 {
				var createArgId string = os.Args[2]
				var createArgSh string = os.Args[3]
				var createArgDoi string = os.Args[4]
				var createArgPoi string = os.Args[5]
				var createArgOw string = os.Args[6]
				result, err = contract.SubmitTransaction("createBl", createArgId, createArgSh, createArgDoi, createArgPoi, createArgOw)
				if err != nil {
					fmt.Printf("Failed to submit transaction: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(result))

			} else {
				fmt.Println("6 Arguments expected  |  go run fabcar-mod.go [create {BLID} {BLSHIPPER} {DATEOFISSUE} {PLACEOFISSUE} {BLOWNER}]")
			}
		} else if  argument == "query" {
			if len(os.Args) == 3 {
				var queryArg string = os.Args[2]
				result, err = contract.EvaluateTransaction("queryBl", queryArg)
				if err != nil {
					fmt.Printf("Failed to evaluate transaction: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(result))

			} else {
				fmt.Println("2 Arguments expected  |  go run fabcar-mod.go [query {BLID}]")
			}
		} else if argument == "changeOwner" {
			if len(os.Args) == 4 {
				var changeArgId string = os.Args[2]
				var changeArgOw string = os.Args[3]
				_, err = contract.SubmitTransaction("changeBlOwner", changeArgId, changeArgOw)
				if err != nil {
					fmt.Printf("Failed to submit transaction: %s\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("3 Arguments expected  |  go run fabcar-mod.go [changeOwner {BLID} {NEWBLOWNER}]")
			}
		} else {
			fmt.Println("No valid argument | try: queryAll , create , query , changeOwner")
		}
	} else {
		fmt.Println("Missing argument(s) | try: queryAll , create , query , changeOwner")
	}
}

/*
	result, err = contract.SubmitTransaction("CreateBl", "CAR10", "HS Bremerhaven", "21.10.2020", "Bremerhaven", "Karin Vosseberg")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.EvaluateTransaction("QueryBl", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	_, err = contract.SubmitTransaction("ChangeBlOwner", "CAR10", "Peter Kelb")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}

	result, err = contract.EvaluateTransaction("QueryBl", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.EvaluateTransaction("QueryAllBls")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
}

*/

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}
