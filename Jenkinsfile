#!groovy

lswci([node: 'docker', mattermost: 'bare-metal-cicd']) {
    lswWithDockerContainer(image: 'artifactory.devleaseweb.com/lswci/golang:1.18') {
        stage("Lint") {
            sh "go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2"
            sh "golangci-lint run --disable-all -E gofmt"
            sh "golangci-lint run --disable-all -E whitespace"
            sh "golangci-lint run --disable-all -E errcheck"
        }
        if (env.BRANCH_NAME == 'master') {
            stage("Build") {
                sh "make release"
            }
            stage("Publish artifacts") {
                # to-do
            }
        }
    }
}
