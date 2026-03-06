pipeline {
    agent any

    options {
        timeout(time: 10, unit: 'MINUTES')
    }

    tools {
        go 'go-1.26'
    }

    environment {
        GO111MODULE = 'on'
    }

    stages {
        stage('Build Go') {
            steps {
                sh '''
                    rm -f document-service
                    go mod tidy
                    go build -o document-service
                '''
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    docker-compose down || true
                    docker-compose build
                    docker-compose up -d
                '''
            }
        }
    }

    post {
        success {
            echo 'Pipeline terminata.'
        }
        failure {
            echo 'Pipeline fallita.'
        }
    }
}