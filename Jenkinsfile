#!groovy

lswci([node: 'docker', mattermost: 'bare-metal-cicd']) {
    name = env.CHANGE_BRANCH ? env.CHANGE_BRANCH.toLowerCase().replace("/", "-") : env.BRANCH_NAME.toLowerCase().replace("/", "-")

    stage('Lint and test build') {
        image = docker.build("${name}-lint", "--target goci .")
        image.inside {
            stage("Lint") {
                sh "make ci"
            }
            stage("Test build") {
                sh "make build"
            }
        }
    }
    if (env.BRANCH_NAME == 'master') {
        image = docker.build("${name}-build", "--target gobuilder .")
        image.inside {
            stage("Build release") {
                sh "make release"
            }

            stage("Publish artifacts") {
                // to-do
            }
        }
    }
}
