package main

import (
	"fmt"
	"strconv"
	"strings"

	// "time"
)

type score struct {
	name  string
	score int
}

func main() {

	
	scores := []score{
		{name: "Dent, Arthur", score: 87},
		{name: "Beeblebrox, Zaphod", score: 42},
		{name: "Prefect, Ford", score: 100},
	}

	loopDemo(scores)
	return

	fmt.Println("Select score to print (1 - 3):")
	var option string
	fmt.Scanln(&option)

	fmt.Println("Student scores")
	fmt.Println(strings.Repeat("-", 14))

	var index int

	switch option {
		case "1":		
			index = 0
		case "2":
			index = 1
		case "3":
			index = 2
		default:
			fmt.Println("Invalid option")
			return
	}
	fmt.Println(scores[index].name, ":", scores[index].score)
	// time.Sleep(10 * time.Second)

	
	// time.Sleep(10 * time.Second)

	
}

func loopDemo(scores []score) {
	shouldContinue := true
	
	for shouldContinue {
		printMenu()
		var option string = "0"
		fmt.Scanln(&option)
		switch option{
		case "1":			
			scores = append(scores, addStudent())
		case "2":
			printReport(scores)
		case "q":
			shouldContinue = false
		}

	}
}

func addStudent() score {
	fmt.Println("Enter a name and score: ")
	var name, rawScore string
	fmt.Scanln(&name, &rawScore)
	s, _ := strconv.Atoi(rawScore)	
	fmt.Printf("Name: %s, Score: %d\n", name, s)
	return score{name: name, score: s}
}

func printReport(scores []score) {
	fmt.Println("Student scores")
	fmt.Println(strings.Repeat("-", 14))
	for _, sc := range scores {
		fmt.Println(sc.name, ":", sc.score)
	}
}

func printMenu(){
	fmt.Println("Menu:")
	fmt.Println("1 - Add score")
	fmt.Println("2 - List scores")
	fmt.Println("q - Quit")
}