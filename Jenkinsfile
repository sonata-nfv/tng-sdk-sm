pipeline {
  agent any
  stages {
    stage('Unit Test') {
      parallel {
        stage('Unit Tests') {
          steps {
            echo 'Performing Unit Tests'
          }
        }
        stage('Running Unit Tests') {
          steps {
            echo 'TODO'
          }
        }
      }
    }
  }
  post {
    always {
      echo 'TODO'
    }
  }
}
