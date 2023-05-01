include .env

build:
	docker build -t ${APP_IMAGE}:${PROD_TAG} --target prod app

build-dev:
	docker build -t ${APP_IMAGE}:${DEV_TAG} --target dev app

push:
	docker push ${APP_IMAGE}:${PROD_TAG}

dev:
	docker run -it --rm -p ${PORT}:${PORT}  \
		-e GOPROXY=direct \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
		--env-file $(shell pwd)/.env \
		-w /app \
		-v $(shell pwd)/app:/app \
		cosmtrek/air -c /app/.air.toml

run:
	docker run --rm -p ${PORT}:${PORT} \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
		--env-file $(shell pwd)/.env \
		${APP_IMAGE}:${PROD_TAG}
