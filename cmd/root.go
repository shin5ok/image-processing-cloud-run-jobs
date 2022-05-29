/*
Copyright © 2022 Shingo <shin5ok@55mp.com>

*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/shin5ok/image-processing-cloud-run-jobs/myimaging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "image-processing-cloud-run-jobs",
	Short: "Image Processing",
	Long: `This is an Image Processing program.
You can run it on local environment, and also on Cloud Run jobs
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  runSizing,
}

func runSizing(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	width, err := cmd.Flags().GetInt("width")
	if err != nil {
		fmt.Println(err)
	}
	callImaging(args, width, debug)
}

func callImaging(args []string, width int, debug bool) {
	var wg sync.WaitGroup
	for _, file := range args {
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().Int("number", 1, "Specify concurency")
	rootCmd.Flags().Bool("debug", false, "debug mode")
	rootCmd.Flags().Int("width", 240, "size of width")
}

func initConfig() {
	fmt.Println("calling init")
}
