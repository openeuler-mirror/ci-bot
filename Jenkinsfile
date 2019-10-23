pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('build') {
            steps {
                go get golang.org/x/crypto/ssh
                go get -v -tags 'fixtures acceptance' ./... 
                go get github.com/wadey/gocovmerge
                go get github.com/mattn/goveralls
                go get golang.org/x/tools/cmd/goimports
                ./script/coverage
            }
        }
    }
}
