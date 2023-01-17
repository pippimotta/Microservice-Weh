package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct {}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo (p RPCPayload, rsp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(),data.LogEntry{
		Name: p.Name,
		Data: p.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}
	*rsp = "Processed payload via RPC: "+ p.Name
	return nil
}