.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fabric && docker-compose up --force-recreate -d
	@echo "Sleep 15 seconds in order to let the environment setup correctly"
	@sleep 15
	@echo "Environment up"

env-up-pro:
	@echo "Start environment ..."
	@cd fabric && docker-compose -f docker-compose-pro.yml up --force-recreate -d
	@echo "Sleep 15 seconds in order to let the environment setup correctly"
	@sleep 15
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fabric && docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./heroes-service

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/heroes-service-* heroes-service
	@docker rm -f -v `docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"

fe-app:
	docker rm -f health-service-fe-app || true
	docker pull registry.cn-beijing.aliyuncs.com/cabernety/health-service-fe-app:latest
	docker run --name health-service-fe-app -d -p 3001:3000 registry.cn-beijing.aliyuncs.com/cabernety/health-service-fe-app:latest

fe:
	docker rm -f health-service-fe || true
	docker pull registry.cn-beijing.aliyuncs.com/cabernety/health-service-fe:latest
	docker run -name health-service-fe -d -p 3000:3000 registry.cn-beijing.aliyuncs.com/cabernety/health-service-fe:latest

