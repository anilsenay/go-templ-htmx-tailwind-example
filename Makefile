run:
	browser-sync start --port 3000 --proxy 'http://localhost:8080' --no-ui --no-open &
	air

build:
	templ
	tailwind
	go build -o ./go-htmx-example main.go

templ:
	templ generate

## tailwind: build tailwind
tailwind:
	tailwindcss -i ./public/styles/tailwind-input.css -o ./public/styles/tailwind-output.css --minify

## tailwind-watch: watch build tailwind
tailwind-watch:
	tailwindcss -i ./public/styles/tailwind-input.css -o ./public/styles/tailwind-output.css --watch