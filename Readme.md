挂号就诊流程
---

# 整体流程如下:
    医院排班 -> 生成排班记录(ArrangementHistory)
       ↓
    患者挂号 -> 生成挂号单(RegisterHistory[state=Register])
       ↓
    患者就诊 -> 更新挂号单状态和就诊时间(RegisterHistory[state=Visiting,visitUnix=time.Now().Unix()])
       ↓
    医生开处方 -> 生成病例(Case), 处方(Prescription), 订单(Order[State=NotPaid])
       ↓
    患者支付 -> 更新订单状态(Order[State=Paid]), 生成支付记录(PaymentHistory)
       ↓
    患者取药 -> 更新订单状态(Order[State=Finished]), 生成药品出库记录(OutboundHistory)

# 接口访问流程如下：

## 0. query get `/query`

    参数 query_string 是一个 json 对象的序列化字符串

    示例：

    ?query_string={"selector":{"docType":{"$eq":"**Supplier**"}}}

    其中 **Supplier** 是实体的名称，具体实体名称对应如下：
    * 医院      Hospital
    * 医生      Doctor
    * 排班记录  ArrangementHistory
    * 用户      User
    * 挂号单    RegisterHistory
    * 处方      Prescription
    * 订单      Order
    * 病例      Case
    * 支付记录   PaymentHistory
    * 出库记录   OutboundHistory


## 1. 生成数据 get `/initData`

    生成各种实体测试数据, 所有的实体数据可以通过 `/query` 接口查询

## 2. 生成排班记录 post `/arrangement`

    参数 hospitalKey (需要访问 `/query` 获取医院数据[Hospital]，拿到 hospitalKey)

    参数 doctorKey (需要访问 `/query` 获取医生数据[Doctor]，拿到 doctorKey)

    参数 visitUnix 就诊时间，1527897600 10 位数字精度到秒，获取方式 time.Now().Unix(), 一般按 上午 10 点，下午 2 点整的时间算，标记上午还是下午

    示例:
```
    {
    	"userKey": "User-14021119870124337X",
    	"registerHistoryKey": "RegisterHistory-59df2650-4948-41d4-be88-a09c198dd722",
    	"state": "Visiting"
    }
```

## 3. 挂号 get `/createRegister`

    参数 arrangementKey (需要访问 `/query` 获取排班记录数据，拿到 arrangementKey)

    参数 userKey (需要访问 `/query` 获取排班记录数据，拿到 userKey)

    通过 query 可以查看生成的挂号单 RegisterHistory


## 4. 就诊 post `/updateRegister`

    参数 userKey 步骤 3 得到的 userKey

    参数 registerHistoryKey (需要访问 `/query` 获取挂号单数据，拿到 registerHistoryKey)

    参数 state 必须为 Visiting

```
    {
        "userKey": "User-14021119870124337X",
        "registerHistoryKey": "RegisterHistory-59df2650-4948-41d4-be88-a09c198dd722",
        "state": "Visiting"
    }
```

## 5. 开处方 post `/updateRegister`

    参数 userKey 步骤 3 得到的 userKey

    参数 registerHistoryKey 步骤 4 得到的 registerHistoryKey

    参数 state 必须为 PendingPayment

    其他参数如下示例，8 个参数一个不能少

```
    {
        "userKey": "User-14021119870124337X",
        "registerHistoryKey": "RegisterHistory-59df2650-4948-41d4-be88-a09c198dd722",
        "state": "PendingPayment",
        "complained":"发热、恶寒、咳嗽2天，右胸掣痛半天。",
        "diagnose":"因外出衣着不慎而始感头痛，连及巅顶，鼻塞声重，时流清涕，微有咳嗽，恶寒发热，无汗。自以为是“感冒”而服“去痛片”未效，但仍坚持工作。次日病情加重，头痛连及项背，周身酸楚无力，下午3时，突然发热、寒战，咳嗽顿作，痰粘而黄，涕浊，不欲饮食，便秘溲黄，遂到×院急诊。",
        "history":"平素身体尚可，未患过肺结核及肺炎，未患过肝炎，去年查肝功无异常；1987年患过“急性胃肠炎”，经治而愈；无心脏、肾脏、血液、内分泌及神经系统疾病，亦无外伤史。",
        "familyHistory":"母亲年过七旬，尚健。父因“脑出血”于1980年去世。",
        "items":[
            {
                "medicalItemKey": "MedicalItem-国药准字Z10983104-16110111-6933968000031",
                "count": "2"
            },
            {
                "medicalItemKey": "MedicalItem-国药准字Z20054729-170601-6934883300435",
                "count": "1"
            }
            ]
    }
```
    通过 query 可以查看生成的订单 Order

    通过 query 可以查看生成的病例 Case

    通过 query 可以查看生成的处方 Prescription

## 6. 支付 post `/updateRegister`

    参数 userKey 步骤 3 得到的 userKey

    参数 registerHistoryKey 步骤 4 得到的 registerHistoryKey

    参数 state 必须为 Paid

    示例：
    ```
        {
        	"userKey": "User-14021119870124337X",
        	"registerHistoryKey": "RegisterHistory-59df2650-4948-41d4-be88-a09c198dd722",
        	"state": "Paid"
        }
    ```
    通过 query 可以查看生成的支付记录 PaymentHistory

## 7. 取药 post `/updateRegister`

    参数 userKey 步骤 3 得到的 userKey

    参数 registerHistoryKey 步骤 4 得到的 registerHistoryKey

    参数 state 必须为 Finished

    示例:
    ```
        {
        	"userKey": "User-14021119870124337X",
        	"registerHistoryKey": "RegisterHistory-59df2650-4948-41d4-be88-a09c198dd722",
        	"state": "Finished"
        }
    ```
    通过 query 可以查看生成的出库记录 OutboundHistory