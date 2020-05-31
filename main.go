package main

import (
	"context"
	_ "database/sql"
	"fmt"
	_ "fmt"
	"github.com/UstinovV/wm_api/mpsv"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	//f, err := os.Open("config.yml")
	//if err != nil {
	//	log.Fatal("Config open error: ",err)
	//}
	//defer f.Close()
	//
	//config := server.NewConfig()
	//decoder := yaml.NewDecoder(f)
	//err = decoder.Decode(&config)
	//if err != nil {
	//	log.Fatal("Config error: ",err)
	//}
	//fmt.Println(config.DBConnection)
	//logger, _ := zap.NewProduction()
	//defer logger.Sync() // flushes buffer, if any
	//sugar := logger.Sugar()
	//
	//serv := server.NewServer(config, sugar)
	//
	//err = serv.Start()
	//if err != nil {
	//	log.Fatal("Server error: ", err)
	//}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cant connect to auth server ", err)
	}
	defer conn.Close()
	fmt.Printf("ASD")
	mpsvClient := mpsv.NewMpsvParserClient(conn)
	//authClient := auth.NewAuthCheckerClient(conn)
	//
	//resp, err := authClient.Test(context.Background(), &auth.Message{Body:"Hello from main "})
	//if err != nil {
	//	log.Fatal("Error when hello ", err)
	//}
	//
	//log.Printf("Response from grpc %s", resp.Body)

	//mpsvOffer := mpsv.MpsvOffer{}
	mpsvUrl := mpsv.MpsvUrl{Url:"http://some.url"}
	stream, err := mpsvClient.ParseMpsvUrl(context.Background(), &mpsvUrl)
	if err != nil {
		log.Fatal("Stream error :", err)
	}
	for {
		o, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", mpsvClient, err)
		}
		log.Println(o.Title)
	}
}
