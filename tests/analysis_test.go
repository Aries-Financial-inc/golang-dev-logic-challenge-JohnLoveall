package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-JohnLoveall/model"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-JohnLoveall/routes"
	"github.com/stretchr/testify/assert"
)

func TestOptionsContractModelValidation(t *testing.T) {
	// Test valid OptionsContract
	validContract := model.OptionsContract{
		Type:           model.Call,
		StrikePrice:    100,
		Bid:            5,
		Ask:            6,
		ExpirationDate: time.Now().AddDate(0, 1, 0), // 1 month in the future
		LongShort:      model.Long,
	}

	// Manually validate the model (usually done in a handler)
	assert.Equal(t, model.Call, validContract.Type)
	assert.Equal(t, 100.0, validContract.StrikePrice)
	assert.Equal(t, 5.0, validContract.Bid)
	assert.Equal(t, 6.0, validContract.Ask)
	assert.True(t, validContract.ExpirationDate.After(time.Now()))
	assert.Equal(t, model.Long, validContract.LongShort)
}

func TestAnalysisEndpoint(t *testing.T) {
	router := routes.SetupRouter()

	// Create a sample request body
	contracts := []model.OptionsContract{
		{
			Type:           model.Call,
			StrikePrice:    100,
			Bid:            5,
			Ask:            6,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Long,
		},
		{
			Type:           model.Put,
			StrikePrice:    90,
			Bid:            3,
			Ask:            4,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Short,
		},
	}
	requestBody, _ := json.Marshal(contracts)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// // Check if the response contains the expected fields
	assert.Contains(t, response, "graph_data")
	assert.Contains(t, response, "max_profit")
	assert.Contains(t, response, "max_loss")
	assert.Contains(t, response, "break_even_points")
}

func TestIntegration(t *testing.T) {
	router := routes.SetupRouter()

	// Simulate a broader scenario
	// Step 1: Create a batch of sample request bodies and perform analysis on each
	contractsBatch := [][]model.OptionsContract{
		{
			{
				Type:           model.Call,
				StrikePrice:    100,
				Bid:            5,
				Ask:            6,
				ExpirationDate: time.Now().AddDate(0, 1, 0),
				LongShort:      model.Long,
			},
			{
				Type:           model.Put,
				StrikePrice:    90,
				Bid:            3,
				Ask:            4,
				ExpirationDate: time.Now().AddDate(0, 1, 0),
				LongShort:      model.Short,
			},
		},
		{
			{
				Type:           model.Call,
				StrikePrice:    110,
				Bid:            6,
				Ask:            7,
				ExpirationDate: time.Now().AddDate(0, 1, 0),
				LongShort:      model.Long,
			},
			{
				Type:           model.Put,
				StrikePrice:    85,
				Bid:            2,
				Ask:            3,
				ExpirationDate: time.Now().AddDate(0, 1, 0),
				LongShort:      model.Short,
			},
		},
	}

	for _, contracts := range contractsBatch {
		requestBody, _ := json.Marshal(contracts)

		// Create a new HTTP request
		req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check if the response contains the expected fields
		assert.Contains(t, response, "graph_data")
		assert.Contains(t, response, "max_profit")
		assert.Contains(t, response, "max_loss")
		assert.Contains(t, response, "break_even_points")
	}

	// Step 2: Test with more than 4 contracts
	tooManyContracts := []model.OptionsContract{
		{
			Type:           model.Call,
			StrikePrice:    100,
			Bid:            5,
			Ask:            6,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Long,
		},
		{
			Type:           model.Put,
			StrikePrice:    90,
			Bid:            3,
			Ask:            4,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Short,
		},
		{
			Type:           model.Call,
			StrikePrice:    110,
			Bid:            6,
			Ask:            7,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Long,
		},
		{
			Type:           model.Put,
			StrikePrice:    85,
			Bid:            2,
			Ask:            3,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Short,
		},
		{
			Type:           model.Call,
			StrikePrice:    120,
			Bid:            7,
			Ask:            8,
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			LongShort:      model.Long,
		},
	}

	requestBody, _ := json.Marshal(tooManyContracts)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check if the response contains the expected error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Maximum of 4 options contracts allowed", response["error"])
}
