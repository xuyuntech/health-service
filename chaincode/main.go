package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Hospital struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone1    string `json:"phone1"`
	Phone2    string `json:"phone2"`
	CountryID string `json:"country_id"`
}

type User struct {
	ObjectType string `json:"docType"`
	Sid        string `json:"sid"`
	Name       string `json:"name"`
	Birthday   string `json:"birthday"`
	Gender     string `json:"gender"`
	Age        string `json:"age"`
	TotalSpend string `json:"totalSpend"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Store Init ->>")
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	fmt.Println("function ->>", function)
	if function == "queryAllHospital" {
		return s.queryAllHospital(APIstub, args)
	} else if function == "initHospital" {
		return s.initHospital(APIstub)
	} else if function == "query" {
		return s.query(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initUser(stub shim.ChaincodeStubInterface) sc.Response {
	objectType := "user"
	users := []User{
		User{objectType, "14021119870124337X", "蔡悦", "1987-01-24", "1", "31", "18218.00", "北京朝阳区朝阳北路", "18618441311", "330785652@qq.com"},
	}

	i := 0
	for i < len(users) {
		user := users[i]
		userKey := fmt.Sprintf("user-%s", user.Sid)
		if userAsBytes, err := stub.GetState(userKey); err != nil {
			return shim.Error("Failed to get user: " + err.Error())
		} else if userAsBytes != nil {
			fmt.Println("This user already exists: " + userKey)
			return shim.Error("This user already exists: " + userKey)
		}
		storeAsBytes, _ := json.Marshal(user)
		if err := stub.PutState(userKey, storeAsBytes); err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Added", user)
		i = i + 1

	}

	// Notify listeners that an event "initHospital" have been executed (check line 19 in the file invoke.go)
	if err := stub.SetEvent("initUser", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (s *SmartContract) queryAllUser(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	startKey := "0"
	endKey := "999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Sprintf("- queryAllStore:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initHospital(APIstub shim.ChaincodeStubInterface) sc.Response {
	store := []Hospital{
		Hospital{"凤祥园店", "唐山市路北区龙泽路与裕华道交叉口西行50米道南（裕东楼北门）", "0315", "5268016", "130203"},
		Hospital{"察院街店", "玉田县钰鼎春园小区104号楼", "", "18131566086", "130203"},
		Hospital{"复兴路店", "唐山市路南区复兴路223号", "0315", "2860826", "130203"},
		Hospital{"乐亭店", "乐亭永安南路（老大东方南行300米路东）", "0315", "8131897", "130203"},
	}

	i := 0
	for i < len(store) {
		fmt.Println("i is ", i)
		storeAsBytes, _ := json.Marshal(store[i])
		APIstub.PutState(strconv.Itoa(i+1), storeAsBytes)
		fmt.Println("Added", store[i])
		i = i + 1
	}

	// Notify listeners that an event "initHospital" have been executed (check line 19 in the file invoke.go)
	if err := APIstub.SetEvent("initHospital", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryAllHospital(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Sprintf("- queryAllStore:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

/*
	* main function *
 calls the Start function
 The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
