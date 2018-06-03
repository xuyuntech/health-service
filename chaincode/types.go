package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Entity interface {
	GetKey() string
	GetObjectType() string
}

/**
 * @apiDefine Hospital 医院
 * @apiExample {golang} 医院
// 数据来源于有赞数据，是我们的线下门店
type Hospital struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone1     string `json:"phone1"`
	Phone2     string `json:"phone2"`
	CountryID  string `json:"country_id"` 	// 邮编
}
*/
type Hospital struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Phone1     string `json:"phone1"`
	Phone2     string `json:"phone2"`
	CountryID  string `json:"country_id"`
}

func (h *Hospital) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Name)
}
func (h *Hospital) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine MedicalItem 药品
 * @apiExample {golang} 药品
type MedicalItem struct {
	ObjectType       string `json:"docType"`
	Title            string `json:"title"`
	Quantity         string `json:"quantity"`
	Price            string `json:"price"`
	SupplierID       string `json:"supplierID"`       // 生产厂家
	BarCode          string `json:"barCode"`          // 条码
	BatchNumber      string `json:"batchNumber"`      //批号
	PermissionNumber string `json:"permissionNumber"` // 批准文号
	ProductionDate   string `json:"productionDate"`
	ExpiredDate      string `json:"expiredDate"`
}
*/
type MedicalItem struct {
	ObjectType       string `json:"docType"`
	Title            string `json:"title"`
	Quantity         string `json:"quantity"`
	Price            string `json:"price"`
	SupplierID       string `json:"supplierID"`       // 生产厂家
	BarCode          string `json:"barCode"`          // 条码
	BatchNumber      string `json:"batchNumber"`      //批号
	PermissionNumber string `json:"permissionNumber"` // 批准文号
	ProductionDate   string `json:"productionDate"`
	ExpiredDate      string `json:"expiredDate"`
}

func (h MedicalItem) GetKey() string {
	return fmt.Sprintf("%s-%s-%s-%s", h.ObjectType, h.PermissionNumber, h.BatchNumber, h.BarCode)
}
func (h MedicalItem) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine RegisterHistory 挂号单
 * @apiExample {golang} 挂号单
type RegisterHistory struct {
	ObjectType      string `json:"docType"`
	UserKey         string `json:"userKey"`
	ArrangementKey  string `json:"arrangementKey"`	// 排班 key
	PrescriptionKey string `json:"prescriptionKey"`	// 处方 key
	OrderKey        string `json:"orderKey"`		// 订单 key
	State           string `json:"state"`     // 状态 [挂号(Register),就诊中(Visiting),已开处方待支付(PendingPayment),已支付待取药(Paid),已取药(Finished)]
	VisitUnix       string `json:"visitUnix"` // 就诊时间
	Created         string `json:"Created"`   // unix time
}
*/
// RegisterHistory 挂号记录
type RegisterHistory struct {
	ObjectType      string `json:"docType"`
	UserKey         string `json:"userKey"`
	ArrangementKey  string `json:"arrangementKey"`
	PrescriptionKey string `json:"prescriptionKey"`
	OrderKey        string `json:"orderKey"`
	State           string `json:"state"`     // 状态 [挂号(Register),就诊中(Visiting),已开处方待支付(PendingPayment),已支付待取药(Paid),已取药(Finished)]
	VisitUnix       string `json:"visitUnix"` // 就诊时间
	Created         string `json:"Created"`   // unix time
}

func (h RegisterHistory) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, uuid.NewV4().String())
}
func (h RegisterHistory) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine ArrangementHistory 排班记录
 * @apiExample {golang} 排班记录
type ArrangementHistory struct {
	ObjectType  string `json:"docType"`
	HospitalKey string `json:"hospitalKey"` // 医院 key
	DoctorKey   string `json:"doctorKey"`  // 医生 key
	VisitUnix   string `json:"visitUnix"` // 出诊日期时间 unix timestamp
}
*/
// ArrangementHistory 排班
type ArrangementHistory struct {
	ObjectType  string `json:"docType"`
	HospitalKey string `json:"hospitalKey"`
	DoctorKey   string `json:"doctorKey"`
	VisitUnix   string `json:"visitUnix"` // 出诊日期时间 unix timestamp
}

func (h ArrangementHistory) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, uuid.NewV4().String())
}
func (h ArrangementHistory) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine Supplier 供应商
 * @apiExample {golang} 供应商
type Supplier struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`  		// 供应商名称
	Address    string `json:"address"` 		// 地址
	ZipCode    string `json:"zipCode"` 		// 邮编
	Telephone  string `json:"Telephone"` 	// 电话
	Fax        string `json:"fax"` 			// 传真
	WebSite    string `json:"webSite"`		// 网址
}
*/
type Supplier struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	ZipCode    string `json:"zipCode"`
	Telephone  string `json:"Telephone"`
	Fax        string `json:"fax"`
	WebSite    string `json:"webSite"`
}

func (h Supplier) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Name)
}
func (h Supplier) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine Doctor 医生
 * @apiExample {golang} 医生
type Doctor struct {
	ObjectType  string `json:"docType"`
	Sid         string `json:"sid"`			//
	Name        string `json:"name"`		//
	Title       string `json:"title"`		// 职称
	Description string `json:"description"`	// 描述，擅长 。。。
	Avatar      string `json:"avatar"`		// 头像
	Phone       string `json:"phone"`		// 手机号
	Email       string `json:"email"`		// 邮箱
	Gender      string `json:"gender"`		// 性别
	Age         string `json:"age"`
}
*/
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

func (h *Doctor) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Sid)
}
func (h *Doctor) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine User 用户
 * @apiExample {golang} 用户
type User struct {
	ObjectType string `json:"docType"`
	Sid        string `json:"sid"`
	Name       string `json:"name"`
	Birthday   string `json:"birthday"`
	Gender     string `json:"gender"`
	Age        string `json:"age"`
	TotalSpend string `json:"totalSpend"`		// 消费总计
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}
*/
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

/**
 * @apiDefine Prescription 处方
 * @apiExample {golang} 处方
type Prescription struct {
	ObjectType         string `json:"docType"`
	RegisterHistoryKey string `json:"registerHistoryKey"`
	Sid                string `json:"sid"`     // 处方编号
	CaseKey            string `json:"caseKey"` // 病例
	Items              string `json:"items"`   // 处方药品列表 [[ItemKey, Count],[]...]
}
*/
// Prescription 处方，相当于下订单
type Prescription struct {
	ObjectType         string `json:"docType"`
	RegisterHistoryKey string `json:"registerHistoryKey"`
	Sid                string `json:"sid"`     // 处方编号
	CaseKey            string `json:"caseKey"` // 病例
	Items              string `json:"items"`   // 处方药品列表 [[ItemKey, Count],[]...]
}

func (h Prescription) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Sid)
}
func (h Prescription) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine Case 病例
 * @apiExample {golang} 病例
type Case struct {
	ObjectType      string `json:"docType"`
	Sid             string `json:"sid"`             // 病例编号
	PrescriptionKey string `json:"prescriptionKey"` // 处方
	Complained      string `json:"complained"`      // 主诉
	Diagnose        string `json:"diagnose"`        // 临床诊断内容, 现病史
	History         string `json:"history"`         // 既往史
	FamilyHistory   string `json:"familyHistory"`   // 家族史
}
*/
// Case 病例
type Case struct {
	ObjectType      string `json:"docType"`
	Sid             string `json:"sid"`             // 病例编号
	PrescriptionKey string `json:"prescriptionKey"` // 处方
	Complained      string `json:"complained"`      // 主诉
	Diagnose        string `json:"diagnose"`        // 临床诊断内容, 现病史
	History         string `json:"history"`         // 既往史
	FamilyHistory   string `json:"familyHistory"`   // 家族史
}

func (h Case) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Sid)
}
func (h Case) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine Order 订单
 * @apiExample {golang} 订单
type Order struct {
	ObjectType      string `json:"docType"`
	PrescriptionKey string `json:"prescriptionKey"` // 处方 key
	Number          string `json:"number"`			// 序列号
	State           string `json:"state"` // NotPaid Paid Finished
	Created         string `json:"created"`
	Items           string `json:"items"`
	Spending        string `json:"spending"` 		// 总金额
}
*/
// Order 订单
type Order struct {
	ObjectType      string `json:"docType"`
	PrescriptionKey string `json:"prescriptionKey"`
	Number          string `json:"number"`
	State           string `json:"state"` // NotPaid Paid Finished
	Created         string `json:"created"`
	Items           string `json:"items"`
	Spending        string `json:"spending"`
}

func (h Order) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Number)
}
func (h Order) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine PaymentHistory 支付记录
 * @apiExample {golang} 支付记录
type PaymentHistory struct {
	ObjectType         string `json:"docType"`
	Number             string `json:"number"`
	Spending           string `json:"spending"` // 此次总消费
	OrderKey           string `json:"orderKey"`
	PrescriptionKey    string `json:"prescriptionKey"` 		// 处方 key
	RegisterHistoryKey string `json:"registerHistoryKey"`	// 挂号单 key
	Created            string `json:"created"`
}
*/
// PaymentHistory 支付记录
type PaymentHistory struct {
	ObjectType         string `json:"docType"`
	Number             string `json:"number"`
	Spending           string `json:"spending"` // 此次总消费
	OrderKey           string `json:"orderKey"`
	PrescriptionKey    string `json:"prescriptionKey"`
	RegisterHistoryKey string `json:"registerHistoryKey"`
	Created            string `json:"created"`
}

func (h PaymentHistory) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Number)
}
func (h PaymentHistory) GetObjectType() string {
	return h.ObjectType
}

/**
 * @apiDefine OutboundHistory 出库记录
 * @apiExample {golang} 出库记录
type OutboundHistory struct {
	ObjectType         string `json:"docType"`
	Number             string `json:"number"`
	OrderKey           string `json:"number"`
	PrescriptionKey    string `json:"prescriptionKey"`
	RegisterHistoryKey string `json:"registerHistoryKey"`
	Created            string `json:"created"`
	Items              string `json:"items"`
}
*/
// OutboundHistory 出库记录
type OutboundHistory struct {
	ObjectType         string `json:"docType"`
	Number             string `json:"number"`
	OrderKey           string `json:"number"`
	PrescriptionKey    string `json:"prescriptionKey"`
	RegisterHistoryKey string `json:"registerHistoryKey"`
	Created            string `json:"created"`
	Items              string `json:"items"`
}

func (h OutboundHistory) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Number)
}
func (h OutboundHistory) GetObjectType() string {
	return h.ObjectType
}
