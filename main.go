package main

import (
	"fmt"
	"html/template"
	"log"
	"mydashboard/api"
	"mydashboard/auth"
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

func renderLoginTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// Page handlers
func indexPage(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"Title": "Dashboard Home",
		"Page":  "home",
	}
	if claims != nil {
		data["User"] = map[string]interface{}{
			"Username": claims.Username,
			"Email":    claims.Email,
		}
	}
	renderTemplate(w, "templates/index.html", data)
}

func salesPage(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"Title": "Sales Dashboard",
		"Page":  "sales",
	}
	if claims != nil {
		data["User"] = map[string]interface{}{
			"Username": claims.Username,
			"Email":    claims.Email,
		}
	}
	renderTemplate(w, "templates/sales.html", data)
}

func inventoryPage(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.GetUserFromContext(r.Context())
	data := map[string]interface{}{
		"Title": "Inventory Dashboard",
		"Page":  "inventory",
	}
	if claims != nil {
		data["User"] = map[string]interface{}{
			"Username": claims.Username,
			"Email":    claims.Email,
		}
	}
	renderTemplate(w, "templates/inventory.html", data)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is already authenticated via cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie != nil {
		claims, err := auth.ValidateToken(cookie.Value)
		if err == nil && claims != nil {
			// Already authenticated, redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	renderLoginTemplate(w, "templates/login.html", map[string]string{
		"Title": "Login",
	})
}

func main() {
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/login", loginPage).Methods("GET")
	router.HandleFunc("/api/auth/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/api/auth/logout", auth.LogoutHandler).Methods("POST")

	// Protected API routes
	router.HandleFunc("/api/data", auth.AuthMiddleware(http.HandlerFunc(api.GetData)).ServeHTTP).Methods("GET")
	router.HandleFunc("/api/sales/json", auth.AuthMiddleware(http.HandlerFunc(api.GetSalesJSON)).ServeHTTP).Methods("GET")
	router.HandleFunc("/api/inventory/json", auth.AuthMiddleware(http.HandlerFunc(api.GetInventoryJSON)).ServeHTTP).Methods("GET")
	router.HandleFunc("/api/auth/me", auth.AuthMiddleware(http.HandlerFunc(auth.MeHandler)).ServeHTTP).Methods("GET")

	// Protected page routes
	router.HandleFunc("/", auth.AuthMiddleware(http.HandlerFunc(indexPage)).ServeHTTP)
	router.HandleFunc("/sales", auth.AuthMiddleware(http.HandlerFunc(salesPage)).ServeHTTP)
	router.HandleFunc("/inventory", auth.AuthMiddleware(http.HandlerFunc(inventoryPage)).ServeHTTP)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
