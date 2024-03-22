package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnusunil243/Job-Portal-User-service/db"
	"github.com/vishnusunil243/Job-Portal-User-service/initializer"
	"github.com/vishnusunil243/Job-Portal-User-service/internal/service"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("main function")
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf(err.Error())
	}
	addr := os.Getenv("DB_KEY")
	DB, err := db.InitDB(addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	companyConn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Fatal("error connecting to company service")
	}
	emailConn, err := grpc.Dial("localhost:8087", grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while connecting to notification service")
	}
	defer func() {
		companyConn.Close()
		emailConn.Close()
	}()

	companyRes := pb.NewCompanyServiceClient(companyConn)
	notificationRes := pb.NewEmailServiceClient(emailConn)
	service.NotificationClient = notificationRes
	service.CompanyClient = companyRes
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
