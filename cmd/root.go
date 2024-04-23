/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gq",
	Short: "A CLI to ask questions about the data",
	Long: `This CLI is used to ask questions about the data you send to it. 
  It will read the data from stdin or pipe and ask questions from the user. 
  The output will be written to stdout.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand(cmd, args)
	},
}

/**
* This is the main method where we control the routing of the command.
 */
func runCommand(cmd *cobra.Command, args []string) error {
	question, _ := cmd.Flags().GetString("question")

	if question == "" {
		return fmt.Errorf("no question provided")
	}

	if isInputFromPipe() {
		readFromPipe(question, os.Stdin, os.Stdout)
	} else {
		if len(args) == 0 {
			return fmt.Errorf("no data provided")
		}
		result := askQuestion(question, args[0])
		write(result, os.Stdout)
	}
	return nil
}

/*
* This function checks if the input is from pipe
 */
func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

/**
* This function asks a question to the provider and returns the answer
 */
func askQuestion(question string, data string) string {
	return strings.ToUpper(question) + strings.ToUpper(data)
}

/**
* This function reads the data from the pipe and asks questions and write it as output
 */
func readFromPipe(question string, reader io.Reader, writer io.Writer) error {
	fmt.Println(question)
	scanner := bufio.NewScanner(bufio.NewReader(reader))
	for scanner.Scan() {
		result := askQuestion(question, scanner.Text())
		write(result, writer)
	}
	return nil
}

/**
* This function writes the result as output
 */
func write(s string, w io.Writer) error {
	_, e := fmt.Fprintln(w, strings.ToUpper(s))
	if e != nil {
		return e
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "x", false, "debug mode")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.gq.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("question", "q", "", "Question about the data sent")
}