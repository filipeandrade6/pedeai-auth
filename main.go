package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// The input type and the output type are defined by the API Gateway.
func handleRequest2(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := sqlx.Connect("postgres", "user=pedeai dbname=main sslmode=disable password=mypassword host=myhost")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	db.Queryx("SELECT * FROM clientes")

	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
		Body:       fmt.Sprintf("Hello!\n"),
	}
	return res, nil
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Print("entrou na lambda!")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
	defer cancel()

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	var usuario = os.Getenv("DB_USER")
	var senha = os.Getenv("DB_PASS")
	var host = os.Getenv("DB_HOST")
	var port = os.Getenv("DB_PORT")
	var db = os.Getenv("DB_NAME")

	var url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", usuario, senha, host, port, db)

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	log.Print("conectou!")

	var id string

	err = conn.QueryRow(context.Background(), "select id from clientes where nome=$1", "FILIPE ANDRADE").Scan(&id)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}

	log.Print("executou a query!")

	fmt.Println(id)

	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain; charset=utf-8"},
		Body:       fmt.Sprintf(id),
	}

	return res, nil
}

func main() {
	lambda.Start(handleRequest)
}
