# Golang Training

Held by xendit peeps :)

## Instruction on how to execute `rest` command with docker

### Build Go app

```
GOOS=linux GARCH=amd64 go build -o engine .
```

### Build Docker Image
```
docker build -t test-image-name . -f Dockerfile.rest
```

### Run the app
```
docker run -t test-image-name
```

## Instruction on how to execute `cron` command with docker

### Build Go app

```
GOOS=linux GARCH=amd64 go build -o engine .
```

### Build Docker Image
```
docker build -t test-image-name . -f Dockerfile.cron
```

### Run the app
```
docker run -t test-image-name
```