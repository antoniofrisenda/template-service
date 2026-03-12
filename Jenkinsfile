pipeline {
    agent none

    stages {
        stage('Kubectl Test') {
            agent { label 'k8s' }
            steps {
                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                    sh 'ls -l "$KUBECONFIG"'
                    sh 'kubectl --kubeconfig="$KUBECONFIG" cluster-info'
                    sh 'kubectl --kubeconfig="$KUBECONFIG" get nodes'
                }
            }
        }

        stage('Build Go') {
            agent { label 'go' }
            steps {
                sh 'rm -f app'
                sh 'go get -u ./...'
                sh 'go mod tidy'
                sh 'go mod download'
                sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app'
            }
        }

        stage('Go Test') {
            agent { label 'go' }
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
            echo 'Pipeline NOT OK!'
        }
    }
}