package blockchain

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

type ArrangementForm struct {
	HospitalKey string `json:"hospitalKey"`
	DoctorKey   string `json:"doctorKey"`
	VisitUnix   int64  `json:"visitUnix"`
}

func (setup *FabricSetup) Arrangement(f *ArrangementForm) (string, error) {

	eventID := "arrangement"

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
		Fcn:          "arrangement",
		Args:         [][]byte{[]byte(f.HospitalKey), []byte(f.DoctorKey), []byte(fmt.Sprintf("%d", f.VisitUnix))},
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
