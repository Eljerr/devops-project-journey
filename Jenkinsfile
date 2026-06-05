pipeline {
    agent any

    environment {
        DOCKERHUB = credentials('dockerhub-creds')
        IMAGE_NAME = "mybinichizuru/project-02-api"
        IMAGE_TAG = "jenkins-${env.GIT_COMMIT.take(7)}"
    }

    stages {
        stage('Checkout Code') {
            steps {
                echo 'Mengambil kode dari Github...'
                git branch: 'main', url: 'https://github.com/Eljerr/devops-project-journey.git'
            }
        }

        stage('Build') {
            steps {
                echo 'Merakit Docker Image...'
                dir('src-app') {
                    sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
                }
            }
        }

        stage('Menjalankan Trivy Scanning') {
            steps {
                echo 'Menjalankan Security Scan pakai Trivy...'
                sh "trivy image --exit-code 1 --severity HIGH,CRITICAL --ignore-unfixed ${IMAGE_NAME}:${IMAGE_TAG}"
            }
        }

        stage('Push to Docker Hub') {
            steps {
                echo 'Login dan Mendorong Image ke Docker Hub...'
                sh 'echo $DOCKERHUB_PSW | docker login -u $DOCKERHUB_USR --password-stdin'
                sh "docker push ${IMAGE_NAME}:${IMAGE_TAG}"
            }
        }
    } 

    post {
        success {
            echo '✅ Pipeline DevSecOps berhasil!'
        }
        failure {
            echo '❌ Pipeline GAGAL! Segera cek Log.'
        }
        always {
            sh "docker rmi ${IMAGE_NAME}:${IMAGE_TAG} || true"
        }
    }
}
