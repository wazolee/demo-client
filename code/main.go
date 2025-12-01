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
	fmt.Println("Menu:")
	fmt.Println("1 - Add score")
	fmt.Println("2 - List scores")
	fmt.Println("q - Quit")
	for shouldContinue {
		var option string = "0"
		fmt.Scanln(&option)
		switch option{
		case "1":
			
			fmt.Println("Enter a name and score: ")
			var name, rawScore string
			fmt.Scanln(&name, &rawScore)
			s, _ := strconv.Atoi(rawScore)
			scores = append(scores, score{name: name, score: s})
			fmt.Printf("Name: %s, Score: %d\n", name, s)		
		case "2":
			fmt.Println("Student scores")
			fmt.Println(strings.Repeat("-", 14))
			for _, sc := range scores {
				fmt.Println(sc.name, ":", sc.score)
			}
		case "q":
			shouldContinue = false
		}

	}
}