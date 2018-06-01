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

type Entity interface{
	GetKey() string
	GetObjectType() string
}

type Hospital struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone1     string `json:"phone1"`
	Phone2     string `json:"phone2"`
	CountryID  string `json:"country_id"`
}
 
func (h Hospital) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Name)
}
func (h Hospital) GetObjectType() string {
	return h.ObjectType
}

type MedicalItem struct {
	ObjectType string `json:"docType"`
	Title string `json:"title"`
	Quantity string `json:"quantity"`
	Price string `json:"price"`
	SupplierID string `json:"supplierID"` // 生产厂家
	BarCode string `json:"barCode"` // 条码
	BatchNumber string `json:"batchNumber"` //批号
	PermissionNumber string `json:"permissionNumber"` // 批准文号
	ProductionDate string `json:"productionDate"`
	ExpiredDate string `json:"expiredDate"`
}

func (h MedicalItem) GetKey() string {
	return fmt.Sprintf("%s-%s-%s-%s", h.ObjectType, h.PermissionNumber, h.BatchNumber, h.BarCode)
}
func (h MedicalItem) GetObjectType() string {
	return h.ObjectType
}

type Supplier struct {
	ObjectType string  `json:"docType"`
	Name string `json:"name"`
	Address string `json:"address"`
	ZipCode string `json:"zipCode"`
	Telephone string `json:"Telephone"`
	Fax string `json:"fax"`
	WebSite string `json:"webSite"`
}

func (h Supplier) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Name)
}
func (h Supplier) GetObjectType() string {
	return h.ObjectType
}

type Doctor struct {
	ObjectType  string `json:"docType"`
	Sid         string `json:"sid"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Age         string `json:"age"`
}

func (h Doctor) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Sid)
}
func (h Doctor) GetObjectType() string {
	return h.ObjectType
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

func (h User) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Sid)
}
func (h User) GetObjectType() string {
	return h.ObjectType
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

func (s *SmartContract) initData(stub shim.ChaincodeStubInterface) sc.Response {
	doctorType := "Doctor"
	userType := "User"
	medicalItemType := "MedicalItem"
	supplierType := "Supplier"
	hospitalType := "Hospital"

	doctors := []Doctor{
		Doctor{doctorType, "60983948578493875", "张秋霞", "副教授,硕士生导师,医学博士", "国家名老中医聂惠民教授学术继承人。擅用经方治疗内科、妇科、儿科等疑难杂病。内科擅长诊治心脑血管疾病（冠心病、高血压、脑缺血、抑郁症）呼吸系统疾病（如哮喘病、慢性支气管炎）脾胃病（萎缩性胃炎、慢性浅表性胃炎、结肠炎、肠易激综合征、）；妇科擅长经血不调带下病（如不孕、月经不调、痛经、带下）；儿科擅长外感发热、咳嗽、消化不良、厌食等。", "", "18689874637", "321234532@qq.com", "0", "40"},
		Doctor{doctorType, "39804958768948768", "左松青", "副主任医师", "擅长妇科、心脑血管、老年性骨关节病 毕业于首都医科大学，副主任医师从事中医临床工作30余年。具有丰富的临床经验和精湛独到的医术。精于中医内科、妇科、外科常见病、多发病等疑难杂病。尤其在治疗心脑血管病、糖尿病、脾胃病、肝胆病、肿瘤化疗后康复、脑供血不足、压力综合症、骨性关节病、痛经、月经不调、带下病、不孕不育症、胎前产后、更年期综合症及皮肤病、痤疮等有丰富的临床诊疗经验和显著的疗效。", "" "18643434237", "32143243232@qq.com", "0", "55"},
	}

	users := []User{
		User{userType, "14021119870124337X", "蔡悦", "1987-01-24", "1", "31", "18218.00", "北京朝阳区朝阳北路", "18618441311", "330785652@qq.com"},
	}

	hospitals := []Hospital{
		Hospital{hospitalType, "凤祥园店", "唐山市路北区龙泽路与裕华道交叉口西行50米道南（裕东楼北门）", "0315", "5268016", "130203"},
		Hospital{hospitalType, "察院街店", "玉田县钰鼎春园小区104号楼", "", "18131566086", "130203"},
		Hospital{hospitalType, "复兴路店", "唐山市路南区复兴路223号", "0315", "2860826", "130203"},
		Hospital{hospitalType, "乐亭店", "乐亭永安南路（老大东方南行300米路东）", "0315", "8131897", "130203"},
	}

	medicalItems := []MedicalItem{}
	suppliers := []MedicalItem{}
	
	if err := putEntitiesToState(doctors); err != nil {
		return err
	}
	if err := putEntitiesToState(users); err != nil {
		return err
	}
	if err := putEntitiesToState(hospitals); err != nil {
		return err
	}
	if err := putEntitiesToState(medicalItems); err != nil {
		return err
	}
	if err := putEntitiesToState(suppliers); err != nil {
		return err
	}

	// Notify listeners that an event "initHospital" have been executed (check line 19 in the file invoke.go)
	if err := stub.SetEvent("initData", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func putEntitiesToState(entities []Entity) error {
	for _, v ：= range entities {
		key := v.GetKey()
		if vAsBytes, err := stub.GetState(key); err != nil {
			return shim.Error(fmt.Sprintf("Failed to get %s: %v", v.GetObjectType(), err))
		} else if vAsBytes != nil {
			fmt.Printf("This %s already exists: %s", v.GetObjectType(), key)
			return shim.Error(fmt.Sprintf("This %s already exists: %s", v.GetObjectType(), key))
		}
		vAsBytes, _ := json.Marshal(v)
		if err := stub.PutState(key, vAsBytes); err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Added", v)
	}
	return nil
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


/*


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


func (s *SmartContract) initDoctor(stub shim.ChaincodeStubInterface) sc.Response {
	objectType := "Doctor"
	doctors := []Doctor{
		Doctor{objectType, "60983948578493875", "张秋霞", "副教授,硕士生导师,医学博士", "国家名老中医聂惠民教授学术继承人。擅用经方治疗内科、妇科、儿科等疑难杂病。内科擅长诊治心脑血管疾病（冠心病、高血压、脑缺血、抑郁症）呼吸系统疾病（如哮喘病、慢性支气管炎）脾胃病（萎缩性胃炎、慢性浅表性胃炎、结肠炎、肠易激综合征、）；妇科擅长经血不调带下病（如不孕、月经不调、痛经、带下）；儿科擅长外感发热、咳嗽、消化不良、厌食等。", "", "18689874637", "321234532@qq.com", "0", "40"},
		Doctor{objectType, "39804958768948768", "左松青", "副主任医师", "擅长妇科、心脑血管、老年性骨关节病 毕业于首都医科大学，副主任医师从事中医临床工作30余年。具有丰富的临床经验和精湛独到的医术。精于中医内科、妇科、外科常见病、多发病等疑难杂病。尤其在治疗心脑血管病、糖尿病、脾胃病、肝胆病、肿瘤化疗后康复、脑供血不足、压力综合症、骨性关节病、痛经、月经不调、带下病、不孕不育症、胎前产后、更年期综合症及皮肤病、痤疮等有丰富的临床诊疗经验和显著的疗效。", "" "18643434237", "32143243232@qq.com", "0", "55"},
	}

	i := 0
	for i < len(doctors) {
		doctor := doctors[i]
		doctorKey := fmt.Sprintf("doctor-%s", doctor.Sid)
		if userAsBytes, err := stub.GetState(doctorKey); err != nil {
			return shim.Error("Failed to get user: " + err.Error())
		} else if userAsBytes != nil {
			fmt.Println("This user already exists: " + doctorKey)
			return shim.Error("This user already exists: " + doctorKey)
		}
		storeAsBytes, _ := json.Marshal(doctor)
		if err := stub.PutState(doctorKey, storeAsBytes); err != nil {
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

func (s *SmartContract) initUser(stub shim.ChaincodeStubInterface) sc.Response {
	objectType := "User"
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

func (s *SmartContract) initHospital(APIstub shim.ChaincodeStubInterface) sc.Response {
	objectType := "Hospital"
	store := []Hospital{
		Hospital{objectType, "凤祥园店", "唐山市路北区龙泽路与裕华道交叉口西行50米道南（裕东楼北门）", "0315", "5268016", "130203"},
		Hospital{objectType, "察院街店", "玉田县钰鼎春园小区104号楼", "", "18131566086", "130203"},
		Hospital{objectType, "复兴路店", "唐山市路南区复兴路223号", "0315", "2860826", "130203"},
		Hospital{objectType, "乐亭店", "乐亭永安南路（老大东方南行300米路东）", "0315", "8131897", "130203"},
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


*/