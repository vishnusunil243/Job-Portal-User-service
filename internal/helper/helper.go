package helper

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func CompareHashedPassword(hashedPass, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
func GetNumberInString(s string) int {
	num := 0
	if len(s) == 0 {
		return 0
	}
	for _, sr := range s {
		if unicode.IsNumber(sr) {
			num = int(sr - '0')
		}
	}
	return num
}
func ConvertStringToDate(str string) (time.Time, error) {
	layout := "02-01-2006"
	date, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("please provide a valid start date")
	}
	return date, nil
}
func GenerateRoomId() string {
	return uuid.New().String()
}
func ConvertStringToTimeStamp(str string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("please provide a valid start date")
	}
	return date, nil
}
func MinioUpload(name string, imageData []byte) (string, error) {
	minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESSKEY"), os.Getenv("MINIO_SECRETKEY"), ""),
		Secure: false,
	})
	if err != nil {
		log.Print("error while initialising minio", err)
		return "", err
	}
	objectName := "images/" + name
	contentType := `image/jpeg`
	n, err := minioClient.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error while uploading to minio", err)
		return "", err
	}
	log.Printf("Successfully uploaded %s of size %v\n", objectName, n)
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), os.Getenv("BUCKET_NAME"), objectName, time.Second*24*60*60, nil)
	if err != nil {
		log.Println("error while generating presigned URL", err)
		return "", err
	}
	return presignedURL.String(), nil
}
func DialGrpc(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}
