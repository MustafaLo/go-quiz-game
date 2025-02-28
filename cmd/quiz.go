package cmd

import (
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

type problem struct{
  problem string
  solution int
}

var problems_file string

var quizCmd = &cobra.Command{
  Use:   "quiz",
  Short: "Start the Quiz Game",
  Long:  `This is an exercise from gopherexercises for a Quiz Game`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Loading Quiz...")
    
    csvBytes, err := readCSVFile(problems_file)
    if err != nil{
      fmt.Println("Failed to parse CSV file: %v", err)
    }

    csvReader := parseCSV(csvBytes)

    for {
      record, err := csvReader.Read()
      if err != nil{
        if err == io.EOF{
          break
        }
        fmt.Println("Error while reading CSV file: %v", err)
        return
      }

      fmt.Println(record)
    }
  },
}

func init(){
  quizCmd.Flags().StringVarP(&problems_file, "problems", "p", "problems.csv", "Problems")
}

