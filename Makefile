include .env

login:
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${RESISTORY_URL}

build:
	docker build -t ${APP_IMAGE}:${PROD_TAG} --target prod app

test:
	cd app && go test ./...

push:
	docker push ${APP_IMAGE}:${PROD_TAG}

dev:
	docker run -it --rm -p ${PORT}:${PORT}  \
		-e GOPROXY=direct \
		--env-file $(shell pwd)/.env \
		-w /app \
		-v $(shell pwd)/app:/app \
		cosmtrek/air -c /app/.air.toml

run:
	docker run --rm -p ${PORT}:${PORT} \
		--env-file $(shell pwd)/.env \
		${APP_IMAGE}:${PROD_TAG}

apply:
	cd terraform && terraform apply -var="line_access_token=${LINE_ACCESS_TOKEN}" -auto-approve