package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	//connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register rpcServer
	err = rpc.Register(new(RPCServer))
	if  err != nil {
		log.Panic()
	}
	go app.rpcListen()
  go app.gRPCListen()



	// start web server

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	log.Println("Starting Service on port", webPort)
	if err != nil {
		log.Panic()
	}
}
func (app *Config)rpcListen() error{

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	log.Printf("Starting RPC Server on port %s", rpcPort)
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(conn)
	}

}
func connectToMongo() (*mongo.Client, error) {
	//create connect options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting", err)
		return nil, err
	}
	log.Println("Mongo service starts!")
	return c, nil

}
