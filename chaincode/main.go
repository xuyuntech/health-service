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
	Sid string `json:"sid"`
	Name string `json:"name"`
	Birthday string `json:"birthday"`
	Gender string `json:"gender"`
	Age string `json:"age"`
	TotalSpend string `json:"total_spend"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Store Init ->>")
	return shim.Success(nil)
}

func (s *SmartContract) initUser(APIstub shim.ChaincodeStubInterface) sc.Response {
	store := []User{
		User{"14021119870124337X", "蔡悦", "1987-01-24", "1", "31", "18218.00", "北京朝阳区朝阳北路", "18618441311", "330785652@qq.com"},
	}

	i := 0
	for i < len(store) {
		storeAsBytes, _ := json.Marshal(store[i])
		APIstub.PutState(fmt.Sprintf("user-%s", store[i].Sid), storeAsBytes)
		fmt.Println("Added", store[i])
		i = i + 1
	}

	// Notify listeners that an event "initHospital" have been executed (check line 19 in the file invoke.go)
	err = APIstub.SetEvent("initUser", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (s *SmartContract) queryAllUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
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
	err = APIstub.SetEvent("initHospital", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	fmt.Println("function ->>", function)
	if function == "queryAllHospital" {
		return s.queryAllHospital(APIstub, args)
	} else if function == "initHospital" {
		return s.initHospital(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
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
