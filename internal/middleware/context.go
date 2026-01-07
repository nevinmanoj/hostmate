package middleware

type contextKey string

const (
	ContextUserKey   contextKey = "userID"
	ContextjwtSecret contextKey = "jwtSecret"
)
