package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func makeRequest(feature, tag int, lv bool) {
	client := &http.Client{}

	query := fmt.Sprintf("http://127.0.0.1:8080/user_banner?tag_id=%d&feature_id=%d&use_last_revision=%v", tag, feature, lv)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Token", "user_token")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

}

func main() {
	ticker := time.NewTicker(time.Millisecond * 1000 / 1000)
	defer ticker.Stop()

	rand.Seed(time.Now().UnixNano())
	i := 0
	for range ticker.C {
		feature := rand.Intn(1000)
		tag := rand.Intn(1000)
		i++
		go makeRequest(feature, tag, i%2 == 0)
	}
}
