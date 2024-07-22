build: | Dockerfile
	@docker build -t tmpl-greeting .
.PHONY: build

clean:
	@rm -f Dockerfile
.PHONY: clean

Dockerfile:
	@tmpl generate \
		-c config.yml \
		-c Dockerfile.yml \
		-o Dockerfile \
		-m Dockerfile.tmpl:/Dockerfile.tmpl \
		-m includes:/includes \
		/Dockerfile.tmpl
