package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"strconv"

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
	} else if function == "find" {
		return s.find(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/**
 * @api {post} /arrangement 排班
 * @apiName arrangement
 * @apiGroup HealthService
 *
 * @apiParam {String} hospitalKey 医院 key
 * @apiParam {String} doctorKey 医生 key
 * @apiParam {String} visitUnix 出诊时间(time.Now().Unix())
 * @apiUse ArrangementHistory
 */
func (s *SmartContract) arrangement(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	hospitalKey := args[0]
	doctorKey := args[1]
	visitUnix := args[2]
	objectType := "ArrangementHistory"
	visitUnixInt, err := strconv.ParseInt(visitUnix, 10, 64)
	if err != nil {
		return shim.Error(fmt.Sprintf("visitUnix 转化错误:(%v)", err))
	}
	vt := time.Unix(visitUnixInt, 0)
	ampm := "am"
	if vt.Hour() >= 12 {
		ampm = "pm"
	}
	entity := &ArrangementHistory{
		objectType,
		hospitalKey,
		doctorKey,
		visitUnix,
		vt.Format("2006-01-02"),
		ampm,
	}
	if err := putState(entity, stub, "arrangement"); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

/**
 * @api {get} /init 初始化数据
 * @apiDescription 生成 医院、医生、用户、药品、供应商 等基础数据
 * @apiName initData
 * @apiGroup HealthService
 * @apiUse Hospital
 * @apiUse Doctor
 * @apiUse User
 * @apiUse MedicalItem
 * @apiUse Supplier
 */
func (s *SmartContract) initData(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if initDataBytes, err := stub.GetState("initData"); err != nil {
		return shim.Error(fmt.Sprintf("get initData err:(%v)", err))
	} else if fmt.Sprintf("%s", initDataBytes) == "true" {
		return shim.Error("重复初始化")
	}

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
		&User{userType, "18748819870304776X", "张三", "1987-03-04", "1", "31", "18218.00", "北京朝阳区朝阳北路", "18678398789", "3838393778@qq.com"},
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

	if err := stub.PutState("initData", []byte("true")); err != nil {
		return shim.Error(fmt.Sprintf("save initData err:(%v)", err))
	}

	// Notify listeners that an event "initHospital" have been executed (check line 19 in the file invoke.go)
	if err := stub.SetEvent("initData", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

/**
 * @api {get} /createRegister 生成挂号记录
 * @apiName createRegister
 * @apiGroup HealthService
 * @apiUse RegisterHistory
 */
// createRegister 挂号 -> 生成挂号记录
func (s *SmartContract) createRegister(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	userKey := args[0]
	arrangementKey := args[1]

	entity := &RegisterHistory{
		"RegisterHistory",
		userKey,
		arrangementKey,
		"",
		"",
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
		user                 User
		userBytes            []byte
		registerHistoryBytes []byte
		err                  error
	)
	if userBytes, err = stub.GetState(userKey); err != nil {
		return shim.Error(fmt.Sprintf("获取用户失败:(%v)", err))
	}
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return shim.Error(fmt.Sprintf("用户序列化失败:(%v)", err))
	}
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
	/**
	 * @api {get} /updateRegister/visiting 就诊
	 * @apiDescription 用户来医院就诊，挂号状态由 Register -> Visiting
	 * @apiName /updateRegister/visiting
	 * @apiGroup HealthService
	 */
	case "Visiting":
		{
			/* ========  就诊 ========= */
			registerHistory.State = "Visiting"
			registerHistory.VisitUnix = fmt.Sprintf("%d", time.Now().Unix())
			break
		}
	/**
		 * @api {get} /updateRegister/prescription 开处方
		 * @apiDescription 医生开处方，同时生成 处方、订单、病例 数据
		 * @apiName /updateRegister/prescription
		 * @apiGroup HealthService
		 * @apiUse Prescription
		 * @apiUse Order
	 	 * @apiUse Case
	*/
	case "PendingPayment":
		{
			/* ========  已开处方待支付 ========= */
			// 获取排班
			var (
				arrangement      ArrangementHistory
				arrangementBytes []byte
				hospital         Hospital
				hospitalBytes    []byte
				doctor           Doctor
				doctorBytes      []byte
			)
			if arrangementBytes, err = stub.GetState(registerHistory.ArrangementKey); err != nil {
				return shim.Error(fmt.Sprintf("获取排班失败:(%v)", err))
			}
			if err := json.Unmarshal(arrangementBytes, &arrangement); err != nil {
				return shim.Error(fmt.Sprintf("排班序列化失败:(%v)", err))
			}
			if doctorBytes, err = stub.GetState(arrangement.DoctorKey); err != nil {
				return shim.Error(fmt.Sprintf("获取医师失败:(%v)", err))
			}
			if err := json.Unmarshal(doctorBytes, &doctor); err != nil {
				return shim.Error(fmt.Sprintf("医师序列化失败:(%v)", err))
			}
			if hospitalBytes, err = stub.GetState(arrangement.HospitalKey); err != nil {
				return shim.Error(fmt.Sprintf("获取医院失败:(%v)", err))
			}
			if err := json.Unmarshal(hospitalBytes, &hospital); err != nil {
				return shim.Error(fmt.Sprintf("医院序列化失败:(%v)", err))
			}
			// 验证挂号单状态是否为 Visiting
			if registerHistory.State != "Visiting" {
				return shim.Error("挂号单还未核销")
			}
			if len(args) != 8 {
				return shim.Error("需要 8 个参数")
			}
			complained := args[3]
			diagnose := args[4]
			history := args[5]
			familyHistory := args[6]
			itemsStr := args[7]
			// todo 参数检查

			// todo 验证药品列表
			var items []struct {
				MedicalItemKey string `json:"medicalItemKey"`
				Count          string `json:"count"`
			}
			if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
				return shim.Error(fmt.Sprintf("药品列表解析失败(%v)", err))
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
				userKey,
				user.Name,
				arrangement.DoctorKey,
				doctor.Name,
				arrangement.HospitalKey,
				hospital.Name,
				fmt.Sprintf("%d", time.Now().Unix()),
			}
			// 生成处方
			prescription := &Prescription{
				prescriptionObjectType,
				registerHistoryKey,
				uuid.NewV4().String(),
				"",
				itemsStr,
				userKey,
				user.Name,
				arrangement.DoctorKey,
				doctor.Name,
				fmt.Sprintf("%d", time.Now().Unix()),
			}
			userCaseKey := userCase.GetKey()
			prescriptionKey := prescription.GetKey()

			userCase.PrescriptionKey = prescriptionKey
			prescription.CaseKey = userCaseKey

			userCaseBytes, _ := json.Marshal(userCase)
			prescriptionBytes, _ := json.Marshal(prescription)
			if err := stub.PutState(userCaseKey, userCaseBytes); err != nil {
				return shim.Error(fmt.Sprintf("保存病例失败: %v", err))
			}
			if err := stub.PutState(prescriptionKey, prescriptionBytes); err != nil {
				if err := stub.DelState(userCaseKey); err != nil {
					fmt.Sprintf("保存处方失败时，删除病例(%s)失败", userCaseKey)
				}
				return shim.Error(fmt.Sprintf("保存处方失败: %v", err))
			}
			// 生成订单
			// TODO 生成订单明细，明细里应该包括售出药品的全量维度信息
			//    计算总消费额
			var (
				spending         float64
				medicalItem      MedicalItem
				medicalItemBytes []byte
			)
			for _, item := range items {
				if medicalItemBytes, err = stub.GetState(item.MedicalItemKey); err != nil {
					return shim.Error(fmt.Sprintf("获取药品失败 %s (%v)", item.MedicalItemKey, err))
				}
				if err := json.Unmarshal(medicalItemBytes, &medicalItem); err != nil {
					return shim.Error(fmt.Sprintf("药品序列化失败(%v)", err))
				}
				p, err := strconv.ParseFloat(medicalItem.Price, 64)
				if err != nil {
					return shim.Error(fmt.Sprintf("价格解析失败"))
				}
				spending += p
			}
			order := &Order{
				"Order",
				prescriptionKey,
				uuid.NewV4().String(),
				"NotPaid",
				fmt.Sprintf("%d", time.Now().Unix()),
				itemsStr,
				fmt.Sprintf("%.2f", spending),
				registerHistoryKey,
				userKey,
				user.Name,
				userCaseKey,
			}
			orderBytes, _ := json.Marshal(order)
			orderKey := order.GetKey()
			if err := stub.PutState(orderKey, orderBytes); err != nil {
				return shim.Error(fmt.Sprintf("保存订单失败(%v)", err))
			}

			// 更改挂号记录, 添加处方外键和订单外键
			registerHistory.State = "Finished"
			registerHistory.PrescriptionKey = prescriptionKey
			registerHistory.OrderKey = orderKey
			break
		}
	/**
	 * @api {get} /updateRegister/payment 支付
	 * @apiDescription 开处方后，用户支付处方费用，订单状态 NotPaid -> Paid，同时生成支付记录
	 * @apiName updateRegister/payment
	 * @apiGroup HealthService
	 * @apiUse PaymentHistory
	 */
	case "Paid":
		/* ========  支付操作 -> 更改订单状态为 Paid，生成支付记录 ========= */
		// 获取订单
		orderKey := registerHistory.OrderKey
		if orderKey == "" {
			return shim.Error("挂号记录没有绑定对应订单")
		}
		var (
			order      Order
			orderBytes []byte
			err        error
		)
		if orderBytes, err = stub.GetState(orderKey); err != nil {
			return shim.Error(fmt.Sprintf("没有找到对应订单 %s (%v)", orderKey, err))
		}
		if err := json.Unmarshal(orderBytes, &order); err != nil {
			return shim.Error(fmt.Sprintf("订单序列化失败(%v)", err))
		}
		order.State = "Paid" // 更改订单状态为已支付
		orderBytesN, _ := json.Marshal(&order)
		if err := stub.PutState(orderKey, orderBytesN); err != nil {
			return shim.Error(fmt.Sprintf("保存订单失败(%v)", err))
		}
		// 生成支付记录
		paymentHistory := &PaymentHistory{
			"PaymentHistory",
			uuid.NewV4().String(),
			order.Spending,
			orderKey,
			order.PrescriptionKey,
			registerHistoryKey,
			fmt.Sprintf("%d", time.Now().Unix()),
			userKey,
			user.Name,
		}
		paymentHistoryKey := paymentHistory.GetKey()
		paymentHistoryBytes, _ := json.Marshal(paymentHistory)
		if err := stub.PutState(paymentHistoryKey, paymentHistoryBytes); err != nil {
			return shim.Error(fmt.Sprintf("保存支付记录失败(%v)", err))
		}

		break
	/**
	 * @api {get} /updateRegister/finished 取药
	 * @apiDescription 用户取药，完成订单。订单状态 Paid -> Finished，同时生成出库记录
	 * @apiName updateRegister/finished
	 * @apiGroup HealthService
	 * @apiUse OutboundHistory
	 */
	case "Finished":
		/* ========  已取药 -> 修改订单状态为完成，生成出库记录 ========= */
		orderKey := registerHistory.OrderKey
		if orderKey == "" {
			return shim.Error("订单未生成")
		}
		var (
			orderBytes []byte
			err        error
			order      Order
		)
		if orderBytes, err = stub.GetState(orderKey); err != nil {
			return shim.Error(fmt.Sprintf("获取订单失败 %s (%v)", orderKey, err))
		}
		if err := json.Unmarshal(orderBytes, &order); err != nil {
			return shim.Error(fmt.Sprintf("订单序列化失败(%v)", err))
		}
		order.State = "Finished"
		orderBytes, err = json.Marshal(&order)
		if err := stub.PutState(orderKey, orderBytes); err != nil {
			return shim.Error(fmt.Sprintf("更新订单状态失败(%v)", err))
		}

		// 生成出库记录
		outboundHistory := &OutboundHistory{
			"OutboundHistory",
			uuid.NewV4().String(),
			orderKey,
			order.PrescriptionKey,
			registerHistoryKey,
			fmt.Sprintf("%d", time.Now().Unix()),
			order.Items,
		}
		outboundHistoryKey := outboundHistory.GetKey()
		outboundHistoryBytes, _ := json.Marshal(outboundHistory)
		if err := stub.PutState(outboundHistoryKey, outboundHistoryBytes); err != nil {
			return shim.Error(fmt.Sprintf("生成出库记录失败(%v)", err))
		}
		break
	default:
		return shim.Error(fmt.Sprintf("不可用的挂号记录状态(%s)", state))
	}
	if b, err := json.Marshal(&registerHistory); err != nil {
		return shim.Error(fmt.Sprintf("序列化回写记录失败(%v)", err))
	} else if err := stub.PutState(registerHistoryKey, b); err != nil {
		return shim.Error(fmt.Sprintf("回写挂号记录失败(%v)", err))
	}

	// Notify listeners that an event "updateRegister" have been executed (check line 19 in the file invoke.go)
	if err := stub.SetEvent("updateRegister", []byte{}); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) find(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]
	var (
		oBytes []byte
		err    error
	)
	oBytes, err = stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("获取失败: (%v)", err))
	}
	return shim.Success(oBytes)
}

/**
 * @api {get} /query 查询实体列表
 * @apiName query
 * @apiGroup HealthService
 * @apiParam {String} query_string {"selector":{"docType":{"$eq":"实体名称"}}}
 * @apiExample {curl} Hospital
 *    http://localhost:8080/query?query_string={"selector":{"docType":{"$eq":"Hospital"}}}
 */
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
