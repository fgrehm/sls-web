.PHONY: build build.watch serve.api serve.client work release release.build release.package

SHELL = /bin/bash
VERSION = 0.0.1

default: build

release:
	hub release create -d \
		-a release/sls-web-osx-v$(VERSION).tgz \
		-a release/sls-web-linux-v$(VERSION).tgz \
		"v$(VERSION)"
	@echo 'Now go to https://github.com/fgrehm/sls-web/releases and publish the release if it looks good'

release.build:
	mkdir -p release/osx
	mkdir -p release/linux
	rm -rf release/linux/client && rm -rf release/osx/client
	cd src/client && npm run build
	cp -R src/client/dist release/linux/client
	cp -R src/client/dist release/osx/client
	env GOOS=linux GOARCH=amd64 go build -o release/linux/sls-web ./src/cli
	env GOOS=darwin GOARCH=amd64 go build -o release/osx/sls-web ./src/cli
	if ! [[ -f release/linux/san-lite-solver ]]; then \
		curl -sL 'https://sites.google.com/site/afonsosales/tools/linux64/san-lite-solver?attredirects=0' > release/linux/san-lite-solver; \
		chmod +x release/linux/san-lite-solver; \
	fi
	if ! [[ -f release/osx/san-lite-solver ]]; then \
		curl -sL 'https://sites.google.com/site/afonsosales/tools/mac64/san-lite-solver?attredirects=0' > release/osx/san-lite-solver; \
		chmod +x release/osx/san-lite-solver; \
		fi

release.package:
	rm -f sls-web-linux-*.tgz
	cd release && tar cvzf sls-web-linux-v$(VERSION).tgz --transform 's,^linux,sls-web,' linux
	cd release && tar cvzf sls-web-osx-v$(VERSION).tgz --transform 's,^osx,sls-web,' osx

build:
	go build -o tmp/sls-web ./src/cli

build.watch:
	@make build || true
	reflex -r '\.go$$' -- sh -c 'clear && echo "Building..." && make build && echo "Done!"'

serve.api:
	reflex -s -r '\.go$$' -- sh -c 'make build && ./tmp/sls-web serve'

serve.client:
	cd src/client && npm run dev

work:
	reflex -s -r 'sls-web$$' -- sh -c './tmp/sls-web work'
