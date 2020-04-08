package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)

type SmartContract struct {
}

type DemandCommit struct {
	Type string `json:"type"`//0
	MissionId string `json:"missionId"`
	DemandId string `json:"demandId"`
	DemandCommitterId string `json:"demandCommitterId"`
	DemandCommitterName string `json:"demandCommitterName"`
	DemandDocHash string `json:"demandDocHash"`
	DemandDocName string `json:"demandDocName"`
	TestSoftwareName string `json:"testSoftwareName"`
	UpdateTime string `json:"updateTime"`
}

	type DemandCheck struct {
		Type string `json:"type"`//1
		MissionId string `json:"missionId"`
		DemandId string `json:"demandId"`
		DemandCheckerId string `json:"demandCheckerId"`
		DemandCheckerName string `json:"demandCheckerName"`
		CheckResult string `json:"checkResult"`
		UpdateTime string `json:"updateTime"`
	}

	type TestReport struct {
		Type string `json:"type"`//2
		MissionId string `json:"missionId"`
		TestReportId string `json:"testReportId"`
		TestReportName string `json:"testReportName"`
		BugReportList string `json:"bugReportList"`
		TestWorkerId string `json:"testWorkerId"`
		TestWorkerName string json:"testWorkerName"`
		UpdateTime string `json:"updateTime"`
	}

	type ReportCheck struct {
		Type string `json:"type"`//3
		MissionId string `json:"missionId"`
		BugReportId string `json:"bugReportId"`
		ReportScore string `json:"bugReportScore"`
		ReportCheckerId string `json:"reportCheckerId"`
		ReportCheckerName string `json:"reportCheckerName"`
		UpdateTime string `json:"updateTime"`
	}

	type ReportMix struct {
		Type string `json:"type"`//4
		MissionId string `json:"missionId"`
		BugReportList string `json:"bugReportList"`
		ReportMixerId string `json:"reportMixerId"`
		ReportMixerName string `json:"reportMixerName"`
		UpdateTime string `json:"updateTime"`
	}

	type MissionState struct {
		Type string `json:"type"`//5
		MissionId string `json:"missionId"`
		MissionName string `json:"missionName"`
		MissionState string `json:"missionState"`
		UpdateTime string `json:"updateTime"`
	}



func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// err := cid.AssertAttributeValue(APIstub, "software", "*")
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	function, args := APIstub.GetFunctionAndParameters()
	if function == "findOne" {
		return s.findOne(APIstub, args)
	} else if function == "save" {
		return s.save(APIstub, args)
	} else if function == "delete" {
		return s.delete(APIstub, args)
	} else if function == "query" {
		return s.query(APIstub, args)
	} else if function == "queryWithPagination" {
		return s.queryWithPagination(APIstub, args)
	} else if function == "history" {
		return s.history(APIstub, args)
	} 

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) findOne(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(docAsBytes)
}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryWithPagination(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	queryString := args[0]
	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	bookmark := args[2]

	queryResults, err := getQueryResultWithPagination(stub, queryString, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) save(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	doctype := args[1]

	var docAsBytes []byte
	if doctype == "0" {
		doc := DemandCommit{Type:doctype, MissionId: args[2], DemandId: args[3], DemandCommitterId: args[4], DemandCommitterName: args[5],DemandDocHash: args[6], DemandDocName: args[7], TestSoftwareName: args[8], UpdateTime: args[9]}
		docAsBytes, _ = json.Marshal(doc)
	} else if doctype == "1"{
		doc := DemandCheck{Type:doctype, MissionId: args[2], DemandId: args[3], DemandCheckerId: args[4], DemandCheckerName: args[5], CheckResult: args[6], UpdateTime: args[7]}
		docAsBytes, _ = json.Marshal(doc)
	} else if doctype == "2"{
		doc := TestReport{Type:doctype, MissionId: args[2], TestReportId: args[3], TestReportName: args[4], BugReportList: args[5], TestWorkerId: args[6], UpdateTime: args[7]}
		docAsBytes, _ = json.Marshal(doc)
	} else if doctype == "3"{
		doc := ReportCheck{Type:doctype, MissionId: args[2], BugReportId: args[3], ReportScore: args[4], ReportCheckerId: args[5], ReportCheckerName: args[6], UpdateTime: args[7]}
		docAsBytes, _ = json.Marshal(doc)
	} else if doctype == "4"{
		doc := ReportMix{Type:doctype, MissionId: args[2], BugReportList: args[3], ReportMixerId:args[4],ReportMixerName:args[5], UpdateTime: args[6]}
		docAsBytes, _ = json.Marshal(doc)
	} else if doctype == "5"{
		doc := MissionState{Type:doctype, MissionId: args[2], MissionName: args[3], MissionState:args[4], UpdateTime: args[5]}
		docAsBytes, _ = json.Marshal(doc)
	}

	APIstub.PutState(args[0], docAsBytes)

	return shim.Success(nil)
}

func (t *SmartContract) delete(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (s *SmartContract) history(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]
	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		buffer.WriteString(string(response.Value))
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"txId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"doc\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"isDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator,false)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func getQueryResultWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {

	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator,true)
	if err != nil {
		return nil, err
	}
	addPaginationMetadataToQueryResults(buffer, responseMetadata)

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface, forPage bool) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	if(forPage){
		buffer.WriteString("{\"resultList\":")
	}
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"doc\":")
		// Record is a JSON object, so we write as-is

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	if(forPage){
		buffer.WriteString(",")
	}

	return &buffer, nil
}

func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *sc.QueryResponseMetadata) *bytes.Buffer {
	buffer.WriteString("\"responseMetadata\":{\"recordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}")

	return buffer
}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
