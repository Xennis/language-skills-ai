include ../../Makefile.variables

deploy:
	gcloud --project $(GCP_PROJECT) functions deploy CorrectionAI --region $(GCP_REGION) --runtime go116 --trigger-http --allow-unauthenticated
