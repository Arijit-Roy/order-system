package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	orderspb "order-system/grpc"
	"order-system/internal/config"
	"order-system/internal/repository"
	"order-system/internal/service"
	grpctransport "order-system/internal/transport/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"order-system/internal/metrics"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.Load()
	pgDSN := cfg.PostgresDSN()

	pool, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepo(pool)
	svc := service.NewOrderService(repo)

	metrics.Register()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpcserver.NewServer(
		grpcserver.UnaryInterceptor(metrics.UnaryServerInterceptor()),
	)

	go func() {
		mux := http.NewServeMux()

		mux.Handle("/metrics", promhttp.Handler())

		log.Println("Metrics server listening on :8080")

		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatal(err)
		}
	}()

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcSrv, healthServer)

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	orderspb.RegisterOrderServiceServer(grpcSrv, grpctransport.NewOrderGRPCServer(svc))
	reflection.Register(grpcSrv)

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Println("gRPC server listening on :50051")
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
