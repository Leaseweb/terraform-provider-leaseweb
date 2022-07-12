#!groovy

lswci([node: 'docker', mattermost: 'bare-metal-cicd']) {
    lswWithDockerContainer(image: 'artifactory.devleaseweb.com/lswci/golang:1.18') {
        stage("Lint") {
            sh "make ci"
        }
        stage("Test build") {
            sh "make build"
        }
        if (env.BRANCH_NAME == 'master') {
            stage("Build release") {
                sh "make release"
            }
            stage("Publish artifacts") {
                # to-do
            }
        }
    }
}
