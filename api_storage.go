package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
)

func upload_to_storage() {
	accessKey := "YCAJEte2m-fdguRU29QZc2fyB"
	secretKey := "YCM6R35GZwklgQJH5q7trXPgKMBaMNTR7ukD4dyQ"
	endpoint := "https://storage.yandexcloud.net"

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ru-central1"),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	svc := s3.New(sess)

	// Создать новый бакет
	//_, err = svc.CreateBucket(&s3.CreateBucketInput{
	//	Bucket: aws.String("test123speech"),
	//})
	//
	//if err != nil {
	//	log.Fatalf("failed to create bucket: %v", err)
	//}

	// Загрузить объекты в бакет
	file, err := os.Open("mic.ogg")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("cipher"),
		Key:    aws.String("mic.ogg"),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}

	// Получить список объектов в бакете
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String("cipher")})
	if err != nil {
		log.Fatalf("failed to list objects: %v", err)
	}

	for _, item := range resp.Contents {
		if *item.Key == "mic.ogg" {
			log.Println("File: ", *item.Key)
		}
	}

}
