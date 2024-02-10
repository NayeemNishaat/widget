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

		r.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)

		r.Post("/all-sales", app.AllSales)
		r.Post("/all-subscriptions", app.AllSubscriptions)
		r.Post("/get-sale/{id}", app.GetSale)

		r.Post("/refund", app.RefundCharge)
		r.Post("/cancel-subscription", app.CancelSubscription)

		r.Post("/all-users", app.AllUsers)
		r.Post("/all-users/{id}", app.OneUser)
		r.Post("/all-users/edit/{id}", app.EditUser)
		r.Post("/all-users/delete/{id}", app.DeleteUser)
	})

	return mux
}
