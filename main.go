package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	color "github.com/fatih/color"
)

type TickMsg time.Time

const WIDTH = 20
const HEIGHT = 16

type model struct {
	gameState 		 string 
	board     		 [][]string
	snakeCells 		 [][]int  
	score      		 int     
	fruitCell 		 []int 
	currentDirection string
	newDirection 	 string
}

func initialModel() model {
	gameBoard := make([][]string, HEIGHT)
	for i := 0; i < HEIGHT; i ++ {
		gameBoard[i] = make([]string, WIDTH) 
		for j := 0; j < WIDTH; j++ {
			 gameBoard[i][j] = "[]"
		}
	}	

	return model {
		gameState: "starting",
		board: gameBoard,
		snakeCells : [][]int{{HEIGHT/2, (WIDTH/2)-5}, {HEIGHT/2,(WIDTH/2)-4}, {HEIGHT/2,(WIDTH/2)-3}},
		score: 0,
		fruitCell: []int{HEIGHT/2,WIDTH/2 + 6},
		currentDirection: "right",
		newDirection: "right",
	}
}

func (m model) Init() tea.Cmd {
	return tickEvery()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// prevent snake from turning into itself
		case "up":
			if m.currentDirection != "down"{
				m.newDirection= "up"
			}
		case "down":
			if m.currentDirection != "up"{
				m.newDirection = "down"
			}
		case "left":
			if m.currentDirection != "right"{
				m.newDirection = "left"
			}

		case "right":
			if m.currentDirection != "left"{
				m.newDirection= "right"
			}
		case " ":
			m.gameState = "playing"
		}

	case TickMsg:
		if m.gameState == "playing" {
			boolean :=  false
			m, boolean = updateSnake(m)
			if boolean {
				return m, tea.Quit
			}
			m.currentDirection = m.newDirection
		}
		return m, tickEvery()
	}

	return m, nil
}

func (m model) View() string {
	message := ""	
	switch m.gameState {
	case "starting":
		message = "Press SPACE to begin!"
	case "ended":
		message = "Game Over!"
	}

	state := fmt.Sprintf("\n     Score: %d %s\n",m.score, message)
	greenBackground := color.New(color.BgHiGreen)
	redBackground := color.New(color.BgHiRed)

	for i:=0;i <HEIGHT; i++{
		state += "     "
		for j:=0;j <WIDTH; j++{
			if sliceContains(m.snakeCells, []int{i,j}) {
				state += greenBackground.Sprintf("[]") 

			} else if reflect.DeepEqual(m.fruitCell, []int{i,j}) {
				state += redBackground.Sprintf("[]")
			} else {
				state += m.board[i][j]
			}
		}
		state += "\n"
	}
	state += "\n     Press 'q' or Ctrl+C to quit"
	return state
}

func sliceContains(slice [][]int, subslice []int ) (bool) {
	for _, element := range slice{
		if reflect.DeepEqual(element, subslice) {
			return true	
		}
	}
	return false
}

func updateSnake(m model) (model, bool) {
	head := m.snakeCells[0]
	
	switch m.currentDirection{
	case "up":    head = []int{m.snakeCells[0][0] - 1, m.snakeCells[0][1]}
	case "down":  head = []int{m.snakeCells[0][0] + 1, m.snakeCells[0][1]}
	case "right": head = []int{m.snakeCells[0][0], m.snakeCells[0][1] + 1}
	case "left":  head = []int{m.snakeCells[0][0], m.snakeCells[0][1] - 1}
	}

	// calculate new head location based on currentDirection
	newSnakeCells := [][]int{}
	ateFruit := false
	if reflect.DeepEqual(head, m.fruitCell) {
		ateFruit = true
	}

	if (ateFruit) {
		newSnakeCells = append(m.snakeCells, head)
		m.score++
		m.fruitCell = []int{rand.Intn(HEIGHT-1), rand.Intn(WIDTH-1)}
		
	} else {
		newSnakeCells = [][]int{head}
		m.snakeCells = m.snakeCells[:len(m.snakeCells)-1]	
		for _, element := range m.snakeCells {
			newSnakeCells = append(newSnakeCells, element)
		}
	}

	// impossible to crash into first 3 cells so just ignore those ones 
	if sliceContains(m.snakeCells[2:], head) {
		m.gameState = "ended"
		return m, true
	}

	if head[0] > HEIGHT-1 || head[0] < 0 {
		m.gameState = "ended"
		return m, true
	} else if head[1] > WIDTH-1|| head[1] < 0 {
		m.gameState = "ended"
		return m, true
	}

	m.snakeCells = newSnakeCells
	return m, false
}


func tickEvery() tea.Cmd {
	return tea.Every(time.Second / 10, func(t time.Time) tea.Msg {
		return TickMsg(t)	
	})
}

func main() {
	p := tea.NewProgram((initialModel()))

	if err := p.Start(); err != nil {
		fmt.Printf("Error starting tea program ): %v", err )
		os.Exit(1)
	}
}