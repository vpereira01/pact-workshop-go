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

	// Verify the Provider - Tag-based Published Pacts for any known consumers
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: fmt.Sprintf("http://%s", ln.Addr().String()),
		Tags:            []string{"master"},
		PactURLs:        []string{filepath.FromSlash(fmt.Sprintf("%s/goadminservice-gouserservice.json", os.Getenv("PACT_DIR")))},
		ProviderVersion: "1.0.0",
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
