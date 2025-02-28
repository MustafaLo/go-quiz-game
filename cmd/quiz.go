package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func readCSVFile(filename string) ([]byte, error){
  f, err := os.Open(filename)
  if err != nil{
    return nil, err
  }

  defer f.Close()

  data, err := io.ReadAll(f)
  if err != nil{
    return nil, err
  }

  return data, nil
}

func parseCSV(data []byte)(*csv.Reader){
  reader:= csv.NewReader(bytes.NewReader(data))
  return reader
}

func getProblems(reader *csv.Reader)([]problem, error){
  problems := []problem{}

  for {
    record, err := reader.Read()
    if err != nil{
      if err == io.EOF{
        return problems, nil
      }
      return nil, err
    }

    entry := problem{record[0], record[1]}
    problems = append(problems, entry)
  }
}

type problem struct{
  question string
  solution string
}

var problems_file string
var timer int

var quizCmd = &cobra.Command{
  Use:   "quiz",
  Short: "Start the Quiz Game",
  Long:  `This is an exercise from gopherexercises for a Quiz Game`,
  Run: func(cmd *cobra.Command, args []string) {
    csvBytes, err := readCSVFile(problems_file)
    if err != nil{
      fmt.Printf("Failed to parse CSV file: %v", err)
      return
    }

    csvReader := parseCSV(csvBytes)
    problems, err := getProblems(csvReader)
    if err != nil{
      fmt.Printf("Failed to get problems: %v", err)
      return
    }

    score := 0
    scanner := bufio.NewScanner(os.Stdin)
    for index, entry := range problems{
      fmt.Printf("%d. %s", index, entry.question)
      fmt.Print("\nEnter answer: ")

      scanner.Scan()
      err := scanner.Err()
      if err != nil{
        fmt.Printf("Failed to read input: %v", err)
        return
      }

      answer := scanner.Text()
      if answer == entry.solution{
        score += 1
      }
      fmt.Println()
    }

    fmt.Printf("You scored %d out of %d problems right!", score, len(problems))


  },
}

func init(){
  quizCmd.Flags().StringVarP(&problems_file, "problems", "p", "problems.csv", "Problems")
  quizCmd.Flags().IntVarP(&timer, "timer", "t", 30, "Timer for the quiz game")
}

