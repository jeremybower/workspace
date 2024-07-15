.PHONY: build
build: | Dockerfile
	@docker build -t tmpl-greeting .

.PHONY: clean
clean:
	@rm -f Dockerfile

Dockerfile:
	@tmpl -c ../config.yml -c Dockerfile.yml -o Dockerfile Dockerfile.tmpl includes/*.tmpl
