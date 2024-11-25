TAG_NAME:=$(shell date "+%Y%d%m_%H%M%S")

# Use docker buildkit and labeling for all images.
define docker_build
	DOCKER_BUILDKIT=1 docker build --label "part-of=prm" $(1)
endef

.PHONY: help client clean

help:        ## Show this help.
             ##
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

client:      ## Build client binary image.
             ##
	$(call docker_build, -t prm/$@:${TAG_NAME} -f Dockerfile .)

clean:       ## Clean build time resources, including Docker images.
             ##
	docker rmi -f $(shell docker images -f "label=part-of=prm" -q)