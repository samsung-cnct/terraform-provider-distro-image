// -*- mode: groovy -*-
// Jenkins pipeline
// See documents at https://jenkins.io/doc/book/pipeline/jenkinsfile/

podTemplate(label: 'tf-distroimage-go', containers: [
        containerTemplate(name: 'jnlp', image: 'jenkinsci/jnlp-slave:2.62-alpine', args: '${computer.jnlpmac} ${computer.name}'),
        containerTemplate(name: 'golang', image: 'golang:latest', ttyEnabled: true, command: 'cat'),
]) {
    node('tf-distroimage-go') {
        container('golang'){
            stage('checkout') {
                checkout scm
                sh 'go version'
            }    

            stage('build') {
                sh 'go get -v -d -t ./... || true'
                sh 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v'
            }

            stage('test') {
                sh 'TF_ACC=Y go test -v'
            }
        }
    }
}  
