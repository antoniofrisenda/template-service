pipeline {
    agent none

    stages {
        stage('Kubectl Test') {
            agent { label 'k8s' }
            steps {
                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                    sh 'ls -l "$KUBECONFIG"'
                    sh 'kubectl cluster-info'
                    sh 'kubectl get nodes'
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
                sh 'docker build -f .docker/Dockerfile -t document-service .'
            }
        }

        stage('Run Services') {
            agent { label 'docker' }
            environment {
                PORT = '3000'
                DB = 'templates'
                MONGO_URL = 'mongodb://root:pass@mongo:27017/templates?authSource=admin&w=majority'
                AWS_ENDPOINT_URL = 'http://localstack:4566'
                AWS_DEFAULT_REGION = 'us-east-1'
                AWS_ACCESS_KEY_ID = 'test'
                AWS_SECRET_ACCESS_KEY = 'test'
                AWS_S3_BUCKET_NAME = 'document-bucket'
            }
            steps {
                sh 'docker-compose -f .docker/docker-compose.yml down -v'
                sh 'docker-compose -f .docker/docker-compose.yml up -d'
            }
        }

        stage('Init S3 Bucket') {
            agent { label 'docker' }
            steps {
                sh '''
                until docker exec localstack awslocal s3 ls >/dev/null 2>&1; do
                  sleep 2
                done
                docker exec localstack awslocal s3 mb s3://document-bucket
                docker exec localstack awslocal s3 ls
                '''
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