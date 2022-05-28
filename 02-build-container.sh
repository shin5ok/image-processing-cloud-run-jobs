REGION=us-central1
gcloud artifacts repositories create containers --repository-format=docker --location=$REGION
gcloud builds submit --pack=image=$REGION-docker.pkg.dev/$PROJECT/containers/imaging:v1