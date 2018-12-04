package main

import (
	"context"
	"github.com/vllry/k2/api/0.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"log"
	"net"
)

type workloadApiServer struct {
	db *dbConn
}

func startWorkloadApiServer(db *dbConn) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterWorkloadApiServer(s, &workloadApiServer{db})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *workloadApiServer) WorkloadApply(ctx context.Context, query *api.Workload) (*api.WorkloadSubmitted, error) {
	var spec WorkloadSpec
	err := yaml.Unmarshal([]byte(query.Yaml), &spec)
	if err != nil {
		return &api.WorkloadSubmitted{Success: false}, err
	}

	err = submitWorkload(s.db, &spec)
	if err != nil {
		return &api.WorkloadSubmitted{Success: false}, err
	}
	return &api.WorkloadSubmitted{Success: true}, nil
}
