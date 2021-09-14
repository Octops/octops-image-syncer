.PHONY: docker

default: clean build

clean:
	rm -rf bin/*

build: clean
	go build -o bin/octops-image-syncer .

docker:
	docker build -t octops-image-syncer:v0.0.1 .

deploy: docker
	docker save --output bin/octops-image-syncer-v0.0.1.tar octops-image-syncer:v0.0.1
	rsync -v bin/octops-image-syncer-v0.0.1.tar ${SSH_REMOTE}:/home/octops/
	ssh ${SSH_REMOTE} sudo -S k3s ctr images import /home/octops/octops-image-syncer-v0.0.1.tar
	kubectl apply -f hack/install.yaml

destroy:
	kubectl delete -f hack/install.yaml