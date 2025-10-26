package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"github.com/gorilla/mux"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	deviceService service.DeviceService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService service.DeviceService) *Server {
	return &Server{
		listenAddress: listenAddress,
		deviceService: deviceService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	r := mux.NewRouter()

	// Set Content-Type header to application/json
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// Health check
	r.HandleFunc("/api/v0/health", s.Health).Methods(http.MethodGet)

	// Device management
	r.HandleFunc("/api/v0/devices", s.CreateDevice).Methods(http.MethodPost)

	// Transaction signing
	r.HandleFunc("/api/v0/devices/{deviceId}/sign", s.SignTransaction).Methods(http.MethodPost)

	// Device retrieval
	r.HandleFunc("/api/v0/devices/{deviceId}", s.GetDevice).Methods(http.MethodGet)

	// Device retrieval
	r.HandleFunc("/api/v0/devices", s.GetAllDevices).Methods(http.MethodGet)

	return http.ListenAndServe(s.listenAddress, r)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
