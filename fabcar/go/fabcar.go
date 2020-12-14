/*
Copyright 2020 IBM All Rights Reserved.

B/L  chaincode interface program for B/L chaincode Beta Version 1.0

fabcar.go main program containing some chaincode interaction commands for testing written BL chaincode locally

edited for using command line arguments for chaincode interaction

Edited by Malte Garmhausen

eCONBiL - electronical Consignment and Bill of Lading

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

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

	// First query, output (var result) is currently cut out
	result, err := contract.EvaluateTransaction("QueryAllBls")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	// fmt.Println(string(result))

	//the following block implements chaicode invokes and querys according to given arguments
	if len(os.Args) > 1 {

		var argument string = os.Args[1]

		if (argument == "help") || (argument == "H") || (argument == "Help") || (argument == "h") {
			fmt.Println()
			fmt.Println("B/L chaincode Beta Version 1.0 command line interface program")
			fmt.Println()
			fmt.Println("Use this program to interact with deployed chaincode")
			fmt.Println("	--> the current chaincode beta version 1.0ß used chaincode is located in /fabric-samples/chaincode/fabcar/go/fabcar.go")
			fmt.Println("Currently implemented functions for B/L chaincode are query, queryAll and create")
			fmt.Println("Arguments are:")
			fmt.Println("    queryAll")
			fmt.Println("    query")
			fmt.Println("    create")
			fmt.Println("	 help | h | Help | H")

		} else if argument == "queryAll" {
			result, err := contract.EvaluateTransaction("queryAllBls")
			if err != nil {
				fmt.Printf("Failed to evaluate transaction: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(string(result))

		} else if argument == "create" {
			if len(os.Args) == 45 {
				// ==> pos.Args[0] = /tm/go-build... --> default while reading arguments from os input
				// ==> os.Args[1] = "create"
				var createArgBlNr string = os.Args[2]
				var createArgDoi string = os.Args[3]
				var createArgPoi string = os.Args[4]
				var createArgNoBl string = os.Args[5]
				var createArgShName string = os.Args[6]
				var createArgShAddr string = os.Args[7]
				var createArgShCntct string = os.Args[8]
				var createArgShLF string = os.Args[9]
				var createArgCoName string = os.Args[10]
				var createArgCoAddr string = os.Args[11]
				var createArgCoCntct string = os.Args[12]
				var createArgCoLF string = os.Args[13]
				var createArgCrName string = os.Args[14]
				var createArgCrAddr string = os.Args[15]
				var createArgCrCntct string = os.Args[16]
				var createArgCrLF string = os.Args[17]
				var createArgCrTN string = os.Args[18]
				var createArgACName string = os.Args[19]
				var createArgACLF string = os.Args[20]
				var createArgACAddr string = os.Args[21]
				var createArgNPCName string = os.Args[22]
				var createArgNPCAddr string = os.Args[23]
				var createArgNPCLF string = os.Args[24]
				var createArgNPCSame string = os.Args[25]
				var createArgIncoterms string = os.Args[26]
				var createArgFreightCC string = os.Args[27]
				var createArgPrpd string = os.Args[28]
				var createArgCllct string = os.Args[29]
				var createArgPortOfLoading string = os.Args[30]
				var createArgPortOfDischarge string = os.Args[31]
				var createArgPlaceOfReceipt string = os.Args[32]
				var createArgPlaceOfDelivery string = os.Args[33]
				var createArgOVN string = os.Args[34]
				var createArgContNr string = os.Args[35]
				var createArgFullContLoad string = os.Args[36]
				var createArgLessThenContLoad string = os.Args[37]
				var createArgDtOfRecvd string = os.Args[38]
				var createArgShippedOnBoardDate string = os.Args[39]
				var createArgMrksNmbrs string = os.Args[40]
				var createArgNrOfPkg string = os.Args[41]
				var createArgGrossWeight string = os.Args[42]
				var createArgGrossWeightUnit string = os.Args[43]
				var createArgDscrOfGoods string = os.Args[44]

				//Output loop for testing and monitoring purposes during development
				for i := 0; i < 45; i++ {
					fmt.Println("[", strconv.Itoa(i), "]", " ", os.Args[i])
				}
				// fmt.Println(createArgGrossWeightUnit)
				// fmt.Println(os.Args[44])

				fmt.Println("Das B/L mit dem Schlüssel bzw. der BLNumber: ", createArgBlNr, " wurde im Ledger angelegt.")

				result, err = contract.SubmitTransaction("createBl", createArgBlNr, createArgDoi, createArgPoi, createArgNoBl, createArgShName, createArgShAddr, createArgShCntct, createArgShLF, createArgCoName, createArgCoAddr, createArgCoCntct, createArgCoLF, createArgCrName, createArgCrAddr, createArgCrCntct, createArgCrLF, createArgCrTN, createArgACName, createArgACLF, createArgACAddr, createArgNPCName, createArgNPCAddr, createArgNPCLF, createArgNPCSame, createArgIncoterms, createArgFreightCC, createArgPrpd, createArgCllct, createArgPortOfLoading, createArgPortOfDischarge, createArgPlaceOfReceipt, createArgPlaceOfDelivery, createArgOVN, createArgContNr, createArgFullContLoad, createArgLessThenContLoad, createArgDtOfRecvd, createArgShippedOnBoardDate, createArgMrksNmbrs, createArgNrOfPkg, createArgGrossWeight, createArgGrossWeightUnit, createArgDscrOfGoods)
				if err != nil {
					fmt.Printf("Failed to submit transaction: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(result))

			} else {
				fmt.Println("44 Arguments expected  |  go run fabcar.go [create {BLNumber} {DateOfIssue} {PlaceOfIssue} {NrOfBLIssued} {ShipperName} {ShipperAddress} {ShipperContact} {ShipperLegalForm} {ConsigneeName} {ConsigneeAddress} {ConsigneeContact} {ConsigneeLegalForm} {CarrierName} {CarrierAddress} {CarrierContact} {CarrierLegalForm} {CarrierTrailerNr} {AgentCompanyName} {AgentCompanyAddress} {NotifyPartyCompanyName} {NotifyPartyCompanyAddress} {NotifyPartyCompanyLegalForm} {NotifyPartySameAs} {Incoterms} {FreightChargesCurrency} {Prepaid} {Collect} {PortOfLoading} {PortOfDischarge} {PlaceOfReceipt} {PlaceOfDelivery} {OceanVesselName} {Containernumber} {FullContainerLoad} {LessThenContainerLoad} {DateOfRecieved} {ShippedOnBoardDate} {MarksAndNumbers} {NumberOfPackages} {GrossWeight} {GrossWeightUnit} {DescriptionOfGoods}]")

				//Output loop for testing and monitoring purposes during development
				// for i := 0; i < 45; i++ {
				// 	fmt.Println("[", strconv.Itoa(i), "]", " ", os.Args[i])
				// }
				// fmt.Println("length of args: ", len(os.Args))
			}
		} else if argument == "query" {
			if len(os.Args) == 3 {
				var queryArg string = os.Args[2]
				result, err = contract.EvaluateTransaction("queryBl", queryArg)
				if err != nil {
					fmt.Printf("Failed to evaluate transaction: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(result))

			} else {
				fmt.Println("2 Arguments expected  |  go run fabcar-mod.go [query {BLNumber}]")
			}
			// } else if argument == "changeOwner" {
			// 	if len(os.Args) == 4 {
			// 		var changeArgID string = os.Args[2]
			// 		var changeArgOw string = os.Args[3]
			// 		_, err = contract.SubmitTransaction("changeBlOwner", changeArgID, changeArgOw)
			// 		if err != nil {
			// 			fmt.Printf("Failed to submit transaction: %s\n", err)
			// 			os.Exit(1)
			// 		}
			// 	} else {
			// 		fmt.Println("3 Arguments expected  |  go run fabcar-mod.go [changeOwner {BLID} {NEWBLOWNER}]")
			// 	}
		} else {
			fmt.Println("No valid argument | try: queryAll , create , query, help")
		}
	} else {
		fmt.Println("Missing argument(s) | try: queryAll , create , query, help")
	}
}

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
