package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan message2)
var k string
var v string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type message2 struct {
	RollNo    string `json:"rollNo"`
	Message   string `json:"message"`
	TotalWord string `json:"totalWord"`
	TotalChar string `json:"totalChar"`
	Wpm       string `json:"wpm"`
}

var templates *template.Template

func main() {
	initRedis()
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true
	for {
		var msg message2

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
		set(msg.RollNo, msg.Wpm)
		fmt.Println("Word Per Minute data picked from Redis Database")
		fmt.Println(get(msg.RollNo))
	}
}

func handleMessages() {
	for {

		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

var pool *redis.Pool

func initRedis() {

	initPool()

	initStore()
}
func initPool() {
	pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Printf("ERROR: fail init redis: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func initStore() {

	conn := pool.Get()
	defer conn.Close()
}

func set(key string, val string) error {

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
		return err
	}

	return nil
}

func get(key string) (string, error) {

	conn := pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}

	return s, nil
}
