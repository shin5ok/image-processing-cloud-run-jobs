REGION=us-central1
DST_BUCKET=$PROJECT-new

gcloud builds submit --pack=image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1

gcloud beta run jobs update imaging \
  --image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1 \
  --args=pubsubMode \
  --region=$REGION \
  --tasks=2 \
  --task-timeout=5m \
  --set-env-vars=DST_BUCKET=$DST_BUCKET \
  --service-account=imaging-sa@$PROJECT.iam.gserviceaccount.com
