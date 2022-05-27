/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

// pubsubModeCmd represents the pubsubMode command
var pubsubModeCmd = &cobra.Command{
	Use:   "pubsubMode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		topic, _ := cmd.Flags().GetString("topic")
		project := os.Getenv("GOOGLE_CLOUD_PROJECT")
		pullMsgsSync(project, topic)
	},
}

func init() {
	rootCmd.AddCommand(pubsubModeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	pubsubModeCmd.Flags().String("topic", "", "")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pubsubModeCmd.Flags().Get("toggle", "t", false, "Help message for toggle")
}

func pullMsgsSync(projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
	// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
	// the max number of messages the client will hold in memory at a time.
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	fmt.Printf("Received %d messages\n", received)

	return nil
}
