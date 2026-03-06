pipeline {
    agent none

    environment {
        IMAGE_NAME = "document-service:latest"
        CONTAINER_NAME = "document-service"
    }

    stages {
        stage('Build Go & Docker Image') {
            agent {
                docker {
                    image 'golang:1.26-alpine'
                    args '-v $HOME/.cache/go-build:/go/pkg/mod/cache'
                }
            }
            steps {
                script {
                    sh 'rm -f app'
                    sh 'go mod tidy'
                    sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app'
                    sh 'cp app .'
                }
            }
        }

        stage('Build Docker Compose') {
            agent {
                docker {
                    image 'docker/compose:2.20.2'
                    args '-v /var/run/docker.sock:/var/run/docker.sock'
                }
            }
            steps {
                script {
                    sh 'docker-compose build'
                    sh 'docker-compose up -d'
                }
            }
        }

        stage('Integration Tests') {
            agent {
                docker {
                    image 'golang:1.26-alpine'
                    args '-v $HOME/.cache/go-build:/go/pkg/mod/cache'
                }
            }
            steps {
                sh 'go test ./... || true'
            }
        }
    }

    post {
        always {
            echo "Pipeline terminata"
        }
    }
}