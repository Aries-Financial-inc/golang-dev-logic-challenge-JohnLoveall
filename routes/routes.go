package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-JohnLoveall/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/analyze", controllers.AnalysisHandler)

	return router
}
