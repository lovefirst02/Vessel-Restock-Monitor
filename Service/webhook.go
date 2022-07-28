package Service

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func SendWebHook(title, url, value, image string) {
	webhook_url := viper.GetString("WEBHOOK")
	json_str := []byte(fmt.Sprintf(`{"content": null,"embeds": [{"title": "%s","url": "%s","color": null,"fields": [{"name": "size","value": "%s"}],"thumbnail": {"url": "%s"}}],"attachments": []}`, title, url, value, image))
	req, err := http.NewRequest("POST", webhook_url, bytes.NewBuffer(json_str))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {

		fmt.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode >= http.StatusOK {
		fmt.Println(time.Now().Format("[2006-01-02 15:04:05]"), "Send Webhook Sucess")
		return
	}
}
