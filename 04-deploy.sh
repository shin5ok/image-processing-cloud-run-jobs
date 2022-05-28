REGION=us-central1
SUBSCRIPTION=sub
DST_BUCKET=$PROJECT-new

MODE=update # or create
gcloud beta run jobs $MODE imaging \
  --image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1 \
  --args=pubsubMode \
  --region=$REGION \
  --tasks=2 \
  --task-timeout=5m \
  --set-env-vars=PROJECT=$PROJECT,DST_BUCKET=$DST_BUCKET,SUBSCRIPTION=$SUBSCRIPTION \
  --service-account=imaging-sa@$PROJECT.iam.gserviceaccount.com