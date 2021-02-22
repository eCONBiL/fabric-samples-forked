/*
SPDX-License-Identifier: Apache-2.0
*/

/*
modified by Malte Garmhausen

eCONBiL - electronical Consignment and Bill of Lading

changes made for an implementation of the Bill of Lading freight document

B/L chaincode Beta Version 1.1 (work in progress)

This chaincode has all of the collected BL fields implemented and provides functions for initializing the ledger, querying specific B/Ls, querying
all B/Ls from the ledger as well as creating a new B/L.
The initialization-function creates one default B/L into the ledger.

Local Chaincode querys and invokes must be issued through the "interface program" found in /fabric-samples/fabcar/go/fabcar.go which is currently used as
command line interface

*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a B/L
type SmartContract struct {
	contractapi.Contract
}

//BillOfLading struct contains all BL-Data
type BillOfLading struct {

	// BL Info Fields
	BLNumber     string `json:"BLNumber"`
	DateOfIssue  string `json:"DateOfIssue"`
	PlaceOfIssue string `json:"PlaceOfIssue"`

	//Shipper information
	ShipperName      string `json:"ShipperName"` //5
	ShipperAddress   string `json:"ShipperAddress"`
	ShipperContact   string `json:"ShipperContact"`
	ShipperLegalForm string `json:"ShipperLegalForm"`

	//Consignee information
	ConsigneeName      string `json:"ConsigneeName"`
	ConsigneeAddress   string `json:"ConsigneeAddress"` //10
	ConsigneeContact   string `json:"ConsigneeContact"`
	ConsigneeLegalForm string `json:"ConsigneeLegalForm"`

	//Carrier information
	CarrierName          string `json:"CarrierName"`
	CarrierAddress       string `json:"CarrierAddress"`
	CarrierContact       string `json:"CarrierContact"` //15
	CarrierLegalForm     string `json:"CarrierLegalForm"`
	CarrierTrailerNumber string `json:"CarrierTrailerNumber"`

	//Forwarding Agent Company information
	AgentCompanyName      string `json:"AgentCompanyName"`
	AgentCompanyLegalForm string `json:"AgentCompanyLegalForm"`
	AgentCompanyAddress   string `json:"AgentCompanyAddress"` //20

	//Notify Party information
	NotifyPartyCompanyName      string `json:"NotifyPartyCompanyName"`
	NotifyPartyCompanyAddress   string `json:"NotifyPartyCompanyAddress"`
	NotifyPartyCompanyLegalForm string `json:"NotifyPartyCompanyLegalForm"`

	//Term of sale information
	Incoterms string `json:"Incoterms"`

	//Basic freight information
	FreightChargesCurrency string `json:"FreightChargesCurrency"` //25
	Prepaid                bool   `json:"Prepaid"`
	Collect                bool   `json:"Collect"`

	//Transportinfo
	PortOfLoading         string `json:"PortOfLoading"`
	PortOfDischarge       string `json:"PortOfDischarge"`
	PlaceOfReceipt        string `json:"PlaceOfReceipt"` //30
	PlaceOfDelivery       string `json:"PlaceOfDelivery"`
	OceanVesselName       string `json:"OceanVesselName"`
	ContainerNumber       string `json:"ContainerNumber"`
	FullContainerLoad     bool   `json:"FullContainerLoad"`
	LessThenContainerLoad bool   `json:"LessThenContainerLoad"` //35
	ShippedOnBoardDate    string `json:"ShippedOnBoardDate"`

	//Gross info
	MarksAndNumbers            string  `json:"MarksAndNumbers"`
	NumberOfPackages           int     `json:"NumberOfPackages"`
	GrossWeight                int     `json:"GrossWeight"` //40
	GrossWeightUnit            string  `json:"GrossWeightUnit"`
	DescriptionOfGoods         string  `json:"DescriptionOfGoods"`
	DescriptionPerPackage      string  `json:"DescriptionPerPackage"`
	Measurement                float64 `json:"Measurement"`
	MeasurementUnit            string  `json:"MeasurementUnit"` //45
	DeclaredCargoValueAmount   int     `json:"DeclaredCargoValueAmount"`
	DeclaredCargoValueCurrency string  `json:"DeclaredCargoValueCurrency"`
	AdditionalInformation      string  `json:"AdditionalInformation"`
	HazardousMaterial          bool    `json:"HazardousMaterial"`

	// CustomerOrderNumber
	CustomerOrderNumber string `json:"CustomerOrderNumber"` //50

	//Used Conditions (ERA600, Art. 20 a)
	TransportConditions string `json:"TransportConditions"`
	ApplieableLaw       string `json:"ApplieableLaw"`
	PlaceOfJurisdiction string `json:"PlaceOfJurisdiction"`

	//Endorsement info
	OrderDate      string `json:"OrderDate"`
	OrderTo        string `json:"OrderTo"` //55
	OrderAt        string `json:"OrderAt"`
	OrderCheckbox  bool   `json:"OrderCheckbox"`
	BlTransferable bool   `json: "BlTransferable"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *BillOfLading
}

// InitLedger adds a base set of BL to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	bls := []BillOfLading{
		BillOfLading{BLNumber: "TW ECON 1000", DateOfIssue: "10.11.2020", PlaceOfIssue: "Bremerhaven",
			ShipperName: "Autohaus Main GmbH", ShipperAddress: "Hanauerlandstr. 3460314 Frankfurt, Germany", ShipperContact: "ahmain@beispiel.de", ShipperLegalForm: "GmbH",
			ConsigneeName: "German-Cars Ldt.", ConsigneeAddress: "Fue Avenue, A1 518108 Shanghai, China", ConsigneeContact: "86282452253", ConsigneeLegalForm: "Ldt.",
			CarrierName: "MSC Germany S.A. & Co. KG", CarrierAddress: "Hafenstraße 55, 282127 Bremen, Germany", CarrierContact: "deu-bremen@msc.de", CarrierLegalForm: "S.A. & Co. KG", CarrierTrailerNumber: "HB-KK-596",
			AgentCompanyName: "BLG AutoTerminal Bremerhaven GmbH & Co. KG", AgentCompanyLegalForm: "GmbH & Co. KG", AgentCompanyAddress: "Senator-Borttscheller-Str. 1, 27568 Bremerhaven, Germany",
			NotifyPartyCompanyName: "German-Cars Ldt.", NotifyPartyCompanyAddress: "Fue Avenue, A1 518108 Shanghai, China", NotifyPartyCompanyLegalForm: "Ldt.",
			Incoterms:              "FOB (2020)",
			FreightChargesCurrency: "USD", Prepaid: true, Collect: true,
			PortOfLoading: "Bremerhaven Containerterminal", PortOfDischarge: "Shanghai Yangshan", PlaceOfReceipt: "Frankfurt am Main, Adresse, Germany", PlaceOfDelivery: "Shanghai, Adresse, China", OceanVesselName: "MSC Gulsun", ContainerNumber: "OOLU1548378", FullContainerLoad: true, LessThenContainerLoad: false, ShippedOnBoardDate: "02.02.2021",
			MarksAndNumbers: "40' steel Dry Cargo Container No CSQU3054383", NumberOfPackages: 15, GrossWeight: 4250, GrossWeightUnit: "Kg", DescriptionOfGoods: "engines and fitting engine parts packaged together on pallets", DescriptionPerPackage: "abc", Measurement: 40.2, MeasurementUnit: "Feet", DeclaredCargoValueAmount: 75000, DeclaredCargoValueCurrency: "USD", AdditionalInformation: "-", HazardousMaterial: false,
			CustomerOrderNumber: "CO NR 1234",
			TransportConditions: "TransportCond", ApplieableLaw: "ApplLaw", PlaceOfJurisdiction: "POJ",
			OrderDate:      "",
			OrderTo:        "",
			OrderAt:        "",
			OrderCheckbox:  false,
			BlTransferable: false,
		},
	}

	for i, bl := range bls {
		blAsBytes, _ := json.Marshal(bl)
		err := ctx.GetStub().PutState("TW ECON 100"+strconv.Itoa(i), blAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		fmt.Printf("Ledger wurde initialisiert!")
	}

	return nil
}

// CreateBl adds a new BL to the world state with given details
func (s *SmartContract) CreateBl(ctx contractapi.TransactionContextInterface, BLdata string) error {

	// spliting recieved BL-data string in substrings at every "_|_"-seperator  --> output as an array

	splitResult := strings.Split(BLdata, "_|_")

	// check if B/L with given key already exists in ledger
	keyExisting, keyErr := ctx.GetStub().GetState(splitResult[0])
	if keyExisting != nil && keyErr == nil {
		return fmt.Errorf("BL with given key / BlNumber already exists in ledger and will not be created")
	}

	// converting passed arguments from string to correct data type in the following block

	//prepaid and collect are preset for letter of credit
	convPrepaid, converr := strconv.ParseBool(splitResult[24])
	if converr != nil {
		return fmt.Errorf("Error while converting Prepaid to bool", converr.Error())
	}

	convCollect, converr := strconv.ParseBool(splitResult[25])
	if converr != nil {
		return fmt.Errorf("Error while converting Collect to bool", converr.Error())
	}

	convFullContainerLoad, converr := strconv.ParseBool(splitResult[32])
	if converr != nil {
		return fmt.Errorf("Error while converting FullContainerLoad to bool", converr.Error())
	}

	convLessThenContainerLoad, converr := strconv.ParseBool(splitResult[33])
	if converr != nil {
		return fmt.Errorf("Error while converting LessThenContainerLoad to bool", converr.Error())
	}

	convNumberOfPackages, converr := strconv.Atoi(splitResult[36])
	if converr != nil {
		return fmt.Errorf("Error while converting NumberOfPackages to int", converr.Error())
	}

	convGrossWeight, converr := strconv.Atoi(splitResult[37])
	if converr != nil {
		return fmt.Errorf("Error while converting GrossWeight to int", converr.Error())
	}

	convMeasurement, converr := strconv.ParseFloat(splitResult[41], 64)
	if converr != nil {
		return fmt.Errorf("Error while converting Measurement to float", converr.Error())
	}

	convDeclaredCargoValueAmount, converr := strconv.Atoi(splitResult[43])
	if converr != nil {
		return fmt.Errorf("Error while converting DeclaredCargoValueAmount to int", converr.Error())
	}

	convHazardousMaterial, converr := strconv.ParseBool(splitResult[46])
	if converr != nil {
		return fmt.Errorf("Error while converting HazardousMaterial to bool", converr.Error())
	}

	// convCustomerOrderNumber, converr := strconv.Atoi(splitResult[48])
	// if converr != nil {
	// 	return fmt.Errorf("Error while converting aCustomerOrderNumber to int", converr.Error())
	// }

	convOrderCheckbox, converr := strconv.ParseBool(splitResult[54])
	if converr != nil {
		return fmt.Errorf("Error while converting OrderCheckbox to bool", converr.Error())
	}

	// creation of new bl with all fields being initialization of the passed arguments as field data after some fields data type has been converted
	bl := BillOfLading{
		BLNumber:                    splitResult[0],
		DateOfIssue:                 splitResult[1],
		PlaceOfIssue:                splitResult[2],
		ShipperName:                 splitResult[3],
		ShipperAddress:              splitResult[4],
		ShipperContact:              splitResult[5],
		ShipperLegalForm:            splitResult[6],
		ConsigneeName:               splitResult[7],
		ConsigneeAddress:            splitResult[8],
		ConsigneeContact:            splitResult[9],
		ConsigneeLegalForm:          splitResult[10],
		CarrierName:                 splitResult[11],
		CarrierAddress:              splitResult[12],
		CarrierContact:              splitResult[13],
		CarrierLegalForm:            splitResult[14],
		CarrierTrailerNumber:        splitResult[15],
		AgentCompanyName:            splitResult[16],
		AgentCompanyLegalForm:       splitResult[17],
		AgentCompanyAddress:         splitResult[18],
		NotifyPartyCompanyName:      splitResult[19],
		NotifyPartyCompanyAddress:   splitResult[20],
		NotifyPartyCompanyLegalForm: splitResult[21],
		Incoterms:                   splitResult[22],
		FreightChargesCurrency:      splitResult[23],
		Prepaid:                     convPrepaid, // with converted datatype ParseBool (string to bool)
		Collect:                     convCollect, // with converted datatype ParseBool (string to bool)
		PortOfLoading:               splitResult[26],
		PortOfDischarge:             splitResult[27],
		PlaceOfReceipt:              splitResult[28],
		PlaceOfDelivery:             splitResult[29],
		OceanVesselName:             splitResult[30],
		ContainerNumber:             splitResult[31],
		FullContainerLoad:           convFullContainerLoad,     // with converted datatype ParseBool (string to bool)
		LessThenContainerLoad:       convLessThenContainerLoad, // with converted datatype ParseBool (string to bool)
		ShippedOnBoardDate:          splitResult[34],           //splitResult[36],
		MarksAndNumbers:             splitResult[35],
		NumberOfPackages:            convNumberOfPackages, // with converted datatype Atoi (string to int)
		GrossWeight:                 convGrossWeight,      // with converted datatype Atoi (string to int)
		GrossWeightUnit:             splitResult[38],
		DescriptionOfGoods:          splitResult[39],
		DescriptionPerPackage:       splitResult[40],
		Measurement:                 convMeasurement, // with converted datatype ParseFloat (string to float32)
		MeasurementUnit:             splitResult[42],
		DeclaredCargoValueAmount:    convDeclaredCargoValueAmount, // with converted datatype atoi (string to int)
		DeclaredCargoValueCurrency:  splitResult[44],
		AdditionalInformation:       splitResult[45],
		HazardousMaterial:           convHazardousMaterial, // with converted datatype parseBool (string to bool)
		CustomerOrderNumber:         splitResult[47],       // with converted datatype atoi (string to int)
		TransportConditions:         splitResult[48],
		ApplieableLaw:               splitResult[49],
		PlaceOfJurisdiction:         splitResult[50],
		OrderDate:                   splitResult[51], //should be empty in first creation because no endorsement has been issued yet
		OrderTo:                     splitResult[52], //should be empty in first creation because no endorsement has been issued yet
		OrderAt:                     splitResult[53], //should be empty in first creation because no endorsement has been issued yet
		OrderCheckbox:               convOrderCheckbox,
		BlTransferable:              true,
	}

	if bl.OrderCheckbox == false {
		bl.BlTransferable = false
	}

	blAsBytes, _ := json.Marshal(bl)

	//write Bl in correct format on blockchain

	return ctx.GetStub().PutState(splitResult[0], blAsBytes)

}

// QueryBl returns the B/L stored in the world state with given id
func (s *SmartContract) QueryBl(ctx contractapi.TransactionContextInterface, blNumber string) (*BillOfLading, error) {
	blAsBytes, err := ctx.GetStub().GetState(blNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if blAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", blNumber)
	}

	bl := new(BillOfLading)
	_ = json.Unmarshal(blAsBytes, bl)

	return bl, nil
}

// QueryAllBls returns all B/L found in world state
func (s *SmartContract) QueryAllBls(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		bl := new(BillOfLading)
		_ = json.Unmarshal(queryResponse.Value, bl)

		queryResult := QueryResult{Key: queryResponse.Key, Record: bl}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeOceanVessel function to change the OceanVesselName
func (s *SmartContract) ChangeOceanVessel(ctx contractapi.TransactionContextInterface, blNumber string, newOceanVesselName string) error {
	bl, err := s.QueryBl(ctx, blNumber)

	if err != nil {
		return err
	}

	bl.OceanVesselName = newOceanVesselName

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)

}

//RedirectContainer function for redirecting a container to a new destination
// func (s *SmartContract) RedirectContainer(ctx contractapi.TransactionContextInterface, blNumber string, newDestination string) error {

// 	bl, err := s.QueryBl(ctx, blNumber)

// 	if err != nil {
// 		return err
// 	}
// 	bl.PlaceOfDelivery = newDestination

// 	blAsBytes, _ := json.Marshal(bl)

// 	return ctx.GetStub().PutState(blNumber, blAsBytes)
// }

// ReturnBlWithoutLoading kann in WebUi als bool bzw. Checkbox gebaut werden
// ReturnBlWithoutLoading sets CargoReceivedDate as return of container before loading it on a ship
// func (s *SmartContract) ReturnBlWithoutLoading(ctx contractapi.TransactionContextInterface, blNumber string) error {

// 	bl, err := s.QueryBl(ctx, blNumber)

// 	if err != nil {
// 		return err
// 	}

// 	currentDate := time.Now()
// 	formattedDate := currentDate.Format("02.01.2006 15:04:05") //formatted in "DD.MM.YYYY hh:mm:ss"

// 	bl.CargoReceivedDate = formattedDate

// 	blAsBytes, _ := json.Marshal(bl)

// 	return ctx.GetStub().PutState(blNumber, blAsBytes)

// }

// LoadOnBoard kann in WebUi als bool bzw. Checkbox gebaut werden
// LoadOnBoard sets the ShippedOnBoardDate
func (s *SmartContract) LoadOnBoard(ctx contractapi.TransactionContextInterface, blNumber string) error {

	bl, err := s.QueryBl(ctx, blNumber)

	if err != nil {
		return err
	}

	if bl.ShippedOnBoardDate != "" {
		return fmt.Errorf("Freight has already been shipped on Board, ShippedOnBoardDate is already set - ShippedOnBoardDate: ", bl.ShippedOnBoardDate)
	}

	currentDate := time.Now()
	formattedDate := currentDate.Format("02.01.2006 15:04:05") // formatted in "DD.MM.YYYY hh:mm:ss"

	bl.ShippedOnBoardDate = formattedDate

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)

}

//TransferBl function for transfering the Bl (endorsement)
func (s *SmartContract) TransferBl(ctx contractapi.TransactionContextInterface, blNumber string, NPCN string, NPCA string, NPCLF string, OT string, OA string) error {

	bl, err := s.QueryBl(ctx, blNumber)

	if err != nil {
		return err
	}

	//check if orderAt party is authorized to issue endorsement
	// if bl.OrderTo == "" {
	// 	if OA != bl.CarrierName {
	// 		return fmt.Errorf("only the carrier can issue the first endorsement. Last argument must be equal to CarrierName")
	// 	}
	// } else {
	// 	if OA != bl.OrderTo {
	// 		return fmt.Errorf("The given OrderAt Party is not the current Owner and not authorised to issue an endorsement. -- current OrderAt: ", bl.OrderAt)
	// 	}
	// }

	if bl.BlTransferable == false && bl.OrderCheckbox == true {
		return fmt.Errorf("The given B/L has already been returned to the carrier - a transfer is not possible -- Carrier:  ", bl.CarrierName)
	}

	// OrderTo should not be set "" after first endorsement -> causes problem with following test for first endorsement
	if OT == "" {
		return fmt.Errorf("NewOrderTo darf nicht leer sein")
	}

	//check first endorsement with Ordercheckbox
	if bl.OrderTo == "" {

		if bl.OrderCheckbox == false {

			return fmt.Errorf("The given B/L has no ordering option selected and is therefore not negotiable")
		}
		if bl.ConsigneeName != "" && OA != bl.ConsigneeName {
			return fmt.Errorf("case or order - orderAt must be Consignee -- given OA: ", OA, " Consignee :", bl.ConsigneeName)
		}
		if bl.ConsigneeName == "" && OA != bl.ShipperName {
			return fmt.Errorf("case to order - orderAt must be Shipper -- given OA ", OA, " Shipper: ", bl.ShipperName)
		}
	} else {
		if OA != bl.OrderTo {
			return fmt.Errorf("The given OrderAt party is not the current owner and not authorized to issue an endorsement -- current OrderAt: ", bl.OrderAt)
		}
	}

	if bl.ShippedOnBoardDate == "" {
		return fmt.Errorf("The freight has not been shipped on board, ShippedOnBoardDate not set")
	}

	if bl.ConsigneeName != "" && bl.OrderCheckbox == false {
		bl.BlTransferable = false
	}

	if bl.BlTransferable != false {
		if NPCN != "" || NPCA != "" || NPCLF != "" {
			bl.NotifyPartyCompanyName = NPCN
			bl.NotifyPartyCompanyAddress = NPCA
			bl.NotifyPartyCompanyLegalForm = NPCLF
		}

		currentDate := time.Now()
		formattedDate := currentDate.Format("02.01.2006 15:04:05") // formatted in "DD.MM.YYYY hh:mm:ss"

		bl.OrderDate = formattedDate
		bl.OrderTo = OT // nächsterInhaber
		bl.OrderAt = OA // ÜberschreibendePartei
	} else {
		blAsBytes, _ := json.Marshal(bl)
		ctx.GetStub().PutState(blNumber, blAsBytes)
		return fmt.Errorf("given B/L is not negotiable")
	}

	//last endorsement -> orderto is carrier
	if bl.BlTransferable == true && bl.OrderTo == bl.CarrierName && bl.OrderCheckbox == true {
		bl.BlTransferable = false
	}

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)
}

// DepreciationOfBl function for executing a depreciation --> manipulation of freight-fields in B/L
func (s *SmartContract) DepreciationOfBl(ctx contractapi.TransactionContextInterface, blNumber string, newNrPKG string, newGrossWeight string, newGrossWeightUnit string, newDOG string, newDPP string, newMeasurement string, newMeasurementUnit string, newDCVA string, newDCVC string, newAI string, newHM string) error {

	bl, err := s.QueryBl(ctx, blNumber)

	if err != nil {
		return err
	}

	//Convert arguments to correct datatype
	//manipulate bl freight data fields according to which fields are filled in in frontend
	if newNrPKG != "" {
		convNewNrPKG, converr := strconv.Atoi(newNrPKG)
		if converr != nil {
			return fmt.Errorf("Error while converting newNumberOfPackages to int in depreciation", converr.Error())
		}
		bl.NumberOfPackages = convNewNrPKG
	}

	if newGrossWeight != "" {
		convNewGrossWeight, converr := strconv.Atoi(newGrossWeight)
		if converr != nil {
			return fmt.Errorf("Error while converting newGrossWeight to int in depreciation", converr.Error())
		}
		bl.GrossWeight = convNewGrossWeight
	}

	if newGrossWeightUnit != "" {
		bl.GrossWeightUnit = newGrossWeightUnit
	}

	if newDOG != "" {
		bl.DescriptionOfGoods = newDOG
	}

	if newDPP != "" {
		bl.DescriptionPerPackage = newDPP
	}

	if newMeasurement != "" {
		convNewMeasurement, converr := strconv.ParseFloat(newMeasurement, 64)
		if converr != nil {
			return fmt.Errorf("Error while converting newMeasurement to int in depreciation", converr.Error())
		}
		bl.Measurement = convNewMeasurement
	}

	if newMeasurementUnit != "" {
		bl.MeasurementUnit = newMeasurementUnit
	}

	if newDCVA != "" {
		convDCVA, converr := strconv.Atoi(newDCVA)
		if converr != nil {
			return fmt.Errorf("Error while converting newDeclaredCargoValueAmount (newDCVA) to int in depreciation", converr.Error())
		}
		bl.DeclaredCargoValueAmount = convDCVA
	}

	if newDCVC != "" {
		bl.DeclaredCargoValueCurrency = newDCVC
	}

	if newAI != "" {
		bl.AdditionalInformation = newAI
	}
	if newHM != "" {
		convNewHM, converr := strconv.ParseBool(newHM)
		if converr != nil {
			return fmt.Errorf("Error while converting newHazardousMaterial (newHM) to int in depreciation", converr.Error())
		}
		bl.HazardousMaterial = convNewHM
	}

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create modified fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting modified fabcar chaincode: %s", err.Error())
	}
}
