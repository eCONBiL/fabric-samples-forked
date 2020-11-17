/*
SPDX-License-Identifier: Apache-2.0
*/

/*
modified by Malte Garmhausen
changes made for an implementation of the Bill of Lading freight document
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


//BL Felder
// Info describes basic infos from BL
// type Info struct {
//     BLNumber         string `json:"BLNumber"`
//     DateOfIssue      string `json:"DateOfIssue"`
//     PlaceOfIssue     string `json:"PlaceOfIssue"`
//     NumberOfBLIssued int    `json:"NumberOfBLIssued"`
// }
 
// Shipper describes basic infos from Shipper
// type Shipper struct {
//     ShipperName           string `json:"ShipperName"`
//     ShipperAddress        string `json:"ShipperAddress"`
//     ShipperContact        string `json:"ShipperContact"`
//     ShipperLegalForm      string `json:"ShipperLegalForm"`
//     ShipperSignature      string `json:"ShipperSignature"`
//     ShippperSignatureDate string `json:"ShippperSignatureDate"`
//}
/* 
// Consignee describes basic infos from Consignee
type Consignee struct {
    ConsigneeName      string `json:"ConsigneeName"`
    ConsigneeAddress   string `json:"ConsigneeAddress"`
    ConsigneeContact   string `json:"ConsigneeContact"`
    ConsigneeLegalForm string `json:"ConsigneeLegalForm"`
    ConsigneePublicKey string `json:"ConsigneePublicKey"`
}
 
// Carrier describes basic infos from Carrier
type Carrier struct {
    CarrierName          string `json:"CarrierName"`
    CarrierAddress       string `json:"CarrierAddress"`
    CarrierContact       string `json:"CarrierContact"`
    CarrierLegalForm     string `json:"CarrierLegalForm"`
    CarrierPublicKey     string `json:"CarrierPublicKey"`
    CarrierTrailerNumber string `json:"CarrierTrailerNumber"`
    CarrierSignature     string `json:"CarrierSignature"`
    CarrierSignatureDate string `json:"CarrierSignatureDate"`
}
 
// Forwarding Agent Company describes basic infos from Forwarding Agent Company
type ForwardingAgentCompany struct {
    AgentCompanyName      string `json:"AgentCompanyName"`
    AgentCompanyLegalForm string `json:"AgentCompanyLegalForm"`
    AgentCompanyAddress   string `json:"AgentCompanyAddress"`
    AgentCompanyPublicKey string `json:"AgentCompanyPublicKey"`
}
 
// Notify Party describes basic infos from Notify Party
type NotifyParty struct {
    NotifyPartyCompanyName      string `json:"NotifyPartyCompanyName"`
    NotifyPartyCompanyAddress   string `json:"NotifyPartyCompanyAddress"`
    NotifyPartyCompanyLegalForm string `json:"NotifyPartyCompanyLegalForm"`
    NotifyPartyCompanyPublicKey string `json:"NotifyPartyCompanyPublicKey"`
    NotifyPartySameAs           bool   `json:"NotifyPartySameAs"`
}
 
// Term of Sale describes basic infos from Term of Sale
type TermOfSale struct {
    Incoterms string `json:"Incoterms"`
}
 
// Freight describes basic infos from Freight
type Freight struct {
    FreightChargesCurrency string `json:"FreightChargesCurrency"`
    Prepaid                bool   `json:"Prepaid"`
    Collect                bool   `json:"Collect"`
}
 
// TransportInfo describes basic infos from Transport
type TransportInfo struct {
    PortOfLoading         string `json:"PortOfLoading"`
    PortOfDischarge       string `json:"PortOfDischarge"`
    PlaceOfReceipt        string `json:"PlaceOfReceipt"`
    PlaceOfDelivery       string `json:"PlaceOfDelivery"`
    OceanVesselName       string `json:"OceanVesselName"`
    ContainerNumber       int    `json:"ContainerNumber"`
    FullContainerLoad     bool   `json:"FullContainerLoad"`
    LessThenContainerLoad bool   `json:"LessThenContainerLoad"`
    DateofReceived        string `json:"DateofReceived"`
    ShippedOnBoardDate    string `json:"ShippedOnBoardDate"`
}
 
// GrossInfo describes basic infos from Gross
type GrossInfo struct {
    MarksAndNumbers            string `json:"MarksAndNumbers"`
    NumberOfPackages           int    `json:"NumberOfPackages"`
    GrossWeight                int    `json:"GrossWeight"`
    GrossWeightUnit            string `json:"GrossWeightUnit"`
    DescriptionOfGoods         string `json:"DescriptionOfGoods"`
    DescriptionPerPackage      int    `json:"DescriptionPerPackage"`
    Measurement                int    `json:"Measurement"`
    MeasurementUnit            string `json:"MeasurementUnit"`
    DeclaredCargoValueAmount   int    `json:"DeclaredCargoValueAmount"`
    DeclaredCargoValueCurrency string `json:"DeclaredCargoValueCurrency"`
    AdditionalInformation      string `json:"AdditionalInformation"`
    HazardousMaterial          string `json:"HazardousMaterial"`
}
 
// Rest describes basic infos from Rest
type Rest struct {
    CustomerOrderNumber int `json:"CustomerOrderNumber"`
}
// Endorsement describes basic infos from Endorsement
type Endorsement struct {
    
}
*/

type BillOfLading struct {

// BL Info Fields
	BLNumber         string `json:"BLNumber"`
    	DateOfIssue      string `json:"DateOfIssue"`
    	PlaceOfIssue     string `json:"PlaceOfIssue"`
	NumberOfBLIssued int    `json:"NumberOfBLIssued"`
	
//Shipper information	
	ShipperName           string `json:"ShipperName"`
    	ShipperAddress        string `json:"ShipperAddress"`
    	ShipperContact        string `json:"ShipperContact"`
    	ShipperLegalForm      string `json:"ShipperLegalForm"`
    	ShipperSignature      string `json:"ShipperSignature"`
	ShippperSignatureDate string `json:"ShippperSignatureDate"`

//Consignee information
	ConsigneeName      string `json:"ConsigneeName"`
    	ConsigneeAddress   string `json:"ConsigneeAddress"`
    	ConsigneeContact   string `json:"ConsigneeContact"`
    	ConsigneeLegalForm string `json:"ConsigneeLegalForm"`
	ConsigneePublicKey string `json:"ConsigneePublicKey"`
	
//Forwarding Agent Company information
	AgentCompanyName      string `json:"AgentCompanyName"`
    	AgentCompanyLegalForm string `json:"AgentCompanyLegalForm"`
    	AgentCompanyAddress   string `json:"AgentCompanyAddress"`
	AgentCompanyPublicKey string `json:"AgentCompanyPublicKey"`
	
//Notify Party information	
	NotifyPartyCompanyName      string `json:"NotifyPartyCompanyName"`
   	NotifyPartyCompanyAddress   string `json:"NotifyPartyCompanyAddress"`
    	NotifyPartyCompanyLegalForm string `json:"NotifyPartyCompanyLegalForm"`
    	NotifyPartyCompanyPublicKey string `json:"NotifyPartyCompanyPublicKey"`
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
	ContainerNumber       int    `json:"ContainerNumber"`
	FullContainerLoad     bool   `json:"FullContainerLoad"`
	LessThenContainerLoad bool   `json:"LessThenContainerLoad"`
	DateofReceived        string `json:"DateofReceived"`
	ShippedOnBoardDate    string `json:"ShippedOnBoardDate"`

//Gross info
	MarksAndNumbers            string `json:"MarksAndNumbers"`
	NumberOfPackages           int    `json:"NumberOfPackages"`
	GrossWeight                int    `json:"GrossWeight"`
	GrossWeightUnit            string `json:"GrossWeightUnit"`
	DescriptionOfGoods         string `json:"DescriptionOfGoods"`
	DescriptionPerPackage      int    `json:"DescriptionPerPackage"`
	Measurement                int    `json:"Measurement"`
	MeasurementUnit            string `json:"MeasurementUnit"`
	DeclaredCargoValueAmount   int    `json:"DeclaredCargoValueAmount"`
	DeclaredCargoValueCurrency string `json:"DeclaredCargoValueCurrency"`
	AdditionalInformation      string `json:"AdditionalInformation"`
	HazardousMaterial          string `json:"HazardousMaterial"`

// Rest	
	CustomerOrderNumber int `json:"CustomerOrderNumber"`

/*
// Used Conditions (ERA600, Art. 20 a)
	TransportConditions string `json:"TransportConditions"`
	ApplieableLaw	string `json:"ApplieableLaw"`
	PlaceOfJurisdiction string `json:"PlaceOfJurisdiction"`

//Endorsement info 
	//endorsement information still missing, needs to be implemented later
*/
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *BillOfLading
}

// InitLedger adds a base set of BL to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	bls := []BillOfLading{
		BillOfLading{BLNumber: "TW ECON 1001", DateOfIssue: "10.11.2020", PlaceOfIssue: "Bremerhaven", NumberOfBLIssued: 2,
					ShipperName: "Autohaus Main GmbH", ShipperAddress: "Hanauerlandstr. 3460314 Frankfurt, Germany", ShipperContact: "ahmain@beispiel.de", ShipperLegalForm: "GmbH", ShipperSignature: "Hashwert1", ShippperSignatureDate: "09.02.2020",
					ConsigneeName: "German-Cars Ldt.", ConsigneeAddress: "Fue Avenue, A1 518108 Shanghai, China", ConsigneeContact: "86282452253", ConsigneeLegalForm: "Ldt.", ConsigneePublicKey: "ConsigneePublicKey",
					AgentCompanyName: "BLG AutoTerminal Bremerhaven GmbH & Co. KG", AgentCompanyLegalForm: "GmbH & Co. KG", AgentCompanyAddress: "Senator-Borttscheller-Str. 1, 27568 Bremerhaven, Germany", AgentCompanyPublicKey: "AgentCompanyPublicKey",
					NotifyPartyCompanyName: "German-Cars Ldt.", NotifyPartyCompanyAddress: "Fue Avenue, A1 518108 Shanghai, China", NotifyPartyCompanyLegalForm: "Ldt.", NotifyPartyCompanyPublicKey: "NotifyPartyPublicKey", NotifyPartySameAs: true,
					Incoterms: "FOB (2020)",
					FreightChargesCurrency: "USD", Prepaid: true, Collect: true,
					PortOfLoading: "Bremerhaven Containerterminal", PortOfDischarge: "Shanghai Yangshan", PlaceOfReceipt: "Frankfurt am Main, Adresse, Germany", PlaceOfDelivery: "Shanghai, Adresse, China", OceanVesselName: "MSC Gulsun", ContainerNumber: 3, FullContainerLoad: true, LessThenContainerLoad: false, DateofReceived: "08.02.2020", ShippedOnBoardDate: "09.02.2020",
					MarksAndNumbers: "40' steel Dry Cargo Container No CSQU3054383", NumberOfPackages: 15,GrossWeight: 4250, GrossWeightUnit: "Kg", DescriptionOfGoods: "engines and fitting engine parts packaged together on pallets", DescriptionPerPackage: 1, Measurement: 40, MeasurementUnit: "Feet", DeclaredCargoValueAmount: 75000, DeclaredCargoValueCurrency: "USD", AdditionalInformation: "-", HazardousMaterial: "none", 
					CustomerOrderNumber: 1,
					//Information about used conditions 
					//Endorsement Information


				},		 	
		}
			// Shipper{"Autohaus Main Gmbh", "Hanauerlandstr. 3460314 Frankfurt, Germany", "ahmain@beispiel.de", "GmbH", "Hashwert1", "09.02.2020"},
			// Consignee{"German-Cars Ldt.", "Fue Avenue, A1 518108 Shanghai, China", "86282452253", "Ldt.", "ConsigneePublicKey"},
			// Carrier{"MSC Germany S.A. & Co. KG", "Hafenstraße 55, 282127 Bremen, Germany", "deu-bremen@msc.de", "S.A. & Co. KG", "CarrierPublicKey","HB-KK-596", "Carriersignature", "09.02.2020"},
			// ForwardingAgentCompany{"BLG AutoTerminal Bremerhaven GmbH & Co. KG", "GmbH & Co. KG", "Senator-Borttscheller-Str. 1, 27568 Bremerhaven, Germany", "AgentCompanyPublicKey"},
			// NotifyParty{"German-Cars Ldt.", "Fue Avenue, A1 518108 Shanghai, China", "Ldt.", "NotifyPartyPublicKey", true},
			// TermOfSale{"FOB (2020)"},
			// Freight{"USD", true, true},
			// TransportInfo{"Bremerhaven Containerterminal", "Shanghai Yangshan", "Frankfurt am Main, Adresse, Germany", "Shanghai, Adresse, China", "MSC Gulsun", 3, true, false, "08.02.2020", "09.02.2020"},
			// GrossInfo{"40' steel Dry Cargo Container No CSQU3054383", 15, 4250, "Kg", "engines and fitting engine parts packaged together on pallets", 1, 40, "Feet", 75000, "USD", "-", "none"},
			// Rest{1},},
		

	for i, bl := range bls {
		blAsBytes, _ := json.Marshal(bl)
		err := ctx.GetStub().PutState("TW ECON 1001"+strconv.Itoa(i), blAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateBl adds a new BL to the world state with given details
func (s *SmartContract) CreateBl(ctx contractapi.TransactionContextInterface, blNumber string, dateOfIssue string, dateofissue string, placeofissue string, owner string) error {
	bl := BillOfLading{
		//change arguments for this function and the way the arguments are passed into a new BL struct, should look like this 
		// BLNumber: blNumber,
		// DateOfIssue: dateOfIssue;
		
	}

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)
}

// QueryBl returns the car stored in the world state with given id
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

// QueryAllBls returns all cars found in world state
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
