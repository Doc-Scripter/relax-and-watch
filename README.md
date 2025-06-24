# Relax and Watch

A streaming service application built with Go backend and React frontend.

## Deployment on Render

1. Push your code to a GitHub repository
2. Create a new Render account or sign in
3. Click "New" and select "Web Service"
4. Connect your GitHub repository
5. Render will automatically detect the `render.yaml` file and configure the services
6. Click "Create Web Service" to deploy

For more details, see [Render's documentation](https://render.com/docs/deploy-go)

## Local Development

Run the backend:
```
make run
```

Format Go code:
```
make format
```

Restart server:
```
make restart-server
```