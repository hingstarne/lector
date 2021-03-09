@Library('fl_jenkins_shared_library@master') _

pipeline {
    agent {
        kubernetes {
            containerTemplates agentContainers(
                size: "large",
                enableDeploy: params.ENABLE_DEPLOY
            )
        }
    }

    options {
        skipDefaultCheckout(true)
        timeout(time: 25, unit: 'MINUTES')
    }

    parameters {
        booleanParam(name: 'ENABLE_DEPLOY',
                     defaultValue: (env.BRANCH_NAME == 'master' ||
                                    env.TAG_NAME ==~ /(\d{4}-\d{2}-\d{2}_\d{2}-\d{2})$/),
                     description: 'Whether or not to trigger a deployment')

        choice(name: 'TARGET_ENV',
               choices: ['integration', 'production'],
               description: 'Execute the deployment against which environment?')
    }

    stages {
        stage('Build') {
            steps {
                script {
                    env.IMAGE_TAG = build(steps: this,
                                          env: env,
                                          imageTag: env.TAG_NAME)
                }
            }
        }

        stage('Deploy: Integration') {
            when {
                allOf {
                    not { buildingTag() }
                    expression {
                        return (params.ENABLE_DEPLOY &&
                                params.TARGET_ENV.equals('integration'))
                    }
                }
            }
            steps {
                script {
                    deploy(steps: this,
                           environment: 'integration',
                           values: ['image.tag': env.IMAGE_TAG])
                }
            }
        }

        stage('Deploy: Create release') {
            when {
                allOf {
                    branch 'master'
                    expression {
                        return (params.ENABLE_DEPLOY &&
                                params.TARGET_ENV.equals('production'))
                    }
                }
            }
            steps {
                script {
                    gitCreateRelease(steps: this,
                                     repositoryUrl: env.GIT_REPOSITORY_URL)
                }
            }
        }

        stage('Deploy: Production') {
            when {
                tag pattern: /^(\d{4}-\d{2}-\d{2}_\d{2}-\d{2})$/,
                comparator: "REGEXP"
            }
            steps {
                script {
                    deploy(steps: this,
                           environment: 'production',
                           values: ['image.tag': env.IMAGE_TAG])
                }
            }
        }
    }

    post {
        always {
            report(steps: this,
                   environment: params.TARGET_ENV)
        }
        fixed {
            notify(buildStatus: 'BACK TO NORMAL')
        }
        failure {
            notify(buildStatus: 'FAILURE')
        }
    }
}
