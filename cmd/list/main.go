package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "SNS/SQS test program",
	Run:   runRoot,
}

var cmdSns = &cobra.Command{
	Use:   "sns",
	Short: "List SNS queues",
	Args:  cobra.ExactArgs(0),
	Run:   runSns,
}

var cmdSqs = &cobra.Command{
	Use:   "sqs",
	Short: "List SQS queues",
	Args:  cobra.ExactArgs(0),
	Run:   runSqs,
}

func init() {
	rootCmd.AddCommand(cmdSns)
	rootCmd.AddCommand(cmdSqs)
}

func main() {
	// run root command
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("root command error")
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Usage()
}
