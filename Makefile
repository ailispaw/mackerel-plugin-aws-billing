NAME    := mackerel-plugin-aws-billing
IMAGE   := ailispaw/$(NAME)
VERSION := 0.1.0

run:
	docker rm -f $(NAME)
	docker run --name $(NAME) --env-file .env $(IMAGE)

build:
	docker build -t $(IMAGE) .

release:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) -f Dockerfile.release .
	docker tag $(IMAGE) $(IMAGE):$(VERSION)
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

clean:
	docker rm $$(docker ps -q -f "exited!=0")
	docker rmi $$(docker images -q -f "dangling=true")

.PHONY: run build release clean
