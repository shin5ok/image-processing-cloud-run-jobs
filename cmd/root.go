/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/shin5ok/image-processing-cloud-run-jobs/myimaging"
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
	Run:  runSizing,
}

func runSizing(cmd *cobra.Command, args []string) {
	// file := args[0]
	// num, err := cmd.Flags().GetInt("number")
	debug, err := cmd.Flags().GetBool("debug")
	width, err := cmd.Flags().GetInt("width")
	if err != nil {
		fmt.Println(err)
	}
	callImaging(args, width, debug)
}

func callImaging(args []string, width int, debug bool) {
	var wg sync.WaitGroup
	for _, file := range args {
		// for n := 0; n < num; n++ {
		wg.Add(1)
		go func(file string, n string) {
			defer wg.Done()
			s := myimaging.Image{Filename: file}
			if debug {
				fmt.Println("Debug:", s)
			}
			newFilename, _ := s.MakeSmall(width)
			fmt.Println("new filename is " + newFilename)
		}(file, os.Getenv("CLOUD_RUN_JOBS_INDEX"))
	}
	wg.Wait()
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
	cobra.OnInitialize(initConfig)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Int("number", 1, "Specify concurency")
	rootCmd.Flags().Bool("debug", false, "debug mode")
	rootCmd.Flags().Int("width", 240, "size of width")
}

func initConfig() {
	fmt.Println("calling init")
}
