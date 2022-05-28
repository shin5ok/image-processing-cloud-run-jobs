TOPIC=gcs-event

gcloud pubsub create topic $TOPIC
gsutil mb gs://$PROJECT-src
gsutil mb gs://$PROJECT-new
gsutil notification create -f json -t $TOPIC gs://$PROJECT-src