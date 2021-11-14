deploy:
	gcloud functions deploy function-1 \
	--allow-unauthenticated \
	--trigger-http \
	--region=asia-southeast1 \
	--runtime=go113 \
	--entry-point=HelloWorld \
	--ignore-file=httptest \
	--memory=128MB \
	--max-instances=1