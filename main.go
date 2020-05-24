package main

import (
	_ "database/sql"
	"fmt"
	_ "fmt"
	"github.com/UstinovV/wm_api/auth"
	"github.com/UstinovV/wm_api/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func main() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("Config open error: ",err)
	}
	defer f.Close()

	config := server.NewConfig()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Config error: ",err)
	}
	fmt.Println(config.DBConnection)
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	serv := server.NewServer(config, sugar)

	err = serv.Start()
	if err != nil {
		log.Fatal("Server error: ", err)
	}

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cant connect to auth server ", err)
	}
	defer conn.Close()
	fmt.Printf("ASD")
	authClient := auth.NewAuthCheckerClient(conn)

	resp, err := authClient.Test(context.Background(), &auth.Message{Body:"Hello from main "})
	if err != nil {
		log.Fatal("Error when hello ", err)
	}

	fmt.Printf("zxc")
	log.Printf("Response from grpc %s", resp.Body)
}
