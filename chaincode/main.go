package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/satori/go.uuid"
)

type SmartContract struct {
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	fmt.Println("function ->>", function)
	if function == "initData" {
		return s.initData(APIstub, args)
	} else if function == "query" {
		return s.query(APIstub, args)
	} else if function == "createRegister" {
		return s.createRegister(APIstub, args)
	} else if function == "updateRegister" {
		return s.updateRegister(APIstub, args)
	} else if function == "arrangement" {
		return s.arrangement(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) arrangement(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	hospitalKey := args[0]
	doctorKey := args[1]
	visitUnix := args[2]
	objectType := "ArrangementHistory"
	entity := &ArrangementHistory{
		objectType,
		hospitalKey,
		doctorKey,
		visitUnix,
	}
	if err := putState(entity, stub, "arrangement"); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
func (s *SmartContract) initData(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	doctorType := "Doctor"
	userType := "User"
	medicalItemType := "MedicalItem"
	supplierType := "Supplier"
	hospitalType := "Hospital"

	kstzySupplier := Supplier{supplierType, "山东孔圣堂制药有限公司", "山东省邹城经济开发区", "273500", "0537-5300999", "0537-5300999", "http://www.kxtzy.net"}
	xyjSupplier := Supplier{supplierType, "江西新远健药业有限公司", "江西省广昌工业园区新远健大道", "344900", "0795-7378817", "0795-7378817", ""}

	entities := []Entity{
		// Doctors
		&Doctor{doctorType, "60983948578493875", "张秋霞", "副教授,硕士生导师,医学博士", "国家名老中医聂惠民教授学术继承人。擅用经方治疗内科、妇科、儿科等疑难杂病。内科擅长诊治心脑血管疾病（冠心病、高血压、脑缺血、抑郁症）呼吸系统疾病（如哮喘病、慢性支气管炎）脾胃病（萎缩性胃炎、慢性浅表性胃炎、结肠炎、肠易激综合征、）；妇科擅长经血不调带下病（如不孕、月经不调、痛经、带下）；儿科擅长外感发热、咳嗽、消化不良、厌食等。", "", "18689874637", "321234532@qq.com", "0", "40"},
		&Doctor{doctorType, "39804958768948768", "左松青", "副主任医师", "擅长妇科、心脑血管、老年性骨关节病 毕业于首都医科大学，副主任医师从事中医临床工作30余年。具有丰富的临床经验和精湛独到的医术。精于中医内科、妇科、外科常见病、多发病等疑难杂病。尤其在治疗心脑血管病、糖尿病、脾胃病、肝胆病、肿瘤化疗后康复、脑供血不足、压力综合症、骨性关节病、痛经、月经不调、带下病、不孕不育症、胎前产后、更年期综合症及皮肤病、痤疮等有丰富的临床诊疗经验和显著的疗效。", "", "18643434237", "32143243232@qq.com", "0", "55"},
		// Users
		&User{userType, "14021119870124337X", "蔡悦", "1987-01-24", "1", "31", "18218.00", "北京朝阳区朝阳北路", "18618441311", "330785652@qq.com"},
		// Hospitals
		&Hospital{hospitalType, "凤祥园店", "唐山市路北区龙泽路与裕华道交叉口西行50米道南（裕东楼北门）", "0315", "5268016", "130203"},
		&Hospital{hospitalType, "察院街店", "玉田县钰鼎春园小区104号楼", "", "18131566086", "130203"},
		&Hospital{hospitalType, "复兴路店", "唐山市路南区复兴路223号", "0315", "2860826", "130203"},
		&Hospital{hospitalType, "乐亭店", "乐亭永安南路（老大东方南行300米路东）", "0315", "8131897", "130203"},
		// MedicalItems
		&MedicalItem{medicalItemType, "四消丸", "10", "12.30", kstzySupplier.GetKey(), "6933968000031", "16110111", "国药准字Z10983104", "2016-11-30", "2018-11-29"},
		&MedicalItem{medicalItemType, "藿香正气胶囊", "15", "19.80", xyjSupplier.GetKey(), "6934883300435", "170601", "国药准字Z20054729", "2017-06-02", "2019-06-01"},
		// Suppliers
		&kstzySupplier,
		&xyjSupplier,
	}

	for _, vv := range entities {
		key := vv.GetKey()
		if vAsBytes, err := stub.GetState(key); err != nil {
			return shim.Error(fmt.Sprintf("Failed to get %s: %v", vv.GetObjectType(), err))
		} else if vAsBytes != nil {
			fmt.Printf("This %s already exists: %s", vv.GetObjectType(), key)
			return shim.Error(fmt.Sprintf("This %s already exists: %s", vv.GetObjectType(), key))
		}
		vAsBytes, _ := json.Marshal(vv)
		if err := stub.PutState(key, vAsBytes); err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Added", vv)
	}

	// Notify listeners that an event "initHospital" have been executed (check line 19 in the file invoke.go)
	if err := stub.SetEvent("initData", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// createRegister 挂号 -> 生成挂号记录
func (s *SmartContract) createRegister(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	userKey := args[0]
	arrangementKey := args[1]

	entity := &RegisterHistory{
		"RegisterHospitalHistory",
		userKey,
		arrangementKey,
		"Register",
		"",
		fmt.Sprintf("%d", time.Now().Unix())}

	if err := putState(entity, stub, "createRegister"); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// updateRegister 更新挂号记录状态 -> 就诊、已开方待支付、已支付待取药、已取药
func (s *SmartContract) updateRegister(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 参数说明
	// 0 		1					2		[3			4			5			6				7]
	// userKey	registerHistoryKey 	state 	complained 	diagnose 	history 	familyHistory 	itemsStr
	// 用户		挂号记录				状态		主诉			诊断			病史			家族史			药品列表
	userKey := args[0]
	registerHistoryKey := args[1]
	state := args[2]
	if len(args) < 3 {
		return shim.Error("参数至少需要 3 个")
	}
	var (
		registerHistoryBytes []byte
		err                  error
	)

	if registerHistoryBytes, err = stub.GetState(registerHistoryKey); err != nil {
		return shim.Error(err.Error())
	} else if registerHistoryBytes == nil {
		return shim.Error(fmt.Sprintf("该条挂号记录没有找到(%s)", registerHistoryKey))
	}
	var registerHistory RegisterHistory
	if err := json.Unmarshal(registerHistoryBytes, &registerHistory); err != nil {
		return shim.Error(fmt.Sprintf("解析挂号记录出错(%v)", err))
	}
	if registerHistory.UserKey != userKey {
		return shim.Error(fmt.Sprintf("该条挂号记录不属于此用户"))
	}

	caseObjectType := "Case"
	prescriptionObjectType := "Prescription"

	switch state {
	case "Visiting":
		{
			/* ========  就诊 ========= */
			registerHistory.State = "Visiting"
			registerHistory.VisitUnix = fmt.Sprintf("%d", time.Now().Unix())
		}
	case "PendingPayment":
		{
			/* ========  已开处方待支付 ========= */
			if len(args) != 8 {
				return shim.Error("需要 8 个参数")
			}
			complained := args[3]
			diagnose := args[4]
			history := args[5]
			familyHistory := args[6]
			itemsStr := args[7]
			// todo 参数检查
			// 解析药品列表
			var items [][]string
			if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
				return shim.Error("药品列表解析失败")
			}
			// 病例和处方应该是一对一关系
			// 生成新的病例记录
			userCase := &Case{
				caseObjectType,
				uuid.NewV4().String(),
				"",
				complained,
				diagnose,
				history,
				familyHistory,
			}
			// 生成处方, 相当于下订单
			prescription := &Prescription{
				prescriptionObjectType,
				registerHistoryKey,
				uuid.NewV4().String(),
				userCase.GetKey(),
				itemsStr,
			}
			userCase.PrescriptionKey = prescription.GetKey()
			userCaseBytes, _ := json.Marshal(userCase)
			prescriptionBytes, _ := json.Marshal(prescription)
			if err := stub.PutState(userCase.GetKey(), userCaseBytes); err != nil {
				return shim.Error(fmt.Sprintf("保存病例失败: %v", err))
			}
			if err := stub.PutState(prescription.GetKey(), prescriptionBytes); err != nil {
				if err := stub.DelState(userCase.GetKey()); err != nil {
					fmt.Sprintf("保存处方失败时，删除病例(%s)失败", userCase.GetKey())
				}
				return shim.Error(fmt.Sprintf("保存处方失败: %v", err))
			}
			// 更改挂号记录状态，已开处方待支付
			registerHistory.State = "PendingPayment"
		}
	case "Paid":
		/* ========  已支付待取药 ========= */
		registerHistory.State = "Paid"
		// 生成支付记录
	case "Finished":
		/* ========  已取药 ========= */
		registerHistory.State = "Finished"
		// 生成出库记录
	default:
		return shim.Error(fmt.Sprintf("不可用的挂号记录状态(%s)", state))
	}
	if b, err := json.Marshal(&registerHistory); err != nil {
		return shim.Error(fmt.Sprintf("序列化回写记录失败(%v)", err))
	} else if err := stub.PutState(registerHistory.GetKey(), b); err != nil {
		return shim.Error(fmt.Sprintf("回写挂号记录失败(%v)", err))
	}

	// Notify listeners that an event "updateRegister" have been executed (check line 19 in the file invoke.go)
	if err := stub.SetEvent("updateRegister", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {

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

func putState(v interface{}, stub shim.ChaincodeStubInterface, eventName ...string) error {
	entity := v.(Entity)
	key := entity.GetKey()
	if vAsBytes, err := stub.GetState(key); err != nil {
		return errors.New(fmt.Sprintf("Failed to get %s: %v", entity.GetObjectType(), err))
	} else if vAsBytes != nil {
		fmt.Printf("This %s already exists: %s", entity.GetObjectType(), key)
		return errors.New(fmt.Sprintf("This %s already exists: %s", entity.GetObjectType(), key))
	}
	vAsBytes, _ := json.Marshal(entity)
	if err := stub.PutState(key, vAsBytes); err != nil {
		return err
	}
	if len(eventName) > 0 {
		if err := stub.SetEvent(eventName[0], []byte{}); err != nil {
			return err
		}
	}

	return nil
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
