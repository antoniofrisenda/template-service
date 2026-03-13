pipeline {
    agent none

    stages {
        stage('Kube Test') {
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
                sh '''
                    echo "Harbor12345" | docker login 192.168.1.10:8083 \
                        --username "admin" \
                        --password-stdin

                    docker build -f .docker/Dockerfile -t 192.168.1.10:8083/document-service:latest .
                    docker push 192.168.1.10:8083/document-service:latest
                '''
            }
        }

        stage('Init S3 Bucket') {
            agent { label 'k8s' }
            steps {
                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                    sh '''
                        kubectl --context kind-library delete job init-s3-bucket --ignore-not-found
                        kubectl --context kind-library apply -f .docker/.k8s/init-s3-bucket-job.yml
                        kubectl --context kind-library wait --for=condition=complete job/init-s3-bucket --timeout=120s
                    '''
                }
            }
        }

        stage('Deploy') {
            agent { label 'k8s' }
            steps {
                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                    sh '''
                        kubectl --context kind-library apply -f .docker/.k8s/deployment.yml
                        kubectl --context kind-library apply -f .docker/.k8s/service.yml
                        kubectl --context kind-library rollout status deployment/document-service
                    '''
                }
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