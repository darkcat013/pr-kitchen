# Kitchen
Laboratory work nr1 on Network programming course.

# Build and run application in docker
```
docker build -t kitchen .
docker run -p 8080:8080 -it kitchen
```
For Linux: 
```
docker build -t kitchen .
docker run --add-host host.docker.internal:host-gateway -p 8080:8080 -it kitchen
```

# Run application locally
```
go run .
```

# URL
```
http://host.docker.internal:8080
```