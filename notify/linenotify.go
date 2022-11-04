package notify

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	myKey         = "EHLJWHuCMw1tTHVs2N1zecBrbHoGEp7yF3z04xjDVic"
	lineNotifuUrl = "https://notify-api.line.me/api/notify"
)

func SendMsg(msg string) error {
	req, err := http.NewRequest("POST", lineNotifuUrl, strings.NewReader(fmt.Sprintf("message=%s", msg)))
	if err != nil {
		log.Println("http.NewRequest fail")
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+myKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do fail")
		return err
	}

	resp.Body.Close()
	return nil
}
