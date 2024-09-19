package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
)

type User struct {
	CPF string `json:"cpf"`
}

func handleRequest(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	if req.RequestContext.HTTP.Method != "POST" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       http.StatusText(http.StatusMethodNotAllowed),
		}, nil
	}

	var user User
	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
	defer cancel()

	var usuario = os.Getenv("DB_USER")
	var senha = os.Getenv("DB_PASS")
	var host = os.Getenv("DB_HOST")
	var port = os.Getenv("DB_PORT")
	var db = os.Getenv("DB_NAME")

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	var url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", usuario, senha, host, port, db)

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	var id string
	err = conn.QueryRow(context.Background(), "SELECT id FROM clientes WHERE cpf=$1", user.CPF).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusNotFound,
				Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
				Body:       "CPF não encontrado",
			}, nil
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
		Body:       id,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
