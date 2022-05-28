REGION=us-central1
SRC_BUCKET=$PROJECT-src

gcloud iam service-accounts create imaging-sa

gcloud projects add-iam-policy-binding $PROJECT   --role roles/storage.admin --member serviceAccount:imaging-sa@$PROJECT.iam.gserviceaccount.com

gcloud beta run jobs create imaging \
  --image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1 \
  --args=pubsub_mode \
  --tasks=2 \
  --task-timeout=5m \
  --set-env-vars=SRC_BUCKET=$SRC_BUCKET \
  --service-account=imaging-sa@$PROJECT.iam.gserviceaccount.com
