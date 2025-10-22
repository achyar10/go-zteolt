pipeline {
  environment {
    IMAGE_NAME = "go-zte"
    dockerImage = ''
  }
  agent any
  stages {
    stage('Cloning Repository') {
      steps {
        echo 'Notify GitLab'
      }
    }
    stage('Build Image Production') {
      when {
            expression { return env.GIT_BRANCH == "origin/main" }
      }
      steps{
        sh "sed -i 's/${env.IMAGE_NAME}:latest/${env.IMAGE_NAME}:${env.BUILD_ID}/g' deployment-service.yml"
        script {
          dockerImage = docker.build("registry.matik.id/${env.IMAGE_NAME}:${env.BUILD_ID}")
        }
      }
    }
    stage('Pushing Image') {
      environment {
               registryCredential = 'registryID'
           }
      steps{
        script {
          docker.withRegistry('http://registry.matik.id', registryCredential ) {
            dockerImage.push("${env.BUILD_ID}")
          }
        }
      }
    }
    
    stage('Deploying App to Kubernetes Production') {
      when {
            expression { return env.GIT_BRANCH == "origin/main" }
      }
      steps {
        script {
          kubernetesDeploy(configs: "deployment-service.yml", kubeconfigId: "k8s")
        }
      }
    }

    stage('Remove Unused docker image') {
      steps {
          sh "docker image prune -f"
      }
    }
  }
}
