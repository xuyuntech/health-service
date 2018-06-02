package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Entity interface {
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

func (h *Hospital) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, h.Name)
}
func (h *Hospital) GetObjectType() string {
	return h.ObjectType
}

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

// RegisterHospitalHistory 挂号记录
type RegisterHistory struct {
	ObjectType     string `json:"docType"`
	UserKey        string `json:"userKey"`
	ArrangementKey string `json:"arrangementKey"`
	State          string `json:"state"`     // 状态 [挂号(Register),就诊中(Visiting),已开处方待支付(PendingPayment),已支付待取药(Paid),已取药(Finished)]
	VisitUnix      string `json:"visitUnix"` // 就诊时间
	Created        string `json:"Created"`   // unix time
}

func (h RegisterHistory) GetKey() string {
	return fmt.Sprintf("%s-%s", h.ObjectType, uuid.NewV4().String())
}
func (h RegisterHistory) GetObjectType() string {
	return h.ObjectType
}

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
