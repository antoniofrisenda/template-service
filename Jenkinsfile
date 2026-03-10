pipeline {
    agent none

    stages {

        stage('Build Go') {
            agent { label 'golang' }
            steps {
                sh 'rm -f app'
                sh 'go mod download'
                sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app'
            }
        }

        stage('Test Kubectl') {
            agent { label 'k8s' }
            steps {
                sh 'kubectl'
            }
        }

        stage('Test') {
            agent { label 'golang' }
            steps {
                sh 'go test ./... || true'
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
        always {
            echo 'Pipeline completata.'
        }
    }
} 
