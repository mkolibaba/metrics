package router_test

import (
	"context"
	"fmt"
	"github.com/mkolibaba/metrics/internal/common/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"go.uber.org/zap"
)

// ExampleNew демонстрирует создание, запуск и использование сервера метрик.
func ExampleNew() {
	logger := zap.NewNop().Sugar()
	store := inmemory.NewMemStorage()

	// Создание роутера и запуск тестового сервера.
	testRouter := router.New(store, nil, "", nil, logger, rsa.NopDecryptor)
	server := httptest.NewServer(testRouter)
	defer server.Close()

	// Обновление gauge-метрики через URL.
	updateURL := fmt.Sprintf("%s/update/gauge/TestGauge/123.45", server.URL)
	req, _ := http.NewRequest(http.MethodPost, updateURL, nil)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()
	fmt.Printf("Update via URL status: %d\n", resp.StatusCode)

	// Обновление counter-метрики через JSON.
	updateJSONURL := fmt.Sprintf("%s/update/", server.URL)
	jsonBody := `{"id":"TestCounter","type":"counter","delta":10}`
	req, _ = http.NewRequest(http.MethodPost, updateJSONURL, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()
	fmt.Printf("Update via JSON status: %d\n", resp.StatusCode)

	// Получение значения gauge метрики.
	getURL := fmt.Sprintf("%s/value/gauge/TestGauge", server.URL)
	resp, _ = http.Get(getURL)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("Get value status: %d, body: %s\n", resp.StatusCode, string(body))

	// Получение всех метрик на главной странице.
	resp, _ = http.Get(server.URL)
	resp.Body.Close()
	fmt.Printf("List all metrics status: %d\n", resp.StatusCode)

	// Output:
	// Update via URL status: 200
	// Update via JSON status: 200
	// Get value status: 200, body: 123.45
	// List all metrics status: 200
}

// ExampleNew_updateJSONBatch демонстрирует батчевое обновление метрик.
func ExampleNew_updateJSONBatch() {
	logger := zap.NewNop().Sugar()
	store := inmemory.NewMemStorage()
	testRouter := router.New(store, nil, "", nil, logger, rsa.NopDecryptor)
	server := httptest.NewServer(testRouter)
	defer server.Close()

	// Формируем JSON-массив с метриками для обновления.
	batchJSON := `[
		{"id":"BatchGauge","type":"gauge","value":987.6},
		{"id":"BatchCounter","type":"counter","delta":55}
	]`

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/updates/", strings.NewReader(batchJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()

	fmt.Printf("Batch update status: %d\n", resp.StatusCode)

	// Проверка обновления метрики.
	val, err := store.GetGauge(context.Background(), "BatchGauge")
	if err != nil {
		fmt.Println("Error getting gauge")
	}

	fmt.Printf("BatchGauge value after update: %.1f\n", val)

	// Output:
	// Batch update status: 200
	// BatchGauge value after update: 987.6
}
