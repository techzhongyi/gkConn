def config_file = ''
def venv_first = false
def CONFIG_FILE = '''# 基础环境配置
Server:
    Name: gkConn
    Debug: True
    Port: 8081
    LogPath: ./_log/

Redis:
  Host: localhost
  Port: 6379
  Password: redis@6789
  MaxLen: 120000
  Db:
    cache: 11  # 缓存

Factorys:
  - # 车企用户名密码
    Name: name1
    Password: password1
  - # 车企用户名密码
    Name: name2
    Password: password2
'''

pipeline {
    agent {
        node {
          label 'gkDataCenter-dev'
          customWorkspace './workspace/gkDataCenter/gkConn'
        }
    }

    environment {
        GOPATH = "${WORKSPACE}/build"
        GOPROXY = "https://goproxy.cn,direct"
        GOROOT = "/usr/local/go"
    }

    stages {
        stage('@@@-初始化pipline config_file---server_name, port from CONFIG_FILE') {
            steps {
                script {
                    config_file = readYaml text: CONFIG_FILE
                    echo config_file.Server.Name
                    echo "${config_file.Server.Port}"
                    dir('./code') {
                        try {
                            sh 'git checkout .'
                        } catch (Exception ex) {
                            println('ignore, before git clone...')
                        }
                    }
                }
            }
        }

        stage('@@@-拉取代码') {
            steps {
                echo '拉取代码-comlibgo'
                checkout([$class: 'GitSCM', branches: [[name: '*/gk']], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: './comlibgo']], userRemoteConfigs: [[url: 'git@github.com:techzhongyi/comlibgo.git']]])

                echo "拉取代码-${config_file.Server.Name}"
                checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: './code']], userRemoteConfigs: [[url: 'git@github.com:techzhongyi/'+config_file.Server.Name +'.git']]])
            }
        }

        stage('@@@-go build') {
            steps {
                dir('./code') {
                    sh 'go install'
                }
            }
        }

        stage('创建服务配置文件-config.yaml') {
            // todo 判断config_file文件和本地比较是否有变化
            when {
                anyOf {
                    expression { return !fileExists('./code/config.yaml') }
                    expression {
                        def content = readFile(file: './code/config.yaml')
                        return content != CONFIG_FILE
                    }
                }
            }
            steps {
                dir('./code') {
                    writeFile encoding: 'utf-8', file: './config.yaml', text: CONFIG_FILE
                }
                echo "Success-创建服务配置文件"
            }
        }

        stage('启动/重启项目（删除当前端口号进程）') {
            steps {
                echo '准备重启项目-删除项目端口号'
                script {
                    try {
                        pid = sh (
                            script: "lsof -i:${config_file.Server.Port} -t",
                            returnStdout: true
                        ).trim()
                        echo 'kill......' + pid
                        sh "kill -9 ${pid} -sTCP:LISTEN"
                    } catch (Exception ex) {
                        println('ignore, server has down yet...')
                    }
                }
                echo '准备重启项目-service_router.yaml 环境变量注入'
                script {
                    def list = []
                    def read = readYaml(file: '../service_router.yaml')
                    for (element in read) {
                        echo "${element.key} ${element.value}"
                        list.add(element.key + '=' + element.value)
                    }
                    dir('./build/bin/gk') {
                        sh "ls"
                        sh "mv ${config_file.Server.Name} ../../code"
                    }
                    sh 'rm -rf ./_logs'
                    sh 'mkdir ./_logs'
                    sh list.join(' ') + " JENKINS_NODE_COOKIE=dontKillMe nohup ./code/${config_file.Server.Name} > ./_logs/${config_file.Server.Name}.log 2>&1&"
                    // 等待1s钟，看看对应端口号是否起来了，起来了代表成功 否则失败
                    sleep(1)
                    pid = sh (
                        script: "lsof -i:${config_file.Server.Port} -t",
                        returnStdout: true
                    ).trim()
                    if (pid == '') {
                        error '项目启动报错，请查询log 获取最新信息'
                    }
                    echo '项目构建成功，以后台运行！！！' + pid
                }
            }
        }
    }
}
