package metrics

import (
	"context"
	"strings"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
)

// ------------------------------------------------------------
// Count every gRPC request.
//
// We want to answer:
//
// "How many requests did this service receive?"
// ------------------------------------------------------------
var grpcRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests.",
	},
	[]string{
		"method", // CreateOrder, GetOrder...
		"status", // OK, NotFound, Internal...
	},
)

// ------------------------------------------------------------
// Measure how long every request takes.
//
// We want to answer:
//
// "Is CreateOrder becoming slower?"
// ------------------------------------------------------------
var grpcRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Time taken by gRPC requests.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{
		"method",
	},
)

// ------------------------------------------------------------
// Prometheus does NOT automatically know about our metrics.
//
// We must introduce them.
//
// "Hello Prometheus, please watch these."
// ------------------------------------------------------------
func Register() {
	prometheus.MustRegister(
		grpcRequestsTotal,
		grpcRequestDuration,
	)
}

// ------------------------------------------------------------
//
// Instead of putting metric code inside every handler:
//
//	CreateOrder()
//	GetOrder()
//	CancelOrder()
//
// We stand at the front door.
//
// Every request enters through this door.
//
// We measure:
//
// - What method?
// - How long?
// - Success or failure?
//
// Then we let the request continue.
// ------------------------------------------------------------
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {

		// ----------------------------------------
		//
		// Remember when request started.
		//
		// We need this later to calculate
		// how long it took.
		// ----------------------------------------
		start := time.Now()

		// ----------------------------------------
		//
		// defer runs AFTER CreateOrder()
		// finishes.
		//
		// Doesn't matter whether it
		// succeeded or failed.
		//
		// We always record metrics.
		// ----------------------------------------
		defer func() {

			// Example:
			//
			// /grpc.OrderService/CreateOrder
			//
			// becomes
			//
			// grpc.OrderService/CreateOrder
			method := normalizeMethod(info.FullMethod)

			// Convert Go error into
			// gRPC status.
			//
			// Example:
			//
			// OK
			// Internal
			// NotFound
			statusLabel := status.Code(err).String()

			// Count one request.
			grpcRequestsTotal.
				WithLabelValues(method, statusLabel).
				Inc()

			// Record request duration.
			grpcRequestDuration.
				WithLabelValues(method).
				Observe(time.Since(start).Seconds())
		}()

		// ----------------------------------------
		//
		// Finally...
		//
		// Call the REAL gRPC handler.
		//
		// We don't process business logic.
		//
		// We only observe it.
		// ----------------------------------------
		return handler(ctx, req)
	}
}

// ------------------------------------------------------------
//
// gRPC gives:
//
// /grpc.OrderService/CreateOrder
//
// We remove the leading "/" because
// it looks nicer in Prometheus.
// ------------------------------------------------------------
func normalizeMethod(fullMethod string) string {
	return strings.TrimPrefix(fullMethod, "/")
}
