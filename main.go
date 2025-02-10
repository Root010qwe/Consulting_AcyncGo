//go run main.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// URL эндпоинта Django-сервиса для обновления заявки
const MainServiceUpdateURL = "http://localhost:8000/requests/update_result/"

// Токен для псевдо-авторизации (должен совпадать с AUTH_TOKEN в Django)
const AuthToken = "My8Byte"

// Структура для отправки результата (QR удален)
type Result struct {
	PK        int     `json:"pk"`
	TotalCost float64 `json:"total_cost"`
	Token     string  `json:"token"`
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func randomTotalCost() float64 {
	time.Sleep(5 * time.Second)
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Intn(1000) + 100)
}

func SendResult(pk string, url string) {
	totalCost := randomTotalCost()
	var intPK int
	_, err := fmt.Sscanf(pk, "%d", &intPK)
	if err != nil {
		fmt.Println("Error parsing pk:", err)
		return
	}
	data := Result{
		PK:        intPK,
		TotalCost: totalCost,
		Token:     AuthToken,
	}
	resp, err := performPUTRequest(url, data)
	if err != nil {
		fmt.Println("Error sending result:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Result sent successfully for pk:", pk)
	} else {
		fmt.Println("Failed to send result for pk:", pk, "Status:", resp.Status)
	}
}

func main() {
	r := gin.Default()
	r.POST("/set_status", func(c *gin.Context) {
		pk := c.PostForm("pk")
		if pk == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing pk"})
			return
		}
		go SendResult(pk, MainServiceUpdateURL)
		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	r.Run(":8080")
}
