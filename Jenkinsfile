pipeline {
    agent none

    stages {
        stage('Test Kubectl') {
            agent { label 'k8s' }
            steps {
                sh '''
                    # Percorso temporaneo per il kubeconfig
                    export KUBECONFIG=$(mktemp)

                    # Genera un kubeconfig minimo per Docker Desktop
                    kubectl config set-cluster docker-desktop \
                        --server=https://kubernetes.docker.internal:6443 \
                        --insecure-skip-tls-verify=true

                    kubectl config set-context docker-desktop \
                        --cluster=docker-desktop \
                        --user=docker-desktop

                    kubectl config use-context docker-desktop

                    # Verifica la connessione
                    kubectl cluster-info
                '''
            }
        }

        stage('Build Go') {
            agent { label 'golang' }
            steps {
                sh 'rm -f app'
                sh 'go mod download'
                sh 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app'
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
