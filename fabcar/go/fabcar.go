/*
Copyright 2020 IBM All Rights Reserved.

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
			fmt.Println("Use this program to interact with deployed chaincode")
			fmt.Println("Currently implemented functions for B/L chaincode are query, queryAll and create")
			fmt.Println("Arguments are:")
			fmt.Println("    queryAll")
			fmt.Println("    query")
			fmt.Println("    create")
			fmt.Println("	 help | h | Help | H")

		}

		if argument == "queryAll" {
			result, err := contract.EvaluateTransaction("queryAllBls")
			if err != nil {
				fmt.Printf("Failed to evaluate transaction: %s\n", err)
				os.Exit(1)
			}
			fmt.Println(string(result))

		} else if argument == "create" {
			if len(os.Args) == 60 {

				//Output loop for testing and monitoring purposes during development
				for i := 0; i < len(os.Args); i++ {
					fmt.Println("[", strconv.Itoa(i), "]", " ", os.Args[i])
				}

				/*createArgBlNr, createArgDoi, createArgPoi, createArgNoBl, createArgShName, createArgShAddr, createArgShCntct, createArgShLF, createArgCoName, createArgCoAddr, createArgCoCntct, createArgCoLF, createArgCrName, createArgCrAddr, createArgCrCntct, createArgCrLF, createArgCrTN, createArgACName, createArgACLF, createArgACAddr, createArgNPCName, createArgNPCAddr, createArgNPCLF, createArgNPCSame, createArgIncoterms, createArgFreightCC, createArgPrpd, createArgCllct, createArgPortOfLoading, createArgPortOfDischarge, createArgPlaceOfReceipt, createArgPlaceOfDelivery, createArgOVN, createArgContNr, createArgFullContLoad, createArgLessThenContLoad, createArgDtOfRecvd, createArgShippedOnBoardDate, createArgMrksNmbrs, createArgNrOfPkg, createArgGrossWeight, createArgGrossWeightUnit, createArgDscrOfGoods*/

				//converting cl input arguments to format with seperator which can be handled by the chaincode
				blData := ""
				sep := "_|_"

				for j := 2; j < len(os.Args); j++ {
					blData += os.Args[j]
					if j != len(os.Args)-1 {
						blData += sep
					}
				}

				fmt.Println(blData)

				result, err = contract.SubmitTransaction("createBl", blData)
				if err != nil {
					fmt.Printf("Failed to submit transaction: %s\n", err)
					os.Exit(1)
				}
				fmt.Println(string(result))
				fmt.Println("B/L with key / BLNumber: ", os.Args[2], "successfully created.")
			} else {
				fmt.Println("59 Arguments expected  |  go run fabcar.go [create {BLNumber} {DateOfIssue} {PlaceOfIssue} {NrOfBLIssued} {ShipperName} {ShipperAddress} {ShipperContact} {ShipperLegalForm} {ConsigneeName} {ConsigneeAddress} {ConsigneeContact} {ConsigneeLegalForm} {CarrierName} {CarrierAddress} {CarrierContact} {CarrierLegalForm} {CarrierTrailerNr} {AgentCompanyName} {AgentCompanyLegalForm} {AgentCompanyAddress} {NotifyPartyCompanyName} {NotifyPartyCompanyAddress} {NotifyPartyCompanyLegalForm} {NotifyPartySameAs} {Incoterms} {FreightChargesCurrency} {Prepaid} {Collect} {PortOfLoading} {PortOfDischarge} {PlaceOfReceipt} {PlaceOfDelivery} {OceanVesselName} {Containernumber} {FullContainerLoad} {LessThenContainerLoad} {DateOfRecieved} {ShippedOnBoardDate} {MarksAndNumbers} {NumberOfPackages} {GrossWeight} {GrossWeightUnit} {DescriptionOfGoods} {DescriptionPerPackage} {Measurement} {MeasurementUnit} {DeclaredCargoValueAmount} {DeclaredCargoValueCurrency} {AdditionalInformation} {HazardousMaterial} {CustomerOrderNumber} {TransportConditions} {ApplieableLaw} {PlaceOfJurisdiction} {CurrentOwner} {OrderBy} {OrderTo} {OrderAt}]")
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
