pipeline {
    agent {
        label 'upf-agent'
    }
    environment { 
        HTTP_PROXY = 'http://proxy-dmz.intel.com:911'
        HTTPS_PROXY = 'http://proxy-dmz.intel.com:912'
    }
    stages {
        stage('Code Checkout') {
            steps {
                // Check out the code from your version control system
                git 'https://github.com/intel-innersource/frameworks.automation.dtac.agent'
            }
        }
        
        stage('Build') {
            steps {
                container('jenkins-agent-upf') {
                    // Compile code and resolve dependencies
                    sh 'echo "[+] Building..."'
                }
            }
        }
        
        stage('Unit Testing') {
            steps {
                container('jenkins-agent-upf') {
                    // Run unit tests
                    sh 'echo "[+] Running Unit Tests..."'
                    sh 'mkdir .coverage'
                    sh 'go test -race -failfast -coverprofile=.coverage/coverage-unit.txt -covermode=atomic -v ./pfcpiface ./cmd/...'   
                }
            }
        }
        
        stage('Integration Testing') {
            steps {
                // Run integration tests
                sh 'echo "[+] Running Integration Tests..."'
            }
        }
        
        stage('Functionality Testing') {
            steps {
                // Run functionality tests
                sh 'echo "[+] Running Functionality Testing..."'
            }
        }
        
        stage('Performance Testing') {
            steps {
                // Run performance tests
                sh 'echo "[+] Running Performance Testing..."'
            }
        }
        
        stage('Deployment') {
            steps {
                // Deploy UPF to a staging environment
                sh 'echo "[+] Running Deployment Testing..."'
            }
        }
        
        stage('End-to-End Testing') {
            steps {
                // Conduct end-to-end tests
                sh 'echo "[+] Running End-to-End Testing..."'
            }
        }
        
        stage('Reporting') {
            steps {
                // Generate test reports and metrics
                sh 'echo "[+] Reporting..."'
            }
        }
        
        stage('Artifact Management') {
            steps {
                // Store successful artifacts
                sh 'echo "[+] Storing Artifacts..."'
            }
        }
        
        stage('Notifications') {
            steps {
                // Send notifications for test results and deployment status
                sh 'echo "[+] Sending Notifications..."'
            }
        }
    }
    post {
        success {
            // Actions to take when the pipeline succeeds
            sh 'echo "[+] Run Succeeded Do Something..."'
        }
        
        failure {
            // Actions to take when the pipeline fails
            sh 'echo "[-] Run Failedd Do Something..."'
        }
    }
}
