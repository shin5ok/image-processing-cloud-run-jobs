REGION=us-central1
SUBSCRIPTION=sub
DST_BUCKET=$PROJECT-new

gcloud builds submit --pack=image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1
