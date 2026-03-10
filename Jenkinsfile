pipeline {
    agent none

    stages {
        stage('Build Go') {
            agent { label 'golang' }
            steps {
                sh 'rm -f app'
                sh 'go get -u ./...'
                sh 'go mod tidy'
                sh 'go mod download'
                sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app'
            }
        }

        stage('Go Test') {
            agent { label 'golang' }
            steps {
                sh 'go test ./...'
            }
        }

        stage('Build Docker Image') {
            agent { label 'docker' }
            steps {
                sh 'docker build -f .docker/service.Dockerfile -t document-service .'
            }
        }

        stage('Run Services') {
            agent { label 'docker' }
            steps {
                sh 'docker-compose -f .docker/docker-compose.yml down -v'
                sh 'docker-compose -f .docker/docker-compose.yml up -d'
            }
        }
    }

    post {
        success{
            echo 'Pipeline OK!'
        }
        failure {
            echo 'Pipeline FAIL!'
        }
    }
}