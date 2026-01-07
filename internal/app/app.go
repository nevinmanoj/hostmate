package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	appbooking "github.com/nevinmanoj/hostmate/internal/app/booking"
	domainbooking "github.com/nevinmanoj/hostmate/internal/domain/booking"
	domainUser "github.com/nevinmanoj/hostmate/internal/domain/user"

	appProperty "github.com/nevinmanoj/hostmate/internal/app/property"
	domainProperty "github.com/nevinmanoj/hostmate/internal/domain/property"

	appUser "github.com/nevinmanoj/hostmate/internal/app/user"

	postgres "github.com/nevinmanoj/hostmate/internal/db"
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
	userReadRepo := postgres.NewUserReadRepository(dbConn)
	userWriteRepo := postgres.NewUserWriteRepository(dbConn)
	propertyReadRepo := postgres.NewPropertyReadRepository(dbConn)
	propertyWriteRepo := postgres.NewPropertyWriteRepository(dbConn)
	// bookingReadRepo := booking.NewBookingRepository(dbConn)
	bookingWriteRepo := postgres.NewBookingWriteRepository(dbConn)

	//Services
	userService := domainUser.NewUserService(userWriteRepo, jwtSecretbyte)
	propertyService := domainProperty.NewPropertyService(propertyWriteRepo, userReadRepo)
	bookingService := domainbooking.NewBookingService(bookingWriteRepo, propertyReadRepo)

	//Handlers
	userHandler := appUser.NewUserHandler(userService)
	propertyHandler := appProperty.NewPropertyHandler(propertyService)
	bookingHandler := appbooking.NewBookingHandler(bookingService)

	//User routes
	r.Route("/users", func(router chi.Router) {
		router.Get("/{id}", userHandler.GetUser)
		router.Post("/login", userHandler.LoginUser)
		router.Post("/register", userHandler.CreateUser)

	})

	//Property routes
	r.Route("/properties", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Get("/", propertyHandler.GetProperties)
		router.Get("/{id}", propertyHandler.GetProperty)
		router.Post("/", propertyHandler.CreateProperty)
		router.Put("/{id}", propertyHandler.UpdateProperty)
		router.Get("/{id}/availability", bookingHandler.CheckAvailability)

	})

	//booking routes

	r.Route("/bookings", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Get("/", bookingHandler.GetBookings)
		router.Get("/{id}", bookingHandler.GetBooking)
		router.Post("/", bookingHandler.CreateBooking)
		router.Put("/{id}", bookingHandler.UpdateBooking)
	})

	//Payment routes
	r.Route("/payments", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("payments endpoint"))
		})
	})

	// attachments routes
	r.Route("/attachments", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("attachments endpoint"))
		})
	})
	return http.ListenAndServe(":8080", r)
}
