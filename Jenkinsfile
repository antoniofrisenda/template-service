pipeline {
  agent none

  stages {

    stage('Build Go Binary') {
      agent { label 'go' } // nodo con Go
      steps {
        sh '''
        go mod tidy
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./src/cmd/app
        '''
      }
    }

    stage('Docker Compose') {
      agent { label 'docker' } // nodo con Docker
      steps {
        sh '''
        docker-compose build
        docker-compose up -d
        '''
      }
    }

    stage('Integration Tests') {
      agent { label 'go' } // nodo con Go
      steps {
        sh 'go test ./...'
      }
    }

  }

  post {
    always {
      agent { label 'docker' } // nodo con Docker
      steps {
        sh 'docker-compose down'
      }
    }
  }
}