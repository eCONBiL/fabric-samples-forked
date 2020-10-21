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

// Car describes basic details of what makes up a car
type BillOfLading struct {
	Shipper   string `json:"make"`
	DateOfIssue  string `json:"dateofissue"`
	PlaceOfIssue string `json:"placeofissue"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *BillOfLading
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	bls := []BillOfLading{
		BillOfLading{Shipper: "Autohaus Main GmbH", DateOfIssue: "01.01.2020", PlaceOfIssue: "Bremerhaven", Owner: "Party1"},
		BillOfLading{Shipper: "Abat AG", DateOfIssue: "29.07.2020", PlaceOfIssue: "Bremen", Owner: "Christian"},
		}

	for i, bl := range bls {
		blAsBytes, _ := json.Marshal(bl)
		err := ctx.GetStub().PutState("BL"+strconv.Itoa(i), blAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateBl(ctx contractapi.TransactionContextInterface, blNumber string, shipper string, dateofissue string, placeofissue string, owner string) error {
	bl := BillOfLading{
		Shipper:   shipper,
		DateOfIssue:  dateofissue,
		PlaceOfIssue: placeofissue,
		Owner:  owner,
	}

	blAsBytes, _ := json.Marshal(bl)

	return ctx.GetStub().PutState(blNumber, blAsBytes)
}

// QueryCar returns the car stored in the world state with given id
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

// QueryAllCars returns all cars found in world state
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

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeBlOwner(ctx contractapi.TransactionContextInterface, blNumber string, newOwner string) error {
	bl, err := s.QueryBl(ctx, blNumber)

	if err != nil {
		return err
	}

	bl.Owner = newOwner

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
