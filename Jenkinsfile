pipeline {
    agent none

    environment {
        IMAGE_NAME = "document-service:latest"
        CONTAINER_NAME = "document-service"
    }

    stages {

        stage('Build Go Binary') {
            agent any
            steps {
                echo "Build Go Binary inside a container"
                sh '''
                docker run --rm -v $PWD:/app -w /app golang:1.26-alpine sh -c '
                    go mod tidy
                    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app
                '
                '''
            }
        }

        stage('Build & Start Services with Docker Compose') {
            agent any
            steps {
                echo "Build and launch containers with Docker Compose"
                sh '''
                docker run --rm -v $PWD:/app -w /app -v /var/run/docker.sock:/var/run/docker.sock docker/compose:2.20.2 sh -c '
                    docker-compose build
                    docker-compose up -d
                '
                '''
            }
        }

        stage('Integration Tests') {
            agent any
            steps {
                echo "Run integration tests inside Go container"
                sh '''
                docker run --rm -v $PWD:/app -w /app golang:1.26-alpine sh -c '
                    go test ./...
                '
                '''
            }
        }
    }

    post {
        always {
            echo "Cleanup containers"
            sh '''
            docker run --rm -v $PWD:/app -w /app -v /var/run/docker.sock:/var/run/docker.sock docker/compose:2.20.2 sh -c '
                docker-compose down
            '
            '''
        }
    }
}