package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	appBooking "github.com/nevinmanoj/hostmate/internal/app/booking"
	appPayemnt "github.com/nevinmanoj/hostmate/internal/app/payment"
	appProperty "github.com/nevinmanoj/hostmate/internal/app/property"
	appUser "github.com/nevinmanoj/hostmate/internal/app/user"

	domainBooking "github.com/nevinmanoj/hostmate/internal/domain/booking"
	domainPayment "github.com/nevinmanoj/hostmate/internal/domain/payment"
	domainProperty "github.com/nevinmanoj/hostmate/internal/domain/property"
	domainUser "github.com/nevinmanoj/hostmate/internal/domain/user"

	postgres "github.com/nevinmanoj/hostmate/internal/db"
	repoBooking "github.com/nevinmanoj/hostmate/internal/db/booking"
	repoPayment "github.com/nevinmanoj/hostmate/internal/db/payment"
	repoProperty "github.com/nevinmanoj/hostmate/internal/db/property"
	repoUser "github.com/nevinmanoj/hostmate/internal/db/user"

	middleware "github.com/nevinmanoj/hostmate/internal/middleware"
)

func Start() error {
	//Router and db connection
	var r *chi.Mux = chi.NewRouter()

	dsn := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtSecretbyte := []byte(jwtSecret)
	dbConn := postgres.NewPostgres(dsn)

	// Global middleware
	r.Use(chimiddle.StripSlashes)

	//auth middleware
	authMiddleware := middleware.Authorization(jwtSecretbyte)

	//Repos
	userReadRepo := repoUser.NewUserReadRepository(dbConn)
	userWriteRepo := repoUser.NewUserWriteRepository(dbConn)
	propertyReadRepo := repoProperty.NewPropertyReadRepository(dbConn)
	propertyWriteRepo := repoProperty.NewPropertyWriteRepository(dbConn)
	bookingReadRepo := repoBooking.NewBookingReadRepository(dbConn)
	bookingWriteRepo := repoBooking.NewBookingWriteRepository(dbConn)
	paymentWrieteRepo := repoPayment.NewPaymentWriteRepository(dbConn)

	//Services
	userService := domainUser.NewUserService(userWriteRepo, jwtSecretbyte)
	propertyService := domainProperty.NewPropertyService(propertyWriteRepo, userReadRepo)
	bookingService := domainBooking.NewBookingService(bookingWriteRepo, propertyReadRepo)
	paymentService := domainPayment.NewPaymentService(paymentWrieteRepo, userReadRepo, bookingReadRepo, propertyReadRepo)

	//Handlers
	userHandler := appUser.NewUserHandler(userService)
	propertyHandler := appProperty.NewPropertyHandler(propertyService)
	bookingHandler := appBooking.NewBookingHandler(bookingService)
	paymentHandler := appPayemnt.NewPaymentHandler(paymentService)

	//User routes
	r.Route("/users", func(router chi.Router) {
		router.Get("/{userId}", userHandler.GetUser)
		router.Post("/login", userHandler.LoginUser)
		router.Post("/register", userHandler.CreateUser)

	})

	//Property routes
	r.Route("/properties", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Get("/", propertyHandler.GetProperties)
		router.Get("/{propertyId}", propertyHandler.GetProperty)
		router.Post("/", propertyHandler.CreateProperty)
		router.Put("/{propertyId}", propertyHandler.UpdateProperty)
		router.Get("/{propertyId}/availability", bookingHandler.CheckAvailability)
		router.Get("/{propertyId}/payments", paymentHandler.GetPaymentsWithPropertyId)
	})

	//booking routes

	r.Route("/bookings", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Get("/", bookingHandler.GetBookings)
		router.Get("/{bookingId}", bookingHandler.GetBooking)
		router.Post("/", bookingHandler.CreateBooking)
		router.Put("/{bookingId}", bookingHandler.UpdateBooking)
		router.Get("/{bookingId}/payments", paymentHandler.GetPaymentsWithBookingId)
		router.Post("/{bookingId}/payments", paymentHandler.CreatePayment)
		router.Put("/{bookingId}/payments/{paymentId}", paymentHandler.UpdatePayment)
	})

	//Payment routes
	r.Route("/payments", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Get("/", paymentHandler.GetPayments)
		router.Get("/{paymentId}", paymentHandler.GetPayment)
	})

	// attachments routes
	r.Route("/attachments", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("attachments endpoint"))
		})
	})
	return http.ListenAndServe(":8080", r)
}
