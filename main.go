package main

import (
	"fmt"
	"html/template"
	"log"
	"mydashboard/api"
	"net/http"

	"github.com/gorilla/mux"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/layout.html", tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// Page handlers
func indexPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/index.html", map[string]string{
		"Title": "Dashboard Home",
		"Page":  "home",
	})
}

func salesPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/sales.html", map[string]string{
		"Title": "Sales Dashboard",
		"Page":  "sales",
	})
}

func inventoryPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/inventory.html", map[string]string{
		"Title": "Inventory Dashboard",
		"Page":  "inventory",
	})
}

func main() {
	router := mux.NewRouter()

	// Page routes
	router.HandleFunc("/sales", salesPage)
	router.HandleFunc("/inventory", inventoryPage)

	// API endpoints
	router.HandleFunc("/api/data", api.GetData).Methods("GET")
	router.HandleFunc("/api/sales/json", api.GetSalesJSON).Methods("GET")
	router.HandleFunc("/api/inventory/json", api.GetInventoryJSON).Methods("GET")

	// Serve index page at root
	router.HandleFunc("/", indexPage)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
