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

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a B/L
type SmartContract struct {
	contractapi.Contract
}

//BillOfLading struct contains all BL-Data
type BillOfLading struct {

	// BL Info Fields
	BLNumber         string `json:"BLNumber"`
	DateOfIssue      string `json:"DateOfIssue"`
	PlaceOfIssue     string `json:"PlaceOfIssue"`
	NumberOfBLIssued int    `json:"NumberOfBLIssued"`

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
	NotifyPartySameAs           bool   `json:"NotifyPartySameAs"`

	//Term of sale information
	Incoterms string `json:"Incoterms"` //25

	//Basic freight information
	FreightChargesCurrency string `json:"FreightChargesCurrency"`
	Prepaid                bool   `json:"Prepaid"`
	Collect                bool   `json:"Collect"`

	//Transportinfo
	PortOfLoading         string `json:"PortOfLoading"`
	PortOfDischarge       string `json:"PortOfDischarge"` //30
	PlaceOfReceipt        string `json:"PlaceOfReceipt"`
	PlaceOfDelivery       string `json:"PlaceOfDelivery"`
	OceanVesselName       string `json:"OceanVesselName"`
	ContainerNumber       string `json:"ContainerNumber"`
	FullContainerLoad     bool   `json:"FullContainerLoad"` //35
	LessThenContainerLoad bool   `json:"LessThenContainerLoad"`
	DateofReceived        string `json:"DateofReceived"`
	ShippedOnBoardDate    string `json:"ShippedOnBoardDate"`

	//Gross info
	MarksAndNumbers            string  `json:"MarksAndNumbers"` //40
	NumberOfPackages           int     `json:"NumberOfPackages"`
	GrossWeight                int     `json:"GrossWeight"`
	GrossWeightUnit            string  `json:"GrossWeightUnit"`
	DescriptionOfGoods         string  `json:"DescriptionOfGoods"`
	DescriptionPerPackage      string  `json:"DescriptionPerPackage"` //45
	Measurement                float64 `json:"Measurement"`
	MeasurementUnit            string  `json:"MeasurementUnit"`
	DeclaredCargoValueAmount   int     `json:"DeclaredCargoValueAmount"`
	DeclaredCargoValueCurrency string  `json:"DeclaredCargoValueCurrency"`
	AdditionalInformation      string  `json:"AdditionalInformation"` //50
	HazardousMaterial          bool    `json:"HazardousMaterial"`

	// CustomerOrderNumber
	CustomerOrderNumber int `json:"CustomerOrderNumber"`

	//Used Conditions (ERA600, Art. 20 a)
	TransportConditions string `json:"TransportConditions"`
	ApplieableLaw       string `json:"ApplieableLaw"`
	PlaceOfJurisdiction string `json:"PlaceOfJurisdiction"` //55

	//Endorsement info
	CurrentOwner string `json:"CurrentOwner"`
	OrderBy      string `json:"OrderBy"`
	OrderTo      string `json:"OrderTo"`
	OrderAt      string `json:"OrderAt"` //59
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *BillOfLading
}

// InitLedger adds a base set of BL to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	bls := []BillOfLading{
		BillOfLading{BLNumber: "TW ECON 1000", DateOfIssue: "10.11.2020", PlaceOfIssue: "Bremerhaven", NumberOfBLIssued: 2,
			ShipperName: "Autohaus Main GmbH", ShipperAddress: "Hanauerlandstr. 3460314 Frankfurt, Germany", ShipperContact: "ahmain@beispiel.de", ShipperLegalForm: "GmbH",
			ConsigneeName: "German-Cars Ldt.", ConsigneeAddress: "Fue Avenue, A1 518108 Shanghai, China", ConsigneeContact: "86282452253", ConsigneeLegalForm: "Ldt.",
			CarrierName: "MSC Germany S.A. & Co. KG", CarrierAddress: "Hafenstraße 55, 282127 Bremen, Germany", CarrierContact: "deu-bremen@msc.de", CarrierLegalForm: "S.A. & Co. KG", CarrierTrailerNumber: "HB-KK-596",
			AgentCompanyName: "BLG AutoTerminal Bremerhaven GmbH & Co. KG", AgentCompanyLegalForm: "GmbH & Co. KG", AgentCompanyAddress: "Senator-Borttscheller-Str. 1, 27568 Bremerhaven, Germany",
			NotifyPartyCompanyName: "German-Cars Ldt.", NotifyPartyCompanyAddress: "Fue Avenue, A1 518108 Shanghai, China", NotifyPartyCompanyLegalForm: "Ldt.", NotifyPartySameAs: true,
			Incoterms:              "FOB (2020)",
			FreightChargesCurrency: "USD", Prepaid: true, Collect: true,
			PortOfLoading: "Bremerhaven Containerterminal", PortOfDischarge: "Shanghai Yangshan", PlaceOfReceipt: "Frankfurt am Main, Adresse, Germany", PlaceOfDelivery: "Shanghai, Adresse, China", OceanVesselName: "MSC Gulsun", ContainerNumber: "OOLU1548378", FullContainerLoad: true, LessThenContainerLoad: false, DateofReceived: "08.02.2020", ShippedOnBoardDate: "09.02.2020",
			MarksAndNumbers: "40' steel Dry Cargo Container No CSQU3054383", NumberOfPackages: 15, GrossWeight: 4250, GrossWeightUnit: "Kg", DescriptionOfGoods: "engines and fitting engine parts packaged together on pallets", DescriptionPerPackage: "abc", Measurement: 40.2, MeasurementUnit: "Feet", DeclaredCargoValueAmount: 75000, DeclaredCargoValueCurrency: "USD", AdditionalInformation: "-", HazardousMaterial: false,
			CustomerOrderNumber: 1,
			TransportConditions: "TransportCond", ApplieableLaw: "ApplLaw", PlaceOfJurisdiction: "POJ",
			CurrentOwner: "HSBHV",
			OrderBy:      "",
			OrderTo:      "",
			OrderAt:      "",
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

	convNumberBlIssued, converr := strconv.Atoi(splitResult[3])
	if converr != nil {
		return fmt.Errorf("Error while converting aNumberOfBLIssued to int", converr.Error())
	}

	convNumberOfPackages, converr := strconv.Atoi(splitResult[39])
	if converr != nil {
		return fmt.Errorf("Error while converting aNumberOfPackages to int", converr.Error())
	}

	convGrossWeight, converr := strconv.Atoi(splitResult[40])
	if converr != nil {
		return fmt.Errorf("Error while converting aGrossWeight to int", converr.Error())
	}

	convNotifyPartySameAs, converr := strconv.ParseBool(splitResult[23])
	if converr != nil {
		return fmt.Errorf("Error while converting aNotifyPartySameAs to bool", converr.Error())
	}

	convPrepaid, converr := strconv.ParseBool(splitResult[26])
	if converr != nil {
		return fmt.Errorf("Error while converting aPrepaid to bool", converr.Error())
	}

	convCollect, converr := strconv.ParseBool(splitResult[27])
	if converr != nil {
		return fmt.Errorf("Error while converting aCollect to bool", converr.Error())
	}

	convFullContainerLoad, converr := strconv.ParseBool(splitResult[34])
	if converr != nil {
		return fmt.Errorf("Error while converting aFullContainerLoad to bool", converr.Error())
	}

	convLessThenContainerLoad, converr := strconv.ParseBool(splitResult[35])
	if converr != nil {
		return fmt.Errorf("Error while converting aLessThenContainerLoad to bool", converr.Error())
	}

	convMeasurement, converr := strconv.ParseFloat(splitResult[44], 64)
	if converr != nil {
		return fmt.Errorf("Error while converting aMeasurement to float", converr.Error())
	}

	convDeclaredCargoValueAmount, converr := strconv.Atoi(splitResult[46])
	if converr != nil {
		return fmt.Errorf("Error while converting aDeclaredCargoValueAmount to int", converr.Error())
	}

	convHazardousMaterial, converr := strconv.ParseBool(splitResult[49])
	if converr != nil {
		return fmt.Errorf("Error while converting aHazardousMaterial to bool", converr.Error())
	}

	convCustomerOrderNumber, converr := strconv.Atoi(splitResult[50])
	if converr != nil {
		return fmt.Errorf("Error while converting aCustomerOrderNumber to int", converr.Error())
	}

	// creation of new bl with all fields being initialization of the passed arguments as field data after some fields data type has been converted
	bl := BillOfLading{
		BLNumber:                    splitResult[0],
		DateOfIssue:                 splitResult[1],
		PlaceOfIssue:                splitResult[2],
		NumberOfBLIssued:            convNumberBlIssued, //with converted datatype atoi (string to int)
		ShipperName:                 splitResult[4],
		ShipperAddress:              splitResult[5],
		ShipperContact:              splitResult[6],
		ShipperLegalForm:            splitResult[7],
		ConsigneeName:               splitResult[8],
		ConsigneeAddress:            splitResult[9],
		ConsigneeContact:            splitResult[10],
		ConsigneeLegalForm:          splitResult[11],
		CarrierName:                 splitResult[12],
		CarrierAddress:              splitResult[13],
		CarrierContact:              splitResult[14],
		CarrierLegalForm:            splitResult[15],
		CarrierTrailerNumber:        splitResult[16],
		AgentCompanyName:            splitResult[17],
		AgentCompanyLegalForm:       splitResult[18],
		AgentCompanyAddress:         splitResult[19],
		NotifyPartyCompanyName:      splitResult[20],
		NotifyPartyCompanyAddress:   splitResult[21],
		NotifyPartyCompanyLegalForm: splitResult[22],
		NotifyPartySameAs:           convNotifyPartySameAs, // with converted datatype parseBool (string to bool)
		Incoterms:                   splitResult[24],
		FreightChargesCurrency:      splitResult[25],
		Prepaid:                     convPrepaid, // with converted datatype ParseBool (string to bool)
		Collect:                     convCollect, // with converted datatype ParseBool (string to bool)
		PortOfLoading:               splitResult[28],
		PortOfDischarge:             splitResult[29],
		PlaceOfReceipt:              splitResult[30],
		PlaceOfDelivery:             splitResult[31],
		OceanVesselName:             splitResult[32],
		ContainerNumber:             splitResult[33],
		FullContainerLoad:           convFullContainerLoad,     // with converted datatype ParseBool (string to bool)
		LessThenContainerLoad:       convLessThenContainerLoad, // with converted datatype ParseBool (string to bool)
		DateofReceived:              splitResult[36],
		ShippedOnBoardDate:          splitResult[37],
		MarksAndNumbers:             splitResult[38],
		NumberOfPackages:            convNumberOfPackages, // with converted datatype Atoi (string to int)
		GrossWeight:                 convGrossWeight,      // with converted datatype Atoi (string to int)
		GrossWeightUnit:             splitResult[41],
		DescriptionOfGoods:          splitResult[42],
		DescriptionPerPackage:       splitResult[43],
		Measurement:                 convMeasurement, // with converted datatype ParseFloat (string to float32)
		MeasurementUnit:             splitResult[45],
		DeclaredCargoValueAmount:    convDeclaredCargoValueAmount, // with converted datatype atoi (string to int)
		DeclaredCargoValueCurrency:  splitResult[47],
		AdditionalInformation:       splitResult[48],
		HazardousMaterial:           convHazardousMaterial,   // with converted datatype parseBool (string to bool)
		CustomerOrderNumber:         convCustomerOrderNumber, // with converted datatype atoi (string to int)
		TransportConditions:         splitResult[51],
		ApplieableLaw:               splitResult[52],
		PlaceOfJurisdiction:         splitResult[53],
		CurrentOwner:                splitResult[54],
		OrderBy:                     splitResult[55],
		OrderTo:                     splitResult[56],
		OrderAt:                     splitResult[57],
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

/*

The following block gives a first frame for planned chaincode-functions for interacting with the Bl
--> the current "list of functions" is not completed


func (s *SmartContract) TransferBl() error{

}

func (s *SmartContract) EditGrossInfo() error {

}

func (s *SmartContract) ReturnBl() error{

}

*/

//changeBLOwner could be implemented with the endorsement of the BL

// ChangeBlOwner updates the owner field of car with given id in world state
/*

  Taken out because of missing variable for "owner" in current state of BL-fields

func (s *SmartContract) ChangeBlOwner(ctx contractapi.TransactionContextInterface, blNumber string, newOwner string) error {
	bl, err := s.QueryBl(ctx, blNumber)

	if err != nil {
		return err
	}

	bl.Owner = newOwner

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)
}

*/

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
