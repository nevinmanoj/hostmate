package attachment

type AttachmentParentType string

const (
	AttachmentParentPayment AttachmentParentType = "payments"
	AttachmentParentBooking AttachmentParentType = "bookings"
)
