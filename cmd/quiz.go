package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func readCSVFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseCSV(data []byte) *csv.Reader {
	reader := csv.NewReader(bytes.NewReader(data))
	return reader
}

func getProblems(reader *csv.Reader) ([]problem, error) {
	problems := []problem{}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return problems, nil
			}
			return nil, err
		}

		entry := problem{record[0], record[1]}
		problems = append(problems, entry)
	}
}

func getUserInput(scanner *bufio.Scanner, userInput chan bool)(){
  scanner.Scan()
  userInput <- true
}


func startQuiz(problems []problem)(){
  scanner := bufio.NewScanner(os.Stdin)
  correct := 0

  
  for index, entry := range problems {
    fmt.Printf("%d. %s", index, entry.question)
    fmt.Print("\nEnter answer: ")

    userInput := make(chan bool)
    go getUserInput(scanner, userInput)

    select{
      case <- done:
        return
      case <- userInput:
        err := scanner.Err()
        if err != nil {
          fmt.Printf("Failed to read input: %v", err)
        }
    
        answer := scanner.Text()
        if answer == entry.solution {
          correct += 1
          score <- correct
        }
        fmt.Println()
    }
  }

  done <- true
}

func startTimer(){
	timer := time.NewTimer(time.Duration(limit) * time.Second)
  <-timer.C
  done <- true
  fmt.Println("Times up!")
}

type problem struct {
	question string
	solution string
}

var problems_file string
var limit int
var done chan bool
var score chan int

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Start the Quiz Game",
	Long:  `This is an exercise from gopherexercises for a Quiz Game`,
	Run: func(cmd *cobra.Command, args []string) {
		csvBytes, err := readCSVFile(problems_file)
		if err != nil {
			fmt.Printf("Failed to parse CSV file: %v", err)
			return
		}

		csvReader := parseCSV(csvBytes)
		problems, err := getProblems(csvReader)

		if err != nil {
			fmt.Printf("Failed to get problems: %v", err)
			return
		}

    done = make(chan bool)
    score = make(chan int)

    go startQuiz(problems)
    go startTimer()

    <- done

    fmt.Printf("You scored %d out of %d problems right!", <-score, len(problems))


	},
}

func init() {
	quizCmd.Flags().StringVarP(&problems_file, "problems", "p", "problems.csv", "Problems")
	quizCmd.Flags().IntVarP(&limit, "limit", "l", 2, "Timer for the quiz game")
}
