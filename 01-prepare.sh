TOPIC=gcs-event

gcloud pubsub create topic $TOPIC
gsutil mb gs://$PROJECT-src
gsutil mb gs://$PROJECT-new
gsutil notification create -f json -t $TOPIC gs://$PROJECT-src
gcloud artifacts repositories create containers --repository-format=docker --location=$REGION

if ! gcloud iam service-accounts list | grep imaging-sa > /dev/null;
then
  gcloud iam service-accounts create imaging-sa
  gcloud projects add-iam-policy-binding $PROJECT   --role roles/storage.admin --member serviceAccount:imaging-sa@$PROJECT.iam.gserviceaccount.com
fi
