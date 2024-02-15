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
	mux.Post("/api/v1/is-authenticated", app.checkAuthentication)
	mux.Post("/api/v1/forgot-password", app.sendPasswordResetEmail)
	mux.Post("/api/v1/reset-password", app.resetPassword)

	mux.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(app.Auth)

		r.Post("/virtual-terminal-succeeded", app.virtualTerminalPaymentSucceeded)

		r.Post("/all-sales", app.allSales)
		r.Post("/all-subscriptions", app.allSubscriptions)
		r.Post("/get-sale/{id}", app.getSale)

		r.Post("/refund", app.refundCharge)
		r.Post("/cancel-subscription", app.cancelSubscription)

		r.Post("/all-users", app.allUsers)
		r.Post("/all-users/{id}", app.oneUser)
		r.Post("/all-users/edit/{id}", app.editUser)
		r.Post("/all-users/delete/{id}", app.deleteUser)
	})

	return mux
}
