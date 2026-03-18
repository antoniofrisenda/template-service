pipeline {
    agent none

    options {
        skipDefaultCheckout(true)
    }

    stages {
        stage('Checkout Source') {
            agent { label 'go' }
            steps {
                deleteDir()
                checkout scm
                stash name: 'source', includes: '**', excludes: '.git/**', useDefaultExcludes: false
            }
        }

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
                deleteDir()
                unstash 'source'
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
                deleteDir()
                unstash 'source'
                sh 'go test ./... -coverprofile=coverage.out'
                stash name: 'coverage', includes: 'coverage.out'
            }
        }
        
        stage('Sonar Scan') {
            agent { label 'sonar-scanner' }
            steps {
                deleteDir()
                unstash 'source'
                unstash 'coverage'
                sh '''
                sonar-scanner \
                    -Dproject.settings=sonar.properties \
                    -Dsonar.host.url=http://sonarqube:9000 \
                    -Dsonar.token=sqp_85265a6a0df86dd095f39024454d18e08ef33822
                '''
            }
        }

        stage('Build Docker Image') {
            agent { label 'docker' }
            steps {
                deleteDir()
                unstash 'source'
                sh 'docker build -f .docker/Dockerfile -t document-service .'
            }
        }

        stage('Run Services') {
            agent { label 'docker' }
            steps {
                deleteDir()
                unstash 'source'
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
