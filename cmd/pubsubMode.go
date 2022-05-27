/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/shin5ok/image-processing-cloud-run-jobs/myimaging"
	"github.com/spf13/cobra"
)

// pubsubModeCmd represents the pubsubMode command
var pubsubModeCmd = &cobra.Command{
	Use:   "pubsubMode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		topic, _ := cmd.Flags().GetString("sub")
		timeout, _ := cmd.Flags().GetInt("timeout")
		project := os.Getenv("GOOGLE_CLOUD_PROJECT")
		pullMsgsSync(project, topic, timeout)
	},
}

func init() {
	rootCmd.AddCommand(pubsubModeCmd)
	pubsubModeCmd.Flags().String("sub", "", "")
	pubsubModeCmd.Flags().Int("timeout", 10, "")

}

func pullMsgsSync(projectID, subID string, timeout int) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	/*
			{
		  "kind": "storage#object",
		  "id": "shingogcp-firestore-nativemode/images/images/dog/new_000002.png/1653668760085964",
		  "selfLink": "https://www.googleapis.com/storage/v1/b/shingogcp-firestore-nativemode/o/images%2Fimages%2Fdog%2Fnew_000002.png",
		  "name": "images/images/dog/new_000002.png",
		  "bucket": "shingogcp-firestore-nativemode",
		  "generation": "1653668760085964",
		  "metageneration": "1",
		  "contentType": "image/png",
		  "timeCreated": "2022-05-27T16:26:00.091Z",
		  "updated": "2022-05-27T16:26:00.091Z",
		  "storageClass": "STANDARD",
		  "timeStorageClassUpdated": "2022-05-27T16:26:00.091Z",
		  "size": "69099",
		  "md5Hash": "7tcjvcVgCUzJbdTRKUhfPg==",
		  "mediaLink": "https://www.googleapis.com/download/storage/v1/b/shingogcp-firestore-nativemode/o/images%2Fimages%2Fdog%2Fnew_000002.png?generation=1653668760085964&alt=media",
		  "crc32c": "N5z6iw==",
		  "etag": "CMzLlZiMgPgCEAE="
		}
	*/

	type dataStruct struct {
		Name   string `json:"name"`
		Bucket string `json:"bucket"`
	}

	var received int32
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		var datastruct dataStruct

		fmt.Printf("%+v\n", string(msg.Data))
		json.Unmarshal(msg.Data, &datastruct)
		fmt.Printf("%+v\n", datastruct)
		filePath := fmt.Sprintf("%s/%s", datastruct.Bucket, datastruct.Name)
		fmt.Println("gs://" + filePath)

		go processingImage(datastruct.Bucket, datastruct.Name)

		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}
	fmt.Printf("Received %d messages\n", received)

	return nil
}

func processingImage(bucket, object string) {
	tmpFile := uuid.NewString() + ".jpg"
	downloadFile(bucket, object, tmpFile)
	s := myimaging.Image{Filename: tmpFile}
	newFilename, _ := s.MakeSmall(240)
	uploadFile(bucket, newFilename)
	os.Remove(tmpFile)
	os.Remove(newFilename)
}

func uploadFile(bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("notes.txt")
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Printf("Blob %v uploaded.\n", object)
	return nil
}

func downloadFile(bucket, object string, destFileName string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// destFileName := "file.txt"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	f, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %v", err)
	}

	fmt.Printf("Blob %v downloaded to local file %v\n", object, destFileName)

	return nil

}
