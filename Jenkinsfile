pipeline {
    agent { label 'master' }

    options {
        ansiColor('xterm')
    }

    stages {
        stage('Checkout') {
            steps {
                script {
                env.BUILD_VERSION = VersionNumber projectStartDate: '2019-11-13', versionNumberString: "${BUILD_NUMBER}", versionPrefix: "1.", worstResultForIncrement: 'SUCCESS'
                echo "Version:  ${BUILD_VERSION}"

                GIT_REVISION = sh(returnStdout: true, script: """
                    git rev-parse --short HEAD"""
                ).trim()
                def PWD = pwd()
                echo "Check out REVISION: $GIT_REVISION on $PWD"
                DO_GATHER_ARTIFACT_BRANCH = (['master', 'jenkins'].contains(GIT_BRANCH) ||
                    GIT_BRANCH ==~ /release\-[\d\-\.]+/ ||
                    GIT_BRANCH ==~ /[^\s]+enable_docker_image_push$/)
                checkout changelog: false, poll: false, scm: [$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'jenkins-helper']], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'github-personal-jenkins', url: 'https://github.com/sunshine69/jenkins-helper.git']]]
                utils = load("${WORKSPACE}/jenkins-helper/deployment.groovy")
                }//script
            }//steps
        }//stage

        stage('Generate build scripts') {
            steps {
                script {
                  utils.generate_add_user_script()
                  //utils.generate_aws_environment()
                  withCredentials([usernamePassword(credentialsId: 'github-personal-jenkins', passwordVariable: 'GIT_TOKEN', usernameVariable: 'GIT_USER')]) {
                  withCredentials([usernamePassword(credentialsId: 'window-build-credential', passwordVariable: 'ADMIN_PASSWORD', usernameVariable: 'ADMIN_USER')]) {
                  sh '''cat <<EOF > build-windows.sh
ansible-playbook ansible/build-windows.yaml -i ansible/inventory -e "admin_user=$ADMIN_USER git_user=${GIT_USER} git_token=${GIT_TOKEN} admin_password=${ADMIN_PASSWORD}"
EOF
'''
                    sh 'chmod +x build-windows.sh'
                  }//withCred
                  }//withCred
                  sh '''cat <<EOF > build.sh
ARCH=\$(uname -m)
/usr/local/go/bin/go build --tags "json1 fts5 secure_delete" -ldflags='-s -w' -o gnote-ubuntu1804-\${ARCH}
mkdir -p gnote-ubuntu1804/
cp -a glade icons gnote-ubuntu1804-\${ARCH} gnote-ubuntu1804/
tar czf gnote-ubuntu1804-\${ARCH}.tgz gnote-ubuntu1804
rm -rf gnote-ubuntu1804
EOF
'''
                }//script
            }//steps
        }

        stage('Run the command within the docker environment') {
            steps {
                script {
                    sh "./build-windows.sh"
                    utils.run_build_script([
                    //Make sure you build this image ready - having user jenkins and cache go mod for that user.
                        'docker_image': 'golang-ubuntu1804-build:latest',
                        'docker_net_opt': '',
//define the name here so we can save a image cache in the script save-docker-image-cache.sh
                        'docker_extra_opt': '--name golang-ubuntu1804-build-jenkins',
//Uncomment these when you build with the golang-alpine from scratch. After we
//can commented out as the image is saved
                        'outside_scripts': ['save-docker-image-cache.sh'],
                        //'extra_build_scripts': ['fix-godir-ownership.sh'],
                        //'run_as_user': ['fix-godir-ownership.sh': 'root'],
                    ])
                }//script
            }//steps
        }//stage

        stage('Gather artifacts') {
            steps {
                script {
                    utils.save_build_data(['artifact_class': 'gnote'])

                    if (DO_GATHER_ARTIFACT_BRANCH) {
                      archiveArtifacts allowEmptyArchive: true, artifacts: 'gnote-ubuntu1804*.tgz,gnote-windows-bundle*.zip,gnote.exe', fingerprint: true, onlyIfSuccessful: true
                      if (GIT_BRANCH ==~ /master/ ) {
                        echo "Create a release as this is a master merge ..."
                        withCredentials([usernamePassword(credentialsId: 'github-personal-jenkins', passwordVariable: 'GITHUB_TOKEN', usernameVariable: 'GITHUB_USER')]) {
                            env.REPOSITORY = "gnote"
                            sh """
                            ARTIFACT_FILES=\$(ls *.tgz gnote.exe .tgz gnote-windows-bundle*.zip)
                            git tag v${BUILD_VERSION}; git push http://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/${GITHUB_USER}/${REPOSITORY} --tags
                            GITHUB_USER=$GITHUB_USER REPOSITORY=${REPOSITORY} GITHUB_TOKEN=$GITHUB_TOKEN ARTIFACT_FILES=\${ARTIFACT_FILES} ./create-github-release.sh"""
    // some block
                        }
                      }
                    }
                    else {
                      echo "Not collecting artifacts as branch"
                    }// If GIT_BRANCH
                } //script
            }
        }// Gather artifacts

    }
    post {
        always {
            script {
                utils.apply_maintenance_policy_per_branch()
                currentBuild.description = """Artifact version: ${BUILD_VERSION}
Artifact revision: ${GIT_REVISION}"""
            }
        }
        success {
            script {
              cleanWs cleanWhenFailure: false, cleanWhenNotBuilt: false, cleanWhenUnstable: false, deleteDirs: true, disableDeferredWipeout: true, notFailBuild: true, patterns: [[pattern: 'PerformanceTesting', type: 'INCLUDE'], [pattern: '*', type: 'INCLUDE']]
            }//script
        }
        failure {
            script {
                //slackSend baseUrl: 'https://xvt.slack.com/services/hooks/jenkins-ci/', botUser: true, channel: '#errcd-activity', message: "@here CRITICAL - ${JOB_NAME} (${BUILD_URL}) branch (${BRANCH_NAME}) revision (${GIT_REVISION}) on (${NODE_NAME})", teamDomain: 'xvt', tokenCredentialId: 'jenkins-ci-integration-token', color: "danger"
                echo 'Build failed'
            }//script
        }
    }
}
