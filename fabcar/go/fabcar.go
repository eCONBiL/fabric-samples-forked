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
			if len(os.Args) == 57 {

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
				fmt.Println("56 Arguments expected  |  go run fabcar.go [create {BLNumber} {DateOfIssue} {PlaceOfIssue} {ShipperName} {ShipperAddress} {ShipperContact} {ShipperLegalForm} {ConsigneeName} {ConsigneeAddress} {ConsigneeContact} {ConsigneeLegalForm} {CarrierName} {CarrierAddress} {CarrierContact} {CarrierLegalForm} {CarrierTrailerNr} {AgentCompanyName} {AgentCompanyLegalForm} {AgentCompanyAddress} {NotifyPartyCompanyName} {NotifyPartyCompanyAddress} {NotifyPartyCompanyLegalForm} {NotifyPartySameAs} {Incoterms} {FreightChargesCurrency} {Prepaid} {Collect} {PortOfLoading} {PortOfDischarge} {PlaceOfReceipt} {PlaceOfDelivery} {OceanVesselName} {Containernumber} {FullContainerLoad} {LessThenContainerLoad} {ShippedOnBoardDate} {MarksAndNumbers} {NumberOfPackages} {GrossWeight} {GrossWeightUnit} {DescriptionOfGoods} {DescriptionPerPackage} {Measurement} {MeasurementUnit} {DeclaredCargoValueAmount} {DeclaredCargoValueCurrency} {AdditionalInformation} {HazardousMaterial} {CustomerOrderNumber} {TransportConditions} {ApplieableLaw} {PlaceOfJurisdiction} {CurrentOwner} {OrderBy} {OrderTo} {OrderAt}]")
				fmt.Println("Your input: ", os.Args)
				fmt.Println("Your input length: ", len(os.Args))
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

		} else if argument == "endorsement" {
			if len(os.Args) == 8 {
				var blN = os.Args[2]
				var NPCN = os.Args[3]
				var NPCA = os.Args[4]
				var NPCLF = os.Args[5]
				var OT = os.Args[6]
				var OA = os.Args[7]

				fmt.Println(blN, NPCN, NPCA, NPCLF, OT, OA)

				result, err = contract.SubmitTransaction("TransferBl", blN, NPCN, NPCA, NPCLF, OT, OA)
				if err != nil {
					fmt.Printf("Failed to evaluate transaction during endorsement: %s\n", err)
					os.Exit(1)
				}
				fmt.Println("Endorsement successfully executed - new Owner (OrderTo) is : ", OT)
			} else {
				fmt.Println("Wrong input for endorsement |  should be: endorsement blNumber NPCN NPCA NPCLF OT OA")
			}
		} else if argument == "depreciation" {
			if len(os.Args) == 14 {
				var bln = os.Args[2]
				var newNPKG = os.Args[3]
				var newGrossWeight = os.Args[4]
				var newGrossWeightUnit = os.Args[5]
				var newDOG = os.Args[6]
				var newDPP = os.Args[7]
				var newMeasurement = os.Args[8]
				var newMeasurementUnit = os.Args[9]
				var newDCVA = os.Args[10]
				var newDCVC = os.Args[11]
				var newAI = os.Args[12]
				var newHM = os.Args[13]

				fmt.Println(bln, newNPKG, newGrossWeight, newGrossWeightUnit, newDOG, newDPP, newMeasurement, newMeasurementUnit, newDCVA, newDCVC, newAI, newHM)
				result, err = contract.SubmitTransaction("DepreciationOfBl", bln, newNPKG, newGrossWeight, newGrossWeightUnit, newDOG, newDPP, newMeasurement, newMeasurementUnit, newDCVA, newDCVC, newAI, newHM)
				if err != nil {
					fmt.Printf("Failed to evaluate transaction during depreciation: %s\n", err)
					os.Exit(1)
				}
				fmt.Println("Depreciation successfully executed - issue 'go run fabcar.go query ", bln, "' to see the result")

			} else {
				fmt.Println("Wrong input for depreciation |  should be: depreciation blNumber newNrPKG newGrossWeight newGrossWeightUnit newDOG newDPP newMeasurement newMeasurementUnit newDCVA newDCVC newAI newHM")
			}
		} else if argument == "changeOceanVessel" {
			if len(os.Args) == 4 {
				var blnr = os.Args[2]
				var newShip = os.Args[3]
				result, err = contract.SubmitTransaction("ChangeOceanVessel", blnr, newShip)
				if err != nil {
					fmt.Printf("Failed to evaluate transaction during change of ocean vessel name: %s\n", err)
					os.Exit(1)
				}
				fmt.Println("change of ocean vessel name successfully executed  - issue 'go run fabcar.go query ", blnr, "' to see the result")
				fmt.Println(blnr)
			} else {
				fmt.Println("Wrong input for changeOceanVessel | should be: changeOceanVessel blNumber newOceanVesselName")
			}

			// } else if argument == "redirectContainer" {
			// 	if len(os.Args) == 4 {
			// 		var blnr = os.Args[2]
			// 		var newDest = os.Args[3]
			// 		result, err = contract.SubmitTransaction("RedirectContainer", blnr, newDest)
			// 		if err != nil {
			// 			fmt.Printf("Failed to evaluate transaction during redirection of container: %s\n", err)
			// 			os.Exit(1)
			// 		}
			// 		fmt.Println("container redirection successfully executed  - issue 'go run fabcar.go query ", blnr, "' to see the result")
			// 	} else {
			// 		fmt.Println("Wrong input for redirectContainer | should be: redirectContainer blNumber newDestination")
			// 	}

			// } else if argument == "returnWithoutLoading" {
			// 	if len(os.Args) == 3 {
			// 		result, err = contract.SubmitTransaction("ReturnBlWithoutLoading", os.Args[2])
			// 		if err != nil {
			// 			fmt.Printf("Failed to evaluate transaction during return of B/L without loading: %s\n", err)
			// 			os.Exit(1)
			// 		}
			// 		fmt.Println("return without loading successfully executed  - issue 'go run fabcar.go query ", os.Args[2], "' to see the result")
			// 	} else {
			// 		fmt.Println("Wrong input for returnWithoutLoading | should be: returnWithoutLoading blNumber")
			// 	}

		} else if argument == "load" {
			if len(os.Args) == 3 {
				result, err = contract.SubmitTransaction("LoadOnBoard", os.Args[2])
				if err != nil {
					fmt.Printf("Failed to evaluate transaction during loading on board: %s\n", err)
					os.Exit(1)
				}
				fmt.Println("loading on board successfully executed  - issue 'go run fabcar.go query ", os.Args[2], "' to see the result")
			} else {
				fmt.Println("Wrong input for load | should be: load blNumber")
			}
		} else {
			fmt.Println("No valid argument | try: queryAll , create , query, help, endorsement, depreciation, changeOceanVessel, redirectContainer, load, returnWithoutLoading")
		}
	} else {
		fmt.Println("Missing argument(s) | try: queryAll , create , query, help, endorsement, depreciation, changeOceanVessel, redirectContainer, load, returnWithoutLoading")
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
