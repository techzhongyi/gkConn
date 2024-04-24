// 注意：build httpCore的时候放到最后，因为要把apis 更新了，先build 其他会更新
def config_file = ''
def need_push = false
def first_clone = true // httpCore 每次都需要build 因为不知道其他项目的apps是否有变化
def CONFIG_FILE = '''# 基础环境配置
Server:
    Name: gkConn
    Port: 8081
    StoragePath: ./_storage
    DelayMinutes: ${DELAY_MINUTES}

Factorys:
  - # 车企用户名密码
    Name:  ${FACTORY_NAME}
    Password:  ${FACTORY_PASS}
'''
pipeline {
    agent {
        node {
          label 'gk-dev'
          customWorkspace './workspace/docker_gkzy/gkConn'
        }
    }

    stages {
        stage('@@@-初始化pipline config_file---server_name, port from CONFIG_FILE') {
            steps {
                script {
                    config_file = readYaml text: CONFIG_FILE
                    echo config_file.Server.Name
                    echo "${config_file.Server.Port}"
                    try {
                        sh 'git checkout .'
                    } catch (Exception ex) {
                        first_clone = true
                        println('ignore, before git clone...')
                    }
                }
            }
        }

        stage('@@@-拉取代码 testing') {
            steps {
                echo '准备拉取代码-gkConn + comlibgo'
                checkout([$class: 'GitSCM', branches: [[name: '*/gkzy']], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: './gkConn']], userRemoteConfigs: [[url: 'git@github.com:techzhongyi/gkConn.git']]])
                checkout([$class: 'GitSCM', branches: [[name: '*/v1.0']], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: './comlibgo']], userRemoteConfigs: [[url: 'git@github.com:techzhongyi/comlibgo.git']]])
            }
        }

        stage('@@@-创建服务配置文件-config.yaml') {
            // todo 判断config_file文件和本地比较是否有变化
            when {
                anyOf {
                    expression { return !fileExists('./gkConn/config.yaml') }
                    expression {
                        def content = readFile(file: './gkConn/config.yaml')
                        return content != CONFIG_FILE
                    }
                }
            }
            steps {
                writeFile encoding: 'utf-8', file: './gkConn/config.yaml', text: CONFIG_FILE
                echo "Success-创建服务配置文件"
            }
        }

        stage('@@@-docker build 构建 判断是否有代码变化') {
            when {
                anyOf {
                    expression { return first_clone }
                    changeset "**/Dockerfile"
                    changeset "**/*.go"
                    changeset "**/*.yaml"
                }
            } // 必须线下git 打 tag，同时完成git merge
            steps {
                script {
                    need_push = true
                }
                sh 'cp ./gkConn/Dockerfile .'
                sh "docker build --tag registry-vpc.cn-qingdao.aliyuncs.com/gkzy_sys/gkConn:latest --no-cache ."
            }
        }
        stage('Push image') {
            when {
                expression { return need_push }
            }
            steps {
                sh 'echo gkzy@6789 | docker login -u gkzy9999 --password-stdin registry-vpc.cn-qingdao.aliyuncs.com'
                sh "docker push registry-vpc.cn-qingdao.aliyuncs.com/gkzy_sys/gkConn:latest"
                echo '上传镜像成功！'
//                 sh "docker rmi registry-vpc.cn-qingdao.aliyuncs.com/zynewsapce/${config_file.Server.Name}:latest"
//                 echo '删除本地镜像！'
            }
        }
    }
}
