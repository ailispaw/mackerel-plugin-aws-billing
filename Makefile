NAME    := mackerel-plugin-aws-billing
IMAGE   := ailispaw/$(NAME)
VERSION := 0.3.0

build:
	docker build -t $(IMAGE) src
	docker tag $(IMAGE) $(IMAGE):builder

run:
	-docker rm -f $(NAME)
	docker run --name $(NAME) --env-file .env $(IMAGE)

release:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) release
	docker tag $(IMAGE) $(IMAGE):$(VERSION)

push:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

clean:
	-docker rm $$(docker ps -a -q)
	-docker rmi $$(docker images -q -f "dangling=true")

.PHONY: build run release push clean
