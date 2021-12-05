package main

import (
	"aws-go-cdk-wkhtmltopdf-lambda/api/documents"

	"github.com/akrylysov/algnhsa"

	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	log.Info("starting handler")
	h := documents.NewHandler(log)
	algnhsa.ListenAndServe(h, &algnhsa.Options{
		RequestType:        algnhsa.RequestTypeAPIGateway,
		BinaryContentTypes: []string{"application/pdf"},
	})
}
