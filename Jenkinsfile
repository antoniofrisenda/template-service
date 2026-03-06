pipeline {
    agent any

    options {
        timeout(time: 10, unit: 'MINUTES')
    }

    environment {
        GO111MODULE = 'on'
    }

    stages {

        stage('Checkout SCM') {
            steps {
                checkout scm
            }
        }

        stage('List Workspace (Debug)') {
            steps {
                sh 'pwd && ls -R'
            }
        }

        stage('Build Go Binary') {
            steps {
                // Cambia 'cmd/document-service' con la cartella reale dei tuoi file .go
                dir('cmd/document-service') {
                    sh '''
                        echo "Cleaning old binary..."
                        rm -f ../../../document-service

                        echo "Tidying modules..."
                        go mod tidy

                        echo "Building Go binary..."
                        go build -o ../../../document-service
                    '''
                }
            }
        }

        stage('Deploy with Docker Compose') {
            steps {
                sh '''
                    echo "Stopping old containers..."
                    docker-compose down || true

                    echo "Building Docker images..."
                    docker-compose build

                    echo "Starting containers..."
                    docker-compose up -d
                '''
            }
        }
    }

    post {
        success {
            echo 'Pipeline terminata con successo.'
        }
        failure {
            echo 'Pipeline fallita.'
        }
    }
}