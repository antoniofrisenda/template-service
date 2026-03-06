pipeline {
    agent any

    options {
        timeout(time: 10, unit: 'MINUTES')
    }

    tools {
        go 'go-1.26'
    }

    environment {
        IMAGE_NAME = "document-service:latest"
        CONTAINER_NAME = "document-service"
        GO111MODULE = 'on'
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