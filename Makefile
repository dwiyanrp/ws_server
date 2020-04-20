run-docker:
	@docker build --rm -f "Dockerfile" -t optimizegobwas:latest "." && docker run --rm -it  -p 8000:8000/tcp optimizegobwas:latest
