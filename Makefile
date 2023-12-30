run:
	air

build:
	templ
	tailwind
	go build -o ./go-htmx-example main.go

templ:
	templ generate

## tailwind: build tailwind
tailwind:
	npx tailwindcss -i ./public/styles/tailwind-input.css -o ./public/styles/tailwind-output.css --minify

## tailwind-watch: watch build tailwind
tailwind-watch:
	npx tailwindcss -i ./public/styles/tailwind-input.css -o ./public/styles/tailwind-output.css --watch