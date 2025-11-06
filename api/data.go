package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SalesData represents a sales record
type SalesData struct {
	Date        string  `json:"date"`
	ProductName string  `json:"productName"`
	Amount      float64 `json:"amount"`
	Quantity    int     `json:"quantity"`
}

// InventoryItem represents an inventory item
type InventoryItem struct {
	ID          int    `json:"id"`
	ProductName string `json:"productName"`
	Stock       int    `json:"stock"`
	Status      string `json:"status"`
}

// GetSalesJSON returns sales data as JSON
func GetSalesJSON(w http.ResponseWriter, r *http.Request) {
	sales := []SalesData{
		{Date: time.Now().AddDate(0, 0, -2).Format("2006-01-02"), ProductName: "Laptop", Amount: 1299.99, Quantity: 2},
		{Date: time.Now().AddDate(0, 0, -1).Format("2006-01-02"), ProductName: "Mouse", Amount: 29.99, Quantity: 5},
		{Date: time.Now().Format("2006-01-02"), ProductName: "Keyboard", Amount: 79.99, Quantity: 3},
		{Date: time.Now().Format("2006-01-02"), ProductName: "Monitor", Amount: 349.99, Quantity: 1},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

// GetInventoryJSON returns inventory data as JSON
func GetInventoryJSON(w http.ResponseWriter, r *http.Request) {
	inventory := []InventoryItem{
		{ID: 1, ProductName: "Laptop", Stock: 15, Status: "In Stock"},
		{ID: 2, ProductName: "Mouse", Stock: 45, Status: "In Stock"},
		{ID: 3, ProductName: "Keyboard", Stock: 8, Status: "Low Stock"},
		{ID: 4, ProductName: "Monitor", Stock: 0, Status: "Out of Stock"},
		{ID: 5, ProductName: "Webcam", Stock: 22, Status: "In Stock"},
		{ID: 6, ProductName: "Headphones", Stock: 5, Status: "Low Stock"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

// GetData returns a simple message for the index page
func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<div class="text-center">
		<div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-green-400 to-green-600 rounded-full mb-4">
			<svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
			</svg>
		</div>
		<h3 class="text-lg font-semibold text-gray-800 mb-2">Data Loaded Successfully!</h3>
		<p class="text-gray-600 mb-4">Server time: %s</p>
		<div class="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-4 text-left">
			<p class="text-sm text-gray-700"><strong>Status:</strong> API is working correctly</p>
			<p class="text-sm text-gray-700"><strong>Response:</strong> 200 OK</p>
		</div>
	</div>`, time.Now().Format("2006-01-02 15:04:05"))
}
