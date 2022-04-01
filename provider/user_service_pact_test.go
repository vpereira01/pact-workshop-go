package provider

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-workshop-go/model"
	"github.com/pact-foundation/pact-workshop-go/provider/repository"
)

func TestPactProvider(t *testing.T) {
	// Setup Pact and related test stuff
	pact := dsl.Pact{
		Consumer: os.Getenv("CONSUMER_NAME"),
		Provider: os.Getenv("PROVIDER_NAME"),
		LogDir:   os.Getenv("LOG_DIR"),
		PactDir:  os.Getenv("PACT_DIR"),
		LogLevel: "INFO",
	}

	ln := startInstrumentedProvider()
	defer ln.Close()

	stateHandlers := types.StateHandlers{
		"User 10 exists": func() error {
			userRepository = &repository.UserRepository{
				Users: map[string]*model.User{
					"sally": {
						FirstName: "Jean-Marie",
						LastName:  "de La Beaujardière😀😍",
						Username:  "sally",
						Type:      "admin",
						ID:        10,
					},
				},
			}
			return nil
		},
		"User 10 does not exist": func() error {
			userRepository = &repository.UserRepository{}
			return nil
		},
	}

	// Verify the Provider - Tag-based Published Pacts for any known consumers
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://%s", ln.Addr().String()),
		Tags:            []string{"master"},
		PactURLs:        []string{filepath.FromSlash(fmt.Sprintf("%s/goadminservice-gouserservice.json", os.Getenv("PACT_DIR")))},
		ProviderVersion: "1.0.0",
		StateHandlers:   stateHandlers,
	})

	if err != nil {
		t.Log(err)
	}
}

func startInstrumentedProvider() net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("API starting: port %s", ln.Addr())

	mux := GetHTTPHandler()
	go http.Serve(ln, mux)

	return ln
}
