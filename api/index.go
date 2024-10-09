package api

import (
	"api-irr/resolver"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Structs for request and response
type IRRRequest struct {
	Spending []float64 `json:"spending"`
	Income   []float64 `json:"income"`
	Code     string    `json:"code"`
}

type IRRResponse struct {
	Status int     `json:"status"`
	Irr    float64 `json:"irr"`
	Error  string  `json:"error"`
}

var (
	app *gin.Engine
)

func myRoute(r *gin.RouterGroup) {
	r.GET("", func(c *gin.Context) {
		c.String(http.StatusOK, "API ready to use!")
	})

	r.POST("/solve", resolveIRR)
}

func init() {
	app = gin.New()
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r := app.Group("")
	myRoute(r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

func resolveIRR(c *gin.Context) {
	var request IRRRequest
	// Bind the JSON request to the IRRRequest struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the request data to ensure it's correctly parsed
	fmt.Printf("Received request: %+v\n", request)

	if request.Code != "resolve" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your request is rejected."})
		return
	}

	// Perform IRR calculation
	length := len(request.Spending)
	if length != len(request.Income) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number of income and outcome is not equal."})
		return
	}

	// Check for input validity
	if length < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number of income and outcome is less than 2."})
		return
	}

	diff := []float64{}
	for i := 0; i < length; i++ {
		diff = append(diff, request.Income[i]-request.Spending[i])
	}

	// Call RealRoots and handle potential errors
	roots, err := resolver.RealRoots(length-1, diff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, IRRResponse{
			Status: 1,
			Error:  err.Error(),
		})
		return
	}

	// Find the IRR based on the root
	var v float64
	for _, root := range roots {
		if root >= 0 && root < 1 {
			v = root
			break
		}
	}

	irr := resolver.GetIRR(v)

	// Send response based on whether IRR is NaN or not
	if !math.IsNaN(irr) {
		c.JSON(http.StatusOK, IRRResponse{
			Status: 0,
			Irr:    irr,
		})
	} else {
		c.JSON(http.StatusOK, IRRResponse{
			Status: 1,
			Irr:    math.NaN(),
			Error:  "IRR calculation resulted in NaN.",
		})
	}
}
