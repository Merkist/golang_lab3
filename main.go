package main

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}

func trapezoidalIntegral(function func(float64) float64, start, end float64, intervals int) float64 {
	h := (end - start) / float64(intervals)
	integral := (function(start) + function(end)) / 2.0

	for i := 1; i < intervals; i++ {
		x := start + float64(i)*h
		integral += function(x)
	}

	return integral * h
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/calculate", func(c *gin.Context) {
		const PI = math.Pi
		values := []string{"Power", "Error1", "Error2", "Price"}
		inputs := make(map[string]float64)
		for _, v := range values {
			val, err := strconv.ParseFloat(c.PostForm(v), 64)
			if err != nil {
				c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Invalid input for " + v})
				return
			}
			inputs[v] = val
		}

		if inputs["Error2"] < 0 && inputs["Error2"] > 100 {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Допустима похибка повинна бути від 0 до 100%"})
			return
		}

		temp := inputs["Power"] * (inputs["Error2"] / 100)
		integralA := inputs["Power"] - temp
		integralB := inputs["Power"] + temp

		// Функція для обчислення ймовірності
		function := func(p float64) float64 {
			return 1 / (inputs["Error1"] * math.Sqrt(2*PI)) * math.Exp(-((p-inputs["Power"])*(p-inputs["Power"]))/(2*inputs["Error1"]*inputs["Error1"]))
		}

		qw1 := trapezoidalIntegral(function, integralA, integralB, 1000)
		w1 := inputs["Power"] * 24 * qw1 * inputs["Price"]
		w2 := inputs["Power"] * 24 * (1 - qw1) * inputs["Price"]
		profit := w1 - w2

		c.HTML(http.StatusOK, "index.html", gin.H{
			"profit": roundToTwoDecimalPlaces(profit),
		})
	})

	r.Run(":8080")
}
