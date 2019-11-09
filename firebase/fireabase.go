package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const serverKey = "key=AAAAvb76I88:APA91bEyyBLLrRzx-JlMBDZX0o275a_8nd4XC3OuM9Oiz_spnVEi8PuvugDEVmCO_FDpB1b2zywyxLLOyHFaudKvcXzGVLA6DEvPhN-GFca3nCjIEsPXFkK9JV_x4sEMWN1mxZV1r2Lp"
const url = "https://fcm.googleapis.com/fcm/send"

// Message struct
type Message struct {
	To   string      `json:"to"`
	Data MessageData `json:"data"`
}

// MessageData struct
type MessageData struct {
	Message MessageDataPayload `json:"message"`
}

// MessageDataPayload message
type MessageDataPayload struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// SendMessage send message expire qrcode
func SendMessage(message Message) {
	jsonStr, _ := json.Marshal(message)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", serverKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
