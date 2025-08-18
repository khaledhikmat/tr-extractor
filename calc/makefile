dockerize: 
	docker buildx build --platform linux/amd64 -t khaledhikmat/tr-inheritance-calc:latest .

push: dockerize
	docker login
	docker push khaledhikmat/tr-inheritance-calc:latest
