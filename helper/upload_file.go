package helper

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context, fileKey, destinationDir string) string {
	// allowed extension
	allowedExtensions := map[string]bool{
		"zip": true,
		"rar": true,
	}

	// config aws s3 bucket
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return ""
	}

	client := s3.NewFromConfig(cfg)

	// file request handler
	file, header, err := c.Request.FormFile(fileKey)
	if err != nil {
		if err == http.ErrMissingFile {
			return ""
		}
		return ""
	}
	defer file.Close()

	// generate filename
	filename := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(filename))[1:]

	if !allowedExtensions[ext] {
		return ""
	}

	filename = GenerateUniqueFile(filename)

	// s3 upload
	uploader := manager.NewUploader(client)
	_, postErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("notarius"),
		Key:    aws.String(destinationDir + "/" + filename),
		Body:   file,
		ACL:    "private",
	})

	if postErr != nil {
		log.Printf("UploadFile, error uploading s3, %v", err)
	}

	// return
	return filename
}

func UpdateFile(c *gin.Context, fileKey, destinationDir string, oldFile string) string {
	allowedExtensions := map[string]bool{
		"zip": true,
		"rar": true,
	}

	/*
		|--------------------------------------------------------------------------
		| Load AWS S3 config
		|--------------------------------------------------------------------------
	*/

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Gagal memuat konfigurasi AWS: %v", err)
		return ""
	}
	client := s3.NewFromConfig(cfg)

	file, header, err := c.Request.FormFile(fileKey)
	if err != nil {
		log.Printf("Gagal membaca file dari request: %v", err)
		return ""
	}
	defer file.Close()

	/*
		|--------------------------------------------------------------------------
		| Periksa ekstensi & destinasi file
		|--------------------------------------------------------------------------
	*/

	filename := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(filename))[1:]

	if !allowedExtensions[ext] {
		log.Printf("Ekstensi file tidak diperbolehkan: %s", ext)
		return ""
	}

	filename = GenerateUniqueFile(filename)

	filePath := destinationDir + "/" + filename
	oldPath := destinationDir + "/" + oldFile

	/*
		|--------------------------------------------------------------------------
		| Hapus file lama dari S3 sebelum mengunggah file baru
		|--------------------------------------------------------------------------
	*/

	if oldPath != "" {
		err = DeleteFile(client, "notarius", oldPath)

		if err != nil {
			log.Printf("Gagal menghapus file lama: %v", err)
		}
	}

	/*
		|--------------------------------------------------------------------------
		| Upload file baru ke S3
		|--------------------------------------------------------------------------
	*/

	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("notarius"),
		Key:    aws.String(filePath),
		Body:   file,
		ACL:    "private",
	})

	if err != nil {
		log.Printf("Gagal mengunggah file ke S3: %v", err)
		return ""
	}

	log.Printf("File berhasil diperbarui: %s", filePath)
	return filename
}

func DeleteFile(client *s3.Client, bucket, filePath string) error {
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filePath),
	})

	/*
		|--------------------------------------------------------------------------
		| Check if failed
		|--------------------------------------------------------------------------
	*/

	if err != nil {
		log.Printf("Gagal menghapus file dari S3: %v", err)
	}
	return err
}
