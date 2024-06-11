package controllers

import (
	"math"
	"net/http"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-JohnLoveall/model"
	"github.com/gin-gonic/gin"
)

// AnalysisResult structure for the response body
type AnalysisResult struct {
	GraphData       []GraphPoint `json:"graph_data"`
	MaxProfit       float64      `json:"max_profit"`
	MaxLoss         float64      `json:"max_loss"`
	BreakEvenPoints []float64    `json:"break_even_points"`
}

// GraphPoint structure for X & Y values of the risk & reward graph
type GraphPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// AnalysisHandler is the HTTP handler for analyzing options contracts
func AnalysisHandler(c *gin.Context) {
	var contracts []model.OptionsContract

	// Bind the JSON payload to the contracts slice
	if err := c.ShouldBindJSON(&contracts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the number of contracts does not exceed the limit of 4
	if len(contracts) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum of 4 options contracts allowed"})
		return
	}

	// Calculate the graph points, max profit, max loss, and break-even points
	graphData := calculateGraphPoints(contracts)
	maxProfit := calculateMaxProfit(graphData)
	maxLoss := calculateMaxLoss(graphData)
	breakEvenPoints := calculateBreakEvenPoints(graphData)

	// Create and send the response
	response := AnalysisResult{
		GraphData:       graphData,
		MaxProfit:       maxProfit,
		MaxLoss:         maxLoss,
		BreakEvenPoints: breakEvenPoints,
	}

	c.JSON(http.StatusOK, response)
}

// calculateGraphPoints calculates the X & Y values for the risk & reward graph
func calculateGraphPoints(contracts []model.OptionsContract) []GraphPoint {
	var graphData []GraphPoint

	// Iterate over a range of underlying prices
	for i := 0; i <= 200; i++ {
		x := float64(i)
		y := 0.0

		// Calculate the total payoff for the current price
		for _, contract := range contracts {
			payoff := calculatePayoff(contract, x)
			if contract.LongShort == model.Long {
				y += payoff
			} else {
				y -= payoff
			}
		}

		// Append the calculated point to the graph data
		graphData = append(graphData, GraphPoint{X: x, Y: y})
	}
	return graphData
}

// calculatePayoff calculates the payoff of an options contract at a given underlying price
func calculatePayoff(contract model.OptionsContract, underlyingPrice float64) float64 {
	if contract.Type == model.Call {
		return math.Max(0, underlyingPrice-contract.StrikePrice) - contract.Ask
	} else {
		return math.Max(0, contract.StrikePrice-underlyingPrice) - contract.Ask
	}
}

// calculateMaxLoss calculates the maximum possible loss from the graph data
func calculateMaxProfit(graphData []GraphPoint) float64 {
	maxProfit := -math.MaxFloat64

	// Find the maximum Y value in the graph data
	for _, point := range graphData {
		if point.Y > maxProfit {
			maxProfit = point.Y
		}
	}
	return maxProfit
}

func calculateMaxLoss(graphData []GraphPoint) float64 {
	maxLoss := math.MaxFloat64

	// Find the minimum Y value in the graph data
	for _, point := range graphData {
		if point.Y < maxLoss {
			maxLoss = point.Y
		}
	}
	return maxLoss
}

// calculateBreakEvenPoints calculates all break-even points from the graph data
func calculateBreakEvenPoints(graphData []GraphPoint) []float64 {
	var breakEvenPoints []float64

	// Iterate through the graph data to find where the Y value crosses zero
	for i := 1; i < len(graphData); i++ {
		if (graphData[i-1].Y < 0 && graphData[i].Y > 0) || (graphData[i-1].Y > 0 && graphData[i].Y < 0) {
			breakEvenPoints = append(breakEvenPoints, graphData[i].X)
		}
	}
	return breakEvenPoints
}
