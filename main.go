package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"simplebank.com/api"
	db "simplebank.com/db/sqlgen"
	"simplebank.com/gapi"
	"simplebank.com/pb"
	"simplebank.com/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	go func() {
		runGrpcServer(config, store)
	}()
	runGinServer(config, store)
}

func runGrpcServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start GRPC server:", err)
	}
}

// func runGatewayServer(config utils.Config, store db.Store) {
// 	server, err := gapi.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("cannot create server:", err)
// 	}

// 	grpcMux := runtime.NewServeMux()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	pb.RegisterSimpleBankHandlerClient(ctx, grpcMux, server)

// 	listener, err := net.Listen("tcp", config.GrpcServerAddress)
// 	if err != nil {
// 		log.Fatal("cannot start server:", err)
// 	}

// 	log.Printf("start gRPC server on %s", listener.Addr().String())
// 	err = grpcServer.Serve(listener)
// 	if err != nil {
// 		log.Fatal("cannot start GRPC server:", err)
// 	}
// }

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
