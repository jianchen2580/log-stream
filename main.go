package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"os"

	"github.com/Shopify/sarama"
	//"io/ioutil"
	"net/http"
)

const (
	prismPattern string = `(?P<%s>[0-9]+)/(?P<appId>[0-9]+)`
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	// TODO: read from paramter
	fmt.Println(c.Param("accountID"))
	fmt.Printf(prismPattern, c.Param("accountID"))
	brokers := []string{"127.0.0.1:9092"}
	topic := "test"
	consumer, err := NewKafkaConsumer(topic, brokers)
	if err != nil {
		panic(err)
	}
	defer consumer.close()
	// TODO: Offset start point
	cp, err := consumer.consumer.ConsumePartition(consumer.topic, 0, sarama.OffsetOldest)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case msg := <-cp.Messages():
			message := Message{}
			err := json.Unmarshal([]byte(msg.Value), &message)
			if err != nil {
				// TODO: handle unmarshal err here
				fmt.Println("unmarshal failed")
			}
			// TODO: debug print
			fmt.Println(message)

			err = conn.WriteMessage(websocket.TextMessage, []byte(msg.Value))
			if err != nil {
				panic(err)
			}

		}
	}
}

func index(c *gin.Context) {
	accountID := c.Query("account_id")
	c.HTML(200, "index.tmpl", gin.H{"account_id": accountID})
}

func Run() error {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//r.Static("/static", "./static")
	//r.Static("/assert", "./assert")
	r.GET("/", index)
	r.GET("/ws/:accountID", func(c *gin.Context) { wshandler(c.Writer, c.Request, c) })
	r.Run(":8080")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "es-puller"
	app.Usage = "work with `es-puller` microservice"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "broker, b", Value: "127.0.0.1:9092", Usage: "kafka broker", EnvVar: "KAFKA_BROKER"},
		cli.StringFlag{Name: "topic, t", Value: "test", Usage: "kafka topic", EnvVar: "KAFKA_TOPIC"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				if err := Run(); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
