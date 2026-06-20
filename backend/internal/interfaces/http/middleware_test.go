package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpapi "github.com/min-legomain/nikkei-trade-journal/backend/internal/interfaces/http"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

// TestInternalOnly_ValidTokenAllowsRequest guarantees that a request carrying
// the correct Bearer token is forwarded to the next handler.
func TestInternalOnly_ValidTokenAllowsRequest(t *testing.T) {
	// Arrange
	h := httpapi.InternalOnly("secret")(okHandler)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w.Code)
	}
}

// TestInternalOnly_InvalidTokenReturns401 guarantees that a wrong token is
// rejected and the inner handler is never reached.
func TestInternalOnly_InvalidTokenReturns401(t *testing.T) {
	// Arrange
	h := httpapi.InternalOnly("secret")(okHandler)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", w.Code)
	}
}

// TestInternalOnly_MissingTokenReturns401 guarantees that a request without
// any Authorization header is also rejected.
func TestInternalOnly_MissingTokenReturns401(t *testing.T) {
	// Arrange
	h := httpapi.InternalOnly("secret")(okHandler)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	// Act
	h.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d, want 401", w.Code)
	}
}
