package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type LineWebhookRequest struct {
	Events []LineEvent `json:"events"`
}

type LineEvent struct {
	Type       string      `json:"type"`
	Message    LineMessage `json:"message"`
	ReplyToken string      `json:"replyToken"`
}

type LineMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LineResponse struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []LineMessage `json:"messages"`
}

func validateSignature(body []byte, signature string) bool {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	if channelSecret == "" {
		log.Println("LINE_CHANNEL_SECRET is not set")
		return false
	}

	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature == expectedSignature
}

func replyMessage(replyToken, text string) error {
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if channelAccessToken == "" {
		return fmt.Errorf("LINE_CHANNEL_ACCESS_TOKEN is not set")
	}

	response := LineResponse{
		ReplyToken: replyToken,
		Messages: []LineMessage{
			{Type: "text", Text: text},
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/reply", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func processMessage(event LineEvent) error {
	// Only process text messages
	if event.Message.Type != "text" {
		return nil
	}

	// Extract the text and create echo response
	text := event.Message.Text
	total, dist, err := calNameNumber(text)
	if err != nil {
		return err
	}

	response := fmt.Sprintf("ชื่อ %s\nผลรวม = %d\n%s", text, total, dist)

	return replyMessage(event.ReplyToken, response)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received", r.Method, r.URL.Path)

	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	// Only accept POST requests for webhook
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate LINE signature
	signature := r.Header.Get("X-Line-Signature")
	if signature == "" || !validateSignature(body, signature) {
		log.Println("Invalid signature")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Parse the webhook request
	var webhookReq LineWebhookRequest
	if err := json.Unmarshal(body, &webhookReq); err != nil {
		log.Printf("Error parsing webhook request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Process each event
	for _, event := range webhookReq.Events {
		if event.Type == "message" {
			log.Printf("Received message: %s", event.Message.Text)
			if err := processMessage(event); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}

	// Return 200 OK to acknowledge receipt
	w.WriteHeader(http.StatusOK)
}

func getAlphabetNumber(a rune) (int, error) {
	// Define Thai character to number mapping
	thaiNumMap := map[rune]int{
		'ก': 1,
		'ด': 1,
		'ท': 1,
		'ถ': 1,
		'ภ': 1,
		'ฤ': 1,
		'ฦ': 1,
		'่': 1,
		'ุ': 1,
		'า': 1,
		'ำ': 1,
		'ข': 2,
		'ช': 2,
		'ง': 2,
		'บ': 2,
		'ป': 2,
		'้': 2,
		'เ': 2,
		'แ': 2,
		'ู': 2,
		'ฆ': 3,
		'ต': 3,
		'ฑ': 3,
		'ฒ': 3,
		'๋': 3,
		'ค': 4,
		'ธ': 4,
		'ญ': 4,
		'ร': 4,
		'ษ': 4,
		'ะ': 4,
		'ิ': 4,
		'โ': 4,
		'ั': 4,
		'ฉ': 5,
		'ณ': 5,
		'ฌ': 5,
		'น': 5,
		'ม': 5,
		'ห': 5,
		'ฎ': 5,
		'ฬ': 5,
		'ฮ': 5,
		'ึ': 5,
		'จ': 6,
		'ล': 6,
		'ว': 6,
		'อ': 6,
		'ใ': 6,
		'ซ': 7,
		'ศ': 7,
		'ส': 7,
		'๊': 7,
		'ี': 7,
		'ื': 7,
		'ผ': 8,
		'ฝ': 8,
		'พ': 8,
		'ฟ': 8,
		'ย': 8,
		'็': 8,
		'ฏ': 9,
		'ฐ': 9,
		'ไ': 9,
		'์': 9,
	}

	number, ok := thaiNumMap[a]
	if !ok {
		return 0, fmt.Errorf("invalid character: %c", a)
	}

	return number, nil
}

func calNameNumber(name string) (int, string, error) {
	// Define Thai character to number mapping

	total := 0
	dist := ""
	for _, a := range name {
		n, err := getAlphabetNumber(a)
		if err != nil {
			return 0, "", err
		}

		total += n
		dist += fmt.Sprintf("%s = %d ", string(a), n)
	}
	return total, dist, nil
}
