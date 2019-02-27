/**
* Copyright 2018 Cognizant. All Rights Reserved.
*
* EAS IPM Blockchain Solutions
*
*/

package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type LifeInsuranceChaincode struct {

}

type User struct{
	FirstName 		string			`json:"first_name"`
	LastName 		string			`json:"last_name"`
	PAN				string			`json:"pan"`
	DOB          	string   		`json:"dob"`
	AnnualIncome	string			`json:"annual_income"`
	PolicyNumber	[10]string		`json:"policy_number"`
	CompanyName		[10]string		`json:"company_name"`
	InsuredAmount	[10]string		`json:"insured_amount"`
	ClaimStatus		[10]string		`json:"claim_status"`
	Comments		[10]string 		`json:"comments"`
	NOP				string			`json:"nop"`				// Number of Policies
}

func (t *LifeInsuranceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// nothing to do
	return shim.Success(nil)
}

func (t *LifeInsuranceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "create" {
		return t.create(stub, args)
	}

	if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error("Error invoking function")
}

// create a policy
func (t *LifeInsuranceChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	FirstName 		:= args[0]
	LastName  		:= args[1]
	PAN 	  		:= args[2]
	DOB			 	:= args[3]
	AnnualIncome 	:= args[4]
	CompanyName		:= args[5]
	InsuredAmount	:= args[6]
	ClaimStatus		:= args[7]
	Comment			:= args[8]

	var NOP , PNO int
	var Policies	[10] string
	var Companies	[10] string
	var InsuredAmt	[10] string
	var ClaimSt		[10] string
	var Comments	[10] string

	if (PAN == "undefined" || PAN == "" || PAN == "null"){
		return shim.Error("Missing PAN ")
	}
	//size := len(args)
	userData, err	:= stub.GetState(PAN)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get data for " + PAN + "\"}"
		return shim.Error(jsonResp)
	} else if userData == nil {
		fmt.Printf("New User data will be created for " + PAN + "\n")
		NOP = 0
	} else if userData != nil {
		userLedger := User{}
		err := json.Unmarshal(userData, &userLedger)
		if err != nil {
			return shim.Error(err.Error())
		}
	
		fmt.Printf(" User data = %s \n" , userData)

		Policies   = userLedger.PolicyNumber
		Companies  = userLedger.CompanyName
		InsuredAmt = userLedger.InsuredAmount
		ClaimSt	 = userLedger.ClaimStatus
		Comments = userLedger.Comments
		NOP, err = strconv.Atoi(userLedger.NOP)
	}
	if NOP == 10 {
		fmt.Printf(" NO MORE THAN 10 CLAIM REQUESTS ALLOWED \n" )
		jsonResp := "{\"Error\":\"No more than 10 claim requests allowed for " + PAN + "\"}"
		return shim.Error(jsonResp)
	}
	PNO = NOP + 1
	fmt.Printf("before policies = %s \n", Policies)

	if ClaimStatus == "APPROVED" {
		Policies[NOP] = CompanyName + PAN + "-POLICY-000" + strconv.Itoa(PNO)
	} else {
		Policies[NOP] = "--NA--"
	}

	Companies[NOP] 	= CompanyName
	InsuredAmt[NOP] = InsuredAmount
	ClaimSt[NOP]	= ClaimStatus
	Comments[NOP] 	= Comment
	NOP++

	fmt.Printf("after policies = %s \n", Policies)

	// Create User and marshal to JSON
	User := &User{FirstName, LastName, PAN, DOB, AnnualIncome, Policies, Companies, InsuredAmt, ClaimSt, Comments, strconv.Itoa(NOP)}
	UserJSONasBytes, err := json.Marshal(User)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Write to ledger
	err = stub.PutState(PAN, UserJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// search a user
func (t *LifeInsuranceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("PAN required to query user data")
	}

	PAN				:= args[0]
	userData, err	:= stub.GetState(PAN)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get data for " + PAN + "\"}"
		return shim.Error(jsonResp)
	} else if userData == nil {
		jsonResp := "{\"Error\":\"User data does not exist for " + PAN + "\"}"
		return shim.Error(jsonResp)
	}

	userLedger := User{}
	err = json.Unmarshal(userData, &userLedger)
	if err != nil {
		return shim.Error(err.Error())
	}

	var Policies [10]string
	Policies = userLedger.PolicyNumber
	fmt.Printf("Policies = %s \n", Policies)
	return shim.Success(userData)

}


func  main()  {
	err := shim.Start(new(LifeInsuranceChaincode))
	if err != nil {
		fmt.Printf("Error starting LifeInsuranceChaincode : %v \n", err)
	}

}
