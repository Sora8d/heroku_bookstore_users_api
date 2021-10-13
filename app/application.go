package app

import (
	"github.com/gin-gonic/gin"
)

//This roter creates a go routine for every request handled, so they shouldnt have common variables
var router = gin.Default()

// The http server is going to be only here and in controller
func StartApplication() {
	mapUrls()

	logger.Info("starting app...")
	router.Run("localhost:8080")
}
