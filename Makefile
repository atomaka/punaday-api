yarnbin = $(shell yarn bin)

deploy:
	sh -c "source .env.production; $(yarnbin)/serverless deploy --stage production"
