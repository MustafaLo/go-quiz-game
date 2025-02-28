package cmd

import(
	"fmt"
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Quiz Game CLI",
	Long: `This is the quiz game cli from gopherexercises that I completed`,
  }
  
  func Execute() {
	if err := rootCmd.Execute(); err != nil {
	  fmt.Println(err)
	  os.Exit(1)
	}
  }

  func init(){
	rootCmd.AddCommand(quizCmd)
  }