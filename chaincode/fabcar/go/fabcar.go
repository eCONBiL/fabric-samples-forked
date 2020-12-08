/*
SPDX-License-Identifier: Apache-2.0
*/

/*
modified by Malte Garmhausen

eCONBiL - electronical Consignment and Bill of Lading

changes made for an implementation of the Bill of Lading freight document

This chaincode has 43 of the collected BL fields implemented and provides functions for initializing the ledger, querying specific B/Ls, querying
all B/Ls from the ledger as well as creating a new B/L.
The initialization-function creates one default B/L into the ledger.

Chaincode querys and invokes must be issued through the "interface program" found in /fabric-samples/fabcar/go/fabcar.go which is currently used as
command line interface

*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
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
	ShipperName      string `json:"ShipperName"`
	ShipperAddress   string `json:"ShipperAddress"`
	ShipperContact   string `json:"ShipperContact"`
	ShipperLegalForm string `json:"ShipperLegalForm"`

	//Consignee information
	ConsigneeName      string `json:"ConsigneeName"`
	ConsigneeAddress   string `json:"ConsigneeAddress"`
	ConsigneeContact   string `json:"ConsigneeContact"`
	ConsigneeLegalForm string `json:"ConsigneeLegalForm"`

	//Carrier information
	CarrierName          string `json:"CarrierName"`
	CarrierAddress       string `json:"CarrierAddress"`
	CarrierContact       string `json:"CarrierContact"`
	CarrierLegalForm     string `json:"CarrierLegalForm"`
	CarrierTrailerNumber string `json:"CarrierTrailerNumber"`

	//Forwarding Agent Company information
	AgentCompanyName      string `json:"AgentCompanyName"`
	AgentCompanyLegalForm string `json:"AgentCompanyLegalForm"`
	AgentCompanyAddress   string `json:"AgentCompanyAddress"`

	//Notify Party information
	NotifyPartyCompanyName      string `json:"NotifyPartyCompanyName"`
	NotifyPartyCompanyAddress   string `json:"NotifyPartyCompanyAddress"`
	NotifyPartyCompanyLegalForm string `json:"NotifyPartyCompanyLegalForm"`
	NotifyPartySameAs           bool   `json:"NotifyPartySameAs"`

	//Term of sale information
	Incoterms string `json:"Incoterms"`

	//Basic freight information
	FreightChargesCurrency string `json:"FreightChargesCurrency"`
	Prepaid                bool   `json:"Prepaid"`
	Collect                bool   `json:"Collect"`

	//Transportinfo
	PortOfLoading         string `json:"PortOfLoading"`
	PortOfDischarge       string `json:"PortOfDischarge"`
	PlaceOfReceipt        string `json:"PlaceOfReceipt"`
	PlaceOfDelivery       string `json:"PlaceOfDelivery"`
	OceanVesselName       string `json:"OceanVesselName"`
	ContainerNumber       string `json:"ContainerNumber"`
	FullContainerLoad     bool   `json:"FullContainerLoad"`
	LessThenContainerLoad bool   `json:"LessThenContainerLoad"`
	DateofReceived        string `json:"DateofReceived"`
	ShippedOnBoardDate    string `json:"ShippedOnBoardDate"`

	//Gross info
	MarksAndNumbers    string `json:"MarksAndNumbers"`
	NumberOfPackages   int    `json:"NumberOfPackages"`
	GrossWeight        int    `json:"GrossWeight"`
	GrossWeightUnit    string `json:"GrossWeightUnit"`
	DescriptionOfGoods string `json:"DescriptionOfGoods"`

	/* removed this fields temporarily because of argument limitation issue in create function */
	// 	DescriptionPerPackage      int    `json:"DescriptionPerPackage"`
	// 	Measurement                int    `json:"Measurement"`
	// 	MeasurementUnit            string `json:"MeasurementUnit"`
	// 	DeclaredCargoValueAmount   int    `json:"DeclaredCargoValueAmount"`
	// 	DeclaredCargoValueCurrency string `json:"DeclaredCargoValueCurrency"`
	// 	AdditionalInformation      string `json:"AdditionalInformation"`
	// 	HazardousMaterial          bool `json:"HazardousMaterial"`

	// // Rest
	// 	CustomerOrderNumber int `json:"CustomerOrderNumber"`

	//Used Conditions (ERA600, Art. 20 a)
	// 	TransportConditions string `json:"TransportConditions"`
	// 	ApplieableLaw	string `json:"ApplieableLaw"`
	// 	PlaceOfJurisdiction string `json:"PlaceOfJurisdiction"`

	//Endorsement info
	// 	endorsement information still missing, needs to be implemented later

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
			MarksAndNumbers: "40' steel Dry Cargo Container No CSQU3054383", NumberOfPackages: 15, GrossWeight: 4250, GrossWeightUnit: "Kg", DescriptionOfGoods: "engines and fitting engine parts packaged together on pallets", /*DescriptionPerPackage: 1, Measurement: 40, MeasurementUnit: "Feet", DeclaredCargoValueAmount: 75000, DeclaredCargoValueCurrency: "USD", AdditionalInformation: "-", HazardousMaterial: false, */
			//CustomerOrderNumber: 1,
			//Information about used conditions
			//Endorsement Information

		},
	}

	for i, bl := range bls {
		blAsBytes, _ := json.Marshal(bl)
		err := ctx.GetStub().PutState("TW ECON 100"+strconv.Itoa(i), blAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateBl adds a new BL to the world state with given details
func (s *SmartContract) CreateBl(ctx contractapi.TransactionContextInterface, aBLNumber string, aDateOfIssue string, aPlaceOfIssue string, aNumberOfBLIssued string, aShipperName string, aShipperAddress string, aShipperContact string, aShipperLegalForm string,
	aConsigneeName string, aConsigneeAddress string, aConsigneeContact string, aConsigneeLegalForm string,
	aCarrierName string, aCarrierAddress string, aCarrierContact string, aCarrierLegalForm string, aCarrierTrailerNumber string,
	aAgentCompanyName string, aAgentCompanyLegalForm string, aAgentCompanyAddress string, aNotifyPartyCompanyName string,
	aNotifyPartyCompanyAddress string, aNotifyPartyCompanyLegalForm string, aNotifyPartySameAs string, aIncoterms string,
	aFreightChargesCurrency string, aPrepaid string, aCollect string, aPortOfLoading string, aPortOfDischarge string, aPlaceOfReceipt string, aPlaceOfDelivery string,
	aOceanVesselName string, aContainerNumber string, aFullContainerLoad string, aLessThenContainerLoad string, aDateofReceived string,
	aShippedOnBoardDate string, aMarksAndNumbers string, aNumberOfPackages string, aGrossWeight string, aGrossWeightUnit string, aDescriptionOfGoods string,
	/* aDescriptionPerPackage int, aMeasurement int, aMeasurementUnit string, aDeclaredCargoValueAmount int, aDeclaredCargoValueCurrency string,
	   aAdditionalInformation string, aHazardousMaterial bool, aCustomerOrderNumber int*/) error {

	// converting passed arguments from string to correct data type in the following block

	convNumberBlIssued, converr := strconv.Atoi(aNumberOfBLIssued)
	if converr != nil {
		return fmt.Errorf("Error while converting aNumberOfBLIssued to int", converr.Error())
	}

	convNumberOfPackages, converr := strconv.Atoi(aNumberOfPackages)
	if converr != nil {
		return fmt.Errorf("Error while converting aNumberOfPackages to int", converr.Error())
	}

	convGrossWeight, converr := strconv.Atoi(aGrossWeight)
	if converr != nil {
		return fmt.Errorf("Error while converting aGrossWeight to int", converr.Error())
	}

	convNotifyPartySameAs, converr := strconv.ParseBool(aNotifyPartySameAs)
	if converr != nil {
		return fmt.Errorf("Error while converting aNotifyPartySameAs to bool", converr.Error())
	}

	convPrepaid, converr := strconv.ParseBool(aPrepaid)
	if converr != nil {
		return fmt.Errorf("Error while converting aPrepaid to bool", converr.Error())
	}

	convCollect, converr := strconv.ParseBool(aCollect)
	if converr != nil {
		return fmt.Errorf("Error while converting aCollect to bool", converr.Error())
	}

	convFullContainerLoad, converr := strconv.ParseBool(aFullContainerLoad)
	if converr != nil {
		return fmt.Errorf("Error while converting aFullContainerLoad to bool", converr.Error())
	}

	convLessThenContainerLoad, converr := strconv.ParseBool(aLessThenContainerLoad)
	if converr != nil {
		return fmt.Errorf("Error while converting aLessThenContainerLoad to bool", converr.Error())
	}

	// creation of new bl with all fields being initialization of the passed arguments as field data after some fields data type has been converted
	bl := BillOfLading{
		BLNumber:                    aBLNumber,
		DateOfIssue:                 aDateOfIssue,
		PlaceOfIssue:                aPlaceOfIssue,
		NumberOfBLIssued:            convNumberBlIssued, //with converted datatype atoi (string to int)
		ShipperName:                 aShipperName,       //5
		ShipperAddress:              aShipperAddress,
		ShipperContact:              aShipperContact,
		ShipperLegalForm:            aShipperLegalForm,
		ConsigneeName:               aConsigneeName,
		ConsigneeAddress:            aConsigneeAddress, //10
		ConsigneeContact:            aConsigneeContact,
		ConsigneeLegalForm:          aConsigneeLegalForm,
		CarrierName:                 aCarrierName,
		CarrierAddress:              aCarrierAddress,
		CarrierContact:              aCarrierContact, //15
		CarrierLegalForm:            aCarrierLegalForm,
		CarrierTrailerNumber:        aCarrierTrailerNumber,
		AgentCompanyName:            aAgentCompanyName,
		AgentCompanyLegalForm:       aAgentCompanyLegalForm,
		AgentCompanyAddress:         aAgentCompanyLegalForm, //20
		NotifyPartyCompanyName:      aNotifyPartyCompanyName,
		NotifyPartyCompanyAddress:   aNotifyPartyCompanyAddress,
		NotifyPartyCompanyLegalForm: aNotifyPartyCompanyLegalForm,
		NotifyPartySameAs:           convNotifyPartySameAs, // with converted datatype parseBool (string to bool)
		Incoterms:                   aIncoterms,            //25
		FreightChargesCurrency:      aFreightChargesCurrency,
		Prepaid:                     convPrepaid, // with converted datatype ParseBool (string to bool)
		Collect:                     convCollect, // with converted datatype ParseBool (string to bool)
		PortOfLoading:               aPortOfLoading,
		PortOfDischarge:             aPortOfDischarge, //30
		PlaceOfReceipt:              aPlaceOfReceipt,
		PlaceOfDelivery:             aPlaceOfDelivery,
		OceanVesselName:             aOceanVesselName,
		ContainerNumber:             aContainerNumber,
		FullContainerLoad:           convFullContainerLoad,     //35	 with converted datatype ParseBool (string to bool)
		LessThenContainerLoad:       convLessThenContainerLoad, // with converted datatype ParseBool (string to bool)
		DateofReceived:              aDateofReceived,
		ShippedOnBoardDate:          aShippedOnBoardDate,
		MarksAndNumbers:             aMarksAndNumbers,
		NumberOfPackages:            convNumberOfPackages, //40		with converted datatype Atoi (string to int)
		GrossWeight:                 convGrossWeight,      // with converted datatype Atoi (string to int)
		GrossWeightUnit:             aGrossWeightUnit,
		DescriptionOfGoods:          aDescriptionOfGoods,

		/* removed these fields temporarily because of arugument limitation issue --> 47 args max */
		// DescriptionPerPackage: aDescriptionPerPackage,
		// Measurement: aMeasurement,
		// MeasurementUnit: aMeasurementUnit,
		// DeclaredCargoValueAmount: aDeclaredCargoValueAmount,
		// DeclaredCargoValueCurrency: aDeclaredCargoValueCurrency,	//45
		// AdditionalInformation: aAdditionalInformation,
		// HazardousMaterial: aHazardousMaterial,
		// CustomerOrderNumber: aCustomerOrderNumber,
	}

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(aBLNumber, blAsBytes)
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
