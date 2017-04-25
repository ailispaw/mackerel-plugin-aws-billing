NAME    := mackerel-plugin-aws-billing
IMAGE   := ailispaw/$(NAME)
VERSION := 0.2.0

run:
	docker rm -f $(NAME)
	docker run --name $(NAME) --env-file .env $(IMAGE)

build:
	docker build -t $(IMAGE) src

release:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) release
	docker tag $(IMAGE) $(IMAGE):$(VERSION)

push:
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

clean:
	docker rm $$(docker ps -q -f "exited!=0")
	docker rmi $$(docker images -q -f "dangling=true")

.PHONY: run build release push clean
