

build-win:
	wails build -platform=windows/amd64

build-win-dev:
	wails build -platform=windows/amd64 -debug -devtools

build-mac:
	wails build -platform=darwin/universal

build-mac-dev:
	wails build -platform=darwin/universal -debug -devtools
