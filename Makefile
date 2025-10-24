OUT ?= bin/ypoker
ICON ?= internal/ui/assets/yChat.png

build:
	@go build -o $(OUT)

run: build
	@./$(OUT)

dev:
	@air

test:
	go test ./...

package-windows:
	@mkdir -p dist
	@CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 fyne package -os windows -icon $(ICON)
	@mv *.exe dist/ || true

package-linux:
	@mkdir -p dist
	@fyne package -os linux -icon $(ICON)
	@mv *.tar.xz dist/ || true
