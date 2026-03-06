pipeline {
    agent {
        docker {
            image 'golang:1.21'
            args '-v /var/jenkins_home/workspace:/workspace'
        }
    }
    stages {
        stage('Build') {
            steps {
                sh 'go mod tidy'
                sh 'go build -o myapp ./...'
            }
        }
    }
}