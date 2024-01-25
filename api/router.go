package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Router() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"https://*", "http://*"}, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/v1/payment-intent", app.getPaymentIntent)

	mux.Get("/api/v1/widget/{id}", app.getWidgetByID)

	mux.Post("/api/v1/create-customer-and-subscribe-to-plan", app.createCustomerAndSubscribeToPlan)

	mux.Post("/api/v1/authenticate", app.createAuthToken)

	return mux
}
