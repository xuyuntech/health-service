package blockchain

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin/json"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

func (setup *FabricSetup) CreateRegister(userKey, arrangementKey string) (string, error) {
	var args []string
	args = append(args, "createRegister")
	args = append(args, userKey)
	args = append(args, arrangementKey)

	eventID := "createRegister"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	// Register a notification handler on the client
	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode evet: %v", err)
	}

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(chclient.Request{
		ChaincodeID:  setup.ChainCodeID,
		Fcn:          args[0],
		Args:         [][]byte{[]byte(args[1]), []byte(args[2])},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	return response.TransactionID.ID, nil

}

type UpdateRegisterForm struct {
	UserKey            string `json:"userKey"`
	RegisterHistoryKey string `json:"registerHistoryKey"`
	State              string `json:"state"`
	Complained         string `json:"complained"`
	Diagnose           string `json:"diagnose"`
	History            string `json:"history"`
	FamilyHistory      string `json:"familyHistory"`
	Items              []struct {
		MedicalItemKey string `json:"medicalItemKey"`
		Count          string `json:"count"`
	} `json:"items"`
}

func (setup *FabricSetup) UpdateRegister(f *UpdateRegisterForm) (string, error) {

	eventID := "updateRegister"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	// Register a notification handler on the client
	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode evet: %v", err)
	}
	// userKey	registerHistoryKey 	state 	complained 	diagnose 	history 	familyHistory 	items
	itemsBytes, _ := json.Marshal(f.Items)
	// Create a request (proposal) and send it
	response, err := setup.client.Execute(chclient.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn:         "updateRegister",
		Args: [][]byte{
			[]byte(f.UserKey), []byte(f.RegisterHistoryKey), []byte(f.State), []byte(f.Complained),
			[]byte(f.Diagnose), []byte(f.History), []byte(f.FamilyHistory), itemsBytes,
		},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	return response.TransactionID.ID, nil
}
