package download

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

// FetchSheetAsJSON reads a CSV from the provided URL and returns JSON
func FetchSheetAsJSON(sheetURL string) ([]map[string]string, error) {
	// Step 1: Make the HTTP GET request
	resp, err := http.Get(sheetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	// Step 2: Read the CSV response
	reader := csv.NewReader(resp.Body)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %v", err)
	}

	// Step 3: Convert CSV to JSON
	var jsonData []map[string]string
	headers := rows[0]
	for _, row := range rows[1:] {
		obj := make(map[string]string)
		for i, header := range headers {
			obj[header] = row[i]
		}
		jsonData = append(jsonData, obj)
	}

	return jsonData, nil
}

func getSheetURL(path, gid string) string {
	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?gid=%s&format=csv", path, gid)
}

type JsonResponse struct {
	Data []map[string]string `json:"data"`
}

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	
	path := r.URL.Query().Get("path")
	gid := r.URL.Query().Get("gid")

	if (path == "" || gid == "") {
		http.Error(w, "Bad Request: 'path' and 'gid' query parameter is required", http.StatusBadRequest)
		return
	}

	sheetURL := getSheetURL(path, gid)
	
	// Fetch and convert the sheet data
	jsonData, err := FetchSheetAsJSON(sheetURL)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, fmt.Sprintf("Internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	resp := JsonResponse{
		Data: jsonData,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
