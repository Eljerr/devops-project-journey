pipeline {
    agent any

    stages {
        stage('Checkout Code') {
            steps {
                echo 'Mengambil kode dari Github...'
                git branch: 'main', url: 'https://github.com/Eljerr/devops-project-journey'
              } 
          }

        stage('Build') {
            steps {
                echo 'Sedang melakukan proses Build Docker Image...'
              }
          }

        stage('Test') {
            steps {
                echo 'Menjalankan Security Scan pakai Trivy...'
              }
          }

        stage('Deploy') {
            steps {
              echo 'Mendorong Image ke Docker Hub & Update Manifest K8s!'
              }
          }
      }
  }
