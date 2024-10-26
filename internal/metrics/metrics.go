package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ordersGivenCounter prometheus.Counter
}

var ordersGivenCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "orders_given_total",
		Help: "Total number of orders were given",
	},
)

func AddOrdersGiven(delta int) {
	ordersGivenCounter.Add(float64(delta))
}

func Init() {
	prometheus.MustRegister(ordersGivenCounter)
}
