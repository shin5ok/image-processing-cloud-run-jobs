REGION=us-central1
SUBSCRIPTION=sub
DST_BUCKET=$PROJECT-new

gcloud beta run jobs update imaging \
  --image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1 \
  --args=pubsubMode \
  --region=$REGION \
  --tasks=2 \
  --task-timeout=5m \
  --set-env-vars=DST_BUCKET=$DST_BUCKET,SUBSCRIPTION=$SUBSCRIPTION \
  --service-account=imaging-sa@$PROJECT.iam.gserviceaccount.com