@Library('jenkins-libs') _
pipeline {
    agent { label 'x86_64-medium' }
    parameters {
        booleanParam(name: 'EXECUTE_SONARQUBE_ANALYSIS', defaultValue: true, description: 'Executes SonarQube Analysis and waits for quality gates to pass')
        booleanParam(name: 'EXECUTE_ECR_PUSH', defaultValue: false, description: 'Pushes the built image into AWS ECR')
        booleanParam(name: 'EXECUTE_UPDATE_MANIFESTS', defaultValue: false, description: 'Updates the Manifests with the new built image')
        choice(name: 'REGION', choices: getAvailableRegionChoices(), description: 'Specifies the target region for the deployment')
        choice(name: 'ENVIRONMENT', choices: ['development', 'production'], description: '[WARNING] DO NOT change this value unless you really know what it does! Sets the manifests folder name in which to update the change.')
        string(name: 'RELEASE_VERSION', defaultValue: '', description: '[NOTE] Triggering a manual production build should only be done only if automatic ones are not functioning properly. Sets a custom release version to be appended to the tagged image (i.e. i-2024-05-02-jMnq2PGaT9-v.1.0.5).')
    }
    environment {
        SCANNER_HOME = tool 'sonarqube-scanner'
        AWS_REGION = 'eu-west-1'
        AWS_ACCOUNT_ID = '429223056556'
        GITHUB_REPO = 'meeting_service'
        ECR_REPO = "${GITHUB_REPO}"
        IMAGE_TAG = sh(script: "/root/build-and-tag.sh ${params.REGION} ${params.RELEASE_VERSION}", returnStdout: true).trim()
        MANIFESTS_REPO_URL = 'https://support%40meetgeek.ai@github.com/meetgeekai/manifests.git'
        DOCKER_BUILDKIT = 1
    }
    stages {
        stage('Build') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'github-developer-key', usernameVariable: 'AUTH_USER', passwordVariable: 'AUTH_TOKEN')]) {
                        latestImage = "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}:latest"
                        docker.withRegistry("https://${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com", "ecr:${AWS_REGION}:aws-ecr-admin-creds") {
                            echo "Building project with image tag: ${IMAGE_TAG}"
                            docker.build("${latestImage}", "--cache-from ${latestImage} --build-arg BUILDKIT_INLINE_CACHE=1 --build-arg AUTH_USER=$AUTH_USER --build-arg AUTH_TOKEN=$AUTH_TOKEN .")
                        }
                    }
                }
            }
        }
        stage('SonarQube Analysis') {
            when {
                expression { params.ENVIRONMENT == 'development' && params.EXECUTE_SONARQUBE_ANALYSIS == true }
            }
            steps {
                withSonarQubeEnv('sonarqube') {
                    sh "${SCANNER_HOME}/bin/sonar-scanner"
                }
                timeout(time: 15, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        stage('Push Image') {
            when {
                expression { env.BRANCH_NAME == 'main' || params.EXECUTE_ECR_PUSH == true }
            }
            steps {
                script {
                    image = docker.image("${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}:latest")
                    docker.withRegistry("https://${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com", "ecr:${AWS_REGION}:aws-ecr-admin-creds") {
                        image.push("${IMAGE_TAG}")
                        image.push("latest")
                    }
                }
            }
        }
        stage('Update Manifests') {
            when {
                expression { env.BRANCH_NAME == 'main' || params.EXECUTE_UPDATE_MANIFESTS == true }
            }
            steps {
                withCredentials([usernamePassword(credentialsId: 'github-developer-key', usernameVariable: '', passwordVariable: 'AUTH_TOKEN')]) {
                    checkout([
                        $class: 'GitSCM',
                        branches: [[name: 'main']],
                        extensions: [[$class: 'LocalBranch', localBranch: "**"]],
                        userRemoteConfigs: [[url: MANIFESTS_REPO_URL, credentialsId: 'github-developer-key']],
                        doGenerateSubmoduleConfigurations: false
                    ])
                    script {
                        def regionsToDeploy = getRegionsToDeploy(params.REGION)

                        for (region in regionsToDeploy) {
                            dir("${params.ENVIRONMENT}/${region}/${GITHUB_REPO}") {
                                sh "yq -i \'.image.repository = \"${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}\"\' values.yaml"
                                sh "yq -i \'.image.tag = \"${IMAGE_TAG}\"\' values.yaml"
                                sh "git add values.yaml"
                                sh "git commit -m \'JenkinsCI updated image tag for ${ECR_REPO} to ${IMAGE_TAG}\'"
                                sh 'git push https://support%40meetgeek.ai:$AUTH_TOKEN@github.com/meetgeekai/manifests.git main'
                            }
                        }
                    }
                }
            }
        }
    }
    post {
        success {
            script { sendSlackSuccess(repo: env.GITHUB_REPO, branch: env.BRANCH_NAME, build_url: env.BUILD_URL) }
        }
        failure {
            script { sendSlackFailure(repo: env.GITHUB_REPO, branch: env.BRANCH_NAME, build_url: env.BUILD_URL) }
        }
    }
}
