package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	url := os.Getenv("JIRA_URL")
	email := os.Getenv("JIRA_USER")
	token := os.Getenv("JIRA_TOKEN")

	if email == "" || token == "" || url == "" {
		fmt.Println("❌ JIRA_USER, JIRA_URL ou JIRA_TOKEN não definidos")
		os.Exit(1)
	}

	jsonBody := []byte(`{
    "fields": {
        "project": {
            "id": "10000"
        },
        "parent": {
            "id": null
        },
        "summary": "Teste create",
        "issuetype": {
            "id": "10006"
        },
        "description": "Created automatically via API.",
        "assignee": {
            "accountId": "712020:6078b92a-adfa-4c62-9b4f-ac6d4f9467a6"
			}
		}
	}`)

	maxRetries := 3
	timeout := 10 * time.Second
	client := &http.Client{Timeout: timeout}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("🚀 Tentativa %d de %d...\n", attempt, maxRetries)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("❌ Erro ao criar requisição:", err)
			os.Exit(1)
		}

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(email, token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Falha na tentativa %d: %s\n", attempt, err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("🔁 Status: %s\n", resp.Status)
		fmt.Println("📩 Corpo da resposta:")
		fmt.Println(string(body))

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Println("✅ Sucesso ao criar issue no Jira.")
			break
		} else {
			fmt.Printf("⚠️ Tentativa %d falhou com status %s\n", attempt, resp.Status)
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Println("🏁 Fim da execução.")
}
