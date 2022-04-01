//go:build pact
// +build pact

package client_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	target "github.com/pact-foundation/pact-workshop-go/consumer/client"
	"github.com/pact-foundation/pact-workshop-go/model"
)

var (
	pact          dsl.Pact
	client        *target.Client
	commonHeaders dsl.MapMatcher
)

func TestMain(m *testing.M) {
	var exitCode int

	// Setup Pact and related test stuff
	pact = dsl.Pact{
		Consumer: os.Getenv("CONSUMER_NAME"),
		Provider: os.Getenv("PROVIDER_NAME"),
		LogDir:   os.Getenv("LOG_DIR"),
		PactDir:  os.Getenv("PACT_DIR"),
		LogLevel: "INFO",
	}

	pact.Setup(true)

	u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", pact.Server.Port))

	client = &target.Client{
		BaseURL: u,
	}

	commonHeaders = dsl.MapMatcher{
		"Content-Type":         dsl.Term("application/json; charset=utf-8", `application\/json`),
		"X-Api-Correlation-Id": dsl.Like("100"),
	}

	// Run all the tests
	exitCode = m.Run()

	// Shutdown the Mock Service and Write pact files to disk
	if err := pact.WritePact(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pact.Teardown()
	os.Exit(exitCode)
}

func TestClientPact_GetUserExist(t *testing.T) {
	// arrange
	pact.
		AddInteraction().
		Given("User 10 exists").
		UponReceiving("A request to get user with id 10").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/user/10", "/user/[0-9]+"),
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Body:    dsl.Match(model.User{ID: 10}),
			Headers: commonHeaders,
		})

	// act & assert
	err := pact.Verify(func() error {
		user, err := client.GetUser(10)

		// Assert basic fact
		if user.ID != 10 {
			return fmt.Errorf("wanted user with ID %d but got %d", 10, user.ID)
		}

		return err
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}

func TestClientPact_GetUserNotExist(t *testing.T) {
	// arrange
	pact.
		AddInteraction().
		Given("User 10 does not exist").
		UponReceiving("A request to get user with id 10").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.Term("/user/10", "/user/[0-9]+"),
		}).
		WillRespondWith(dsl.Response{
			Status:  404,
			Headers: commonHeaders,
		})

	// act & assert
	err := pact.Verify(func() error {
		_, err := client.GetUser(10)

		// Assert basic fact
		if !errors.Is(err, target.ErrNotFound) {
			return fmt.Errorf("expected error %s but got %s", target.ErrNotFound, err)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}
}
