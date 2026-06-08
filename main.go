package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type GameState struct {
	BallX    float64
	BallY    float64
	BallDX   float64
	BallDY   float64
	Paddle1Y float64
	Paddle2Y float64
	Score1   int
	Score2   int
	Winner   int // 0: No winner, 1: Player 1, 2: Player 2
	mu       sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	game := &GameState{
		BallX: 400, BallY: 200,
		BallDX: 4, BallDY: 4,
		Paddle1Y: 150, Paddle2Y: 150,
		Winner:   0,
	}

	// Physics Loop
	go func() {
		for {
			updatePhysics(game)
			time.Sleep(16 * time.Millisecond)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		handleConnection(conn, game)
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updatePhysics(g *GameState) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Stop physics if someone has already won
	if g.Winner != 0 {
		return
	}

	g.BallX += g.BallDX
	g.BallY += g.BallDY

	if g.BallY <= 0 || g.BallY >= 390 {
		g.BallDY *= -1
	}

	if g.BallX <= 30 && g.BallY > g.Paddle1Y && g.BallY < g.Paddle1Y+100 {
		g.BallDX = 4
	}

	if g.BallX >= 760 && g.BallY > g.Paddle2Y && g.BallY < g.Paddle2Y+100 {
		g.BallDX = -4
	}

	// Updated Scoring Logic with Win Condition
	if g.BallX < 0 {
		g.Score2++
		if g.Score2 >= 10 {
			g.Winner = 2
		} else {
			resetBall(g)
		}
	} else if g.BallX > 800 {
		g.Score1++
		if g.Score1 >= 10 {
			g.Winner = 1
		} else {
			resetBall(g)
		}
	}
}

func resetBall(g *GameState) {
	g.BallX = 400
	g.BallY = 200
	g.BallDX *= -1
}

func handleConnection(conn *websocket.Conn, g *GameState) {
	defer conn.Close()
	for {
		g.mu.Lock()
		err := conn.WriteJSON(g)
		g.mu.Unlock()
		if err != nil {
			break
		}

		var input struct {
			Y float64 `json:"y"`
		}
		err = conn.ReadJSON(&input)
		if err != nil {
			break
		}

		g.mu.Lock()
		// Only allow paddle movement if game is active
		if g.Winner == 0 {
			g.Paddle1Y = input.Y
		}
		g.mu.Unlock()

		time.Sleep(16 * time.Millisecond)
	}
}