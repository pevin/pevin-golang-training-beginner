# Golang Training

Held by xendit peeps :)

## Instruction on how to run app thru docker

### Build Docker Image
```
docker build -t test-image-name .
```

### Run the app - rest
```
docker run -t test-image-name
```

### Run the app - cron
```
docker run -t test-image-name /app/run_cron_entrypoint.sh
```