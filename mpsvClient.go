package main

import (
	"context"
	"fmt"
	"github.com/UstinovV/wm_api/database"
	"github.com/UstinovV/wm_api/mpsv"
	"github.com/UstinovV/wm_api/server"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("Config open error: ", err)
	}
	defer f.Close()

	config := server.NewConfig()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Config error: ", err)
	}
	fmt.Println(config.DBConnection)
	db := database.NewDB(config.DBConnection)

	err = db.Open()
	if err != nil {
		log.Fatal("DB err :", err)
	}
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":8011", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cant connect to auth server ", err)
	}
	defer conn.Close()
	mpsvClient := mpsv.NewMpsvParserClient(conn)

	mpsvUrl := mpsv.MpsvUrl{Url: "http://some.url"}
	stream, err := mpsvClient.ParseMpsvUrl(context.Background(), &mpsvUrl)
	if err != nil {
		log.Fatal("Stream error :", err)
	}
	stmt, err := db.Db.Prepare("INSERT INTO temp_offers (title, content, external_id) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal("Error preparing statement: ", err)
	}
	for {
		o, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", mpsvClient, err)
		}

		_, err = stmt.Exec(o.Title, o.Content, o.MpsvId)
		if err != nil {
			fmt.Println("Insert error : ", err)
		}
	}
}
