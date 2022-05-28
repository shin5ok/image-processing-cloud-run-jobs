REGION=us-central1
gcloud beta run jobs list
gcloud beta run jobs execute imaging --region $REGION
gcloud beta run jobs executions list --job imaging --region $REGION