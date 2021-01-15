build:
	docker build -t patoarvizu/residential-dns:latest .

push: build
	docker push patoarvizu/residential-dns:latest