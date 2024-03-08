package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnusunil243/Job-Portal-User-service/db"
	"github.com/vishnusunil243/Job-Portal-User-service/initializer"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(err.Error())
	}
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	services := initializer.Initializer(DB)
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, services)
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen on port 8081")
	}
	log.Printf("user service listening on port 8081")
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to listen on port 8081")
	}
}
