package sample

import (
	"math/rand"
	"time"

	"github.com/ra111eo/TestTaskITgo/service"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomID() string {
	id, _ := service.GetNewId()
	return id
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
