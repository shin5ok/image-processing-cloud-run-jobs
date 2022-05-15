/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "image-processing-cloud-run-jobs",
	Short: "Image Processing",
	Long: `This is an Image Processing program.
You can run it on local environment, and also on Cloud Run jobs
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		num, err := cmd.Flags().GetInt("number")
		if err != nil {
			fmt.Println(err)
		}
		var wg sync.WaitGroup
		for n := 0; n < num; n++ {
			wg.Add(1)
			go func(file string, n int) {
				defer wg.Done()
				fmt.Printf("%d: %s\n", n+1, file)
			}(file, n)
		}
		wg.Wait()
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.image-processing-cloud-run-jobs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Int("number", 1, "Specify concurency")
}
