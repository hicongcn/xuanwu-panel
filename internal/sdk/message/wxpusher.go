package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WxPusher struct {
	AppToken      string   `json:"appToken"`
	Content       string   `json:"content"`
	ContentType   int      `json:"contentType"`
	TopicIds      []int    `json:"topicIds,omitempty"`
	Uids          []string `json:"uids,omitempty"`
	Url           string   `json:"url,omitempty"`
	VerifyPayType int      `json:"verifyPayType,omitempty"`
}

type wxPusherResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Uid       string `json:"uid"`
		TopicId   int    `json:"topicId"`
		MessageId int    `json:"messageId"`
		Code      int    `json:"code"`
		Status    string `json:"status"`
	} `json:"data"`
}

func (w *WxPusher) Send() (string, error) {
	apiUrl := "https://wxpusher.zjiecode.com/api/send/message"

	body, err := json.Marshal(w)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res wxPusherResponse
	if err := json.Unmarshal(respBody, &res); err != nil {
		return string(respBody), err
	}

	if res.Code == 1000 {
		return string(respBody), nil
	}

	return string(respBody), fmt.Errorf("WxPusher error: %s (code: %d)", res.Msg, res.Code)
}
