package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	appAttachment "github.com/nevinmanoj/hostmate/internal/app/attachment"
	appBooking "github.com/nevinmanoj/hostmate/internal/app/booking"
	appPayemnt "github.com/nevinmanoj/hostmate/internal/app/payment"
	appProperty "github.com/nevinmanoj/hostmate/internal/app/property"
	appUser "github.com/nevinmanoj/hostmate/internal/app/user"

	domainAccess "github.com/nevinmanoj/hostmate/internal/domain/access"
	domainAttachment "github.com/nevinmanoj/hostmate/internal/domain/attachment"
	domainBooking "github.com/nevinmanoj/hostmate/internal/domain/booking"
	domainPayment "github.com/nevinmanoj/hostmate/internal/domain/payment"
	domainProperty "github.com/nevinmanoj/hostmate/internal/domain/property"
	domainUser "github.com/nevinmanoj/hostmate/internal/domain/user"

	"github.com/nevinmanoj/hostmate/internal/db/azure"
	postgres "github.com/nevinmanoj/hostmate/internal/db/postgres"
	repoAccess "github.com/nevinmanoj/hostmate/internal/db/postgres/access"
	repoAttachment "github.com/nevinmanoj/hostmate/internal/db/postgres/attachment"
	repoBooking "github.com/nevinmanoj/hostmate/internal/db/postgres/booking"
	repoPayment "github.com/nevinmanoj/hostmate/internal/db/postgres/payment"
	repoProperty "github.com/nevinmanoj/hostmate/internal/db/postgres/property"
	repoUser "github.com/nevinmanoj/hostmate/internal/db/postgres/user"

	middleware "github.com/nevinmanoj/hostmate/internal/middleware"
)

func Start() error {
	//Router and db connection
	var r *chi.Mux = chi.NewRouter()

	//get connection strings and jwt secret
	dsn := os.Getenv("DATABASE_URL")
	azurestr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtSecretbyte := []byte(jwtSecret)

	//postgres
	dbConn := postgres.NewPostgres(dsn)

	//Azure blob client
	azureBlobClient, err := azure.NewAzureBlobClient(azurestr)
	if err != nil {
		return err
	}

	// Global middleware
	r.Use(chimiddle.StripSlashes)

	//auth middleware
	authMiddleware := middleware.Authorization(jwtSecretbyte)

	//Blob storage
	blobStorage := azure.NewBlobStorage(azureBlobClient)

	//Repos
	userReadRepo := repoUser.NewUserReadRepository(dbConn)
	userWriteRepo := repoUser.NewUserWriteRepository(dbConn)
	accessRepo := repoAccess.NewAccessRepository(dbConn)
	propertyReadRepo := repoProperty.NewPropertyReadRepository(dbConn)
	propertyWriteRepo := repoProperty.NewPropertyWriteRepository(dbConn)
	bookingReadRepo := repoBooking.NewBookingReadRepository(dbConn)
	bookingWriteRepo := repoBooking.NewBookingWriteRepository(dbConn)
	paymentWriteRepo := repoPayment.NewPaymentWriteRepository(dbConn)
	attachmentWriteRepo := repoAttachment.NewAttachmentWriteRepository(dbConn)

	//Services
	userService := domainUser.NewUserService(userWriteRepo, jwtSecretbyte)
	accessService := domainAccess.NewAccessService(accessRepo)
	propertyService := domainProperty.NewPropertyService(propertyWriteRepo, userReadRepo, accessService)
	bookingService := domainBooking.NewBookingService(bookingWriteRepo, propertyReadRepo, accessService)
	paymentService := domainPayment.NewPaymentService(paymentWriteRepo, accessService, userReadRepo, bookingReadRepo, propertyReadRepo)
	attachmentService := domainAttachment.NewAttachmentService(attachmentWriteRepo, accessService, blobStorage)

	//Handlers
	userHandler := appUser.NewUserHandler(userService)
	propertyHandler := appProperty.NewPropertyHandler(propertyService)
	bookingHandler := appBooking.NewBookingHandler(bookingService)
	paymentHandler := appPayemnt.NewPaymentHandler(paymentService)
	attachmentHandler := appAttachment.NewAttachmentHandler(attachmentService)

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
		// router.Get("/{bookingId}/attachments", bookingHandler.GetBooking)
		router.Get("/{bookingId}/payments", paymentHandler.GetPaymentsWithBookingId)
		router.Post("/{bookingId}/payments", paymentHandler.CreatePayment)
		router.Put("/{bookingId}/payments/{paymentId}", paymentHandler.UpdatePayment)
	})

	//Payment routes
	r.Route("/payments", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Get("/", paymentHandler.GetPayments)
		router.Get("/{paymentId}", paymentHandler.GetPayment)
		// router.Get("/{paymentId}/attachments", paymentHandler.GetPayment)

	})

	// attachments routes
	r.Route("/attachments", func(router chi.Router) {
		router.Use(authMiddleware)
		router.Post("/request-upload", attachmentHandler.RequestUploadURL)
		router.Post("/confirm-upload", attachmentHandler.ConfirmUpload)
	})
	return http.ListenAndServe(":8080", r)
}
