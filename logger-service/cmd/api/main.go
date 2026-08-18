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

	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
)

// Đổi localhost thành 127.0.0.1 để chạy ổn định trên Windows local
var mongoURL = getEnv("MONGO_URL", "mongodb://admin:password@127.0.0.1:27017/?authSource=admin")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// 1. Connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// 2. Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 3. Close connection when main exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// 4. Set up config
	app := Config{
		Models: data.New(client),
	}

	// Register the RPC server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	// gRPC
	go app.gRPCListen()

	// 5. Start the web server
	go app.serve()

	select {}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("Starting logger service on port %s\n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// Chỉ dùng ApplyURI vì chuỗi URI đã bao gồm sẵn User, Password và authSource=admin
	clientOptions := options.Client().ApplyURI(mongoURL)

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	// Ping the primary to verify connection
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Error pinging mongo:", err)
		return nil, err
	}

	log.Println("Connected to Mongo!")
	return c, nil
}
