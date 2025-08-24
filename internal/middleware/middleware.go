package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
)

type contextKey string

const userIDKey contextKey = "userID"
const userEmailKey contextKey = "userEmail"

// write a mock firebase auth middleware that checks the Authorization header for a Bearer token
// and just uses this token as the user ID in the context
func MockFirebaseAuthMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		userID := parts[1]

		// Add user ID to context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		// Add mock email from user id to context
		ctx = context.WithValue(ctx, userEmailKey, userID+"@example.com")
		r = r.WithContext(ctx)

		log.Printf("Mock Firebase Auth: User %s authenticated for %s", userID, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func FirebaseAuthMw(authClient *auth.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		idToken := parts[1]

		// Verify the Firebase ID token
		token, err := authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("Firebase token verification failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), userIDKey, token.UID)

		// Add user email to context
		// todo: find out how to actually do this
		ctx = context.WithValue(ctx, userEmailKey, token.Claims["email"])
		r = r.WithContext(ctx)

		log.Printf("Firebase Auth: User %s authenticated for %s", token.UID, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

func GetUserEmail(r *http.Request) string {
	if userEmail, ok := r.Context().Value(userEmailKey).(string); ok {
		return userEmail
	}
	return ""
}

func JsonMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func CorsMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func ApacheLogMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
			size:           0,
		}

		next.ServeHTTP(rw, r)

		clientIP := r.RemoteAddr
		if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			clientIP = strings.Split(forwardedFor, ",")[0]
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			clientIP = realIP
		}

		userID := "-"
		if uid := GetUserID(r); uid != "" {
			userID = uid
		}

		timestamp := start.Format("02/Jan/2006:15:04:05 -0700")
		duration := time.Since(start)

		log.Printf("%s - %s [%s] \"%s %s %s\" %d %d %v",
			clientIP,
			userID,
			timestamp,
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			rw.statusCode,
			rw.size,
			duration,
		)
	})
}
