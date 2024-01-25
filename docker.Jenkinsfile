// 注意：build httpCore的时候放到最后，因为要把apis 更新了，先build 其他会更新
def config_file = ''
def need_push = false
def first_clone = true // httpCore 每次都需要build 因为不知道其他项目的apps是否有变化
def CONFIG_FILE = '''# 基础环境配置
Server:
    Name: http_core
    Port: 8885
    ApiPath: [ ../user/apis,  ../assets/apis, ../share/apis, ../check/apis, ../mete_parse/apis]
    StoragePath: ./_storage
    ThirdToken: ${THIRD_TOKEN}

# 系统内鉴权类型
AuthAll:
    - # 后台管理者用户
     roleName: Admin
     ExpiredDuration: 120
    - # core将重新颁发证书-先调用微服务返回成功后刷新（常用于登陆，身份信息变更刷新）
     RoleName: Offer  # 可以和其他权限组合使用
     ExpiredDuration: -1
    - # 为公共方法不鉴权（本系统内接口）
     RoleName: Public  # 只能单独出现，不能和其他组合使用
     ExpiredDuration: -1
    - # 为系统全体已登陆用户
     RoleName: Protected  # 只能单独出现，不能和其他组合使用
     ExpiredDuration: -1

Jwt:
     Iss: ${JWT_ISS}
     Secret: ${JWT_SECRET}
'''

pipeline {
    agent {
        node {
          label 'iot-dev'
          customWorkspace './workspace/docker-jianzai/httpCore'
        }
    }

    stages {
        stage('@@@-初始化pipline config_file---server_name, port from CONFIG_FILE') {
            steps {
                script {
                    config_file = readYaml text: CONFIG_FILE
                    echo config_file.Server.Name
                    echo "${config_file.Server.Port}"
                    dir ('./httpCore') {
                        try {
                            sh 'git checkout .'
                        } catch (Exception ex) {
                            first_clone = true
                            println('ignore, before git clone...')
                        }
                    }
                }
            }
        }

        stage('@@@-拉取代码 testing') {
            steps {
                echo '准备拉取代码-httpCore + comlibgo'
                checkout([$class: 'GitSCM', branches: [[name: '*/jianzai']], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: './httpCore']], userRemoteConfigs: [[url: 'git@github.com:techzhongyi/httpCore.git']]])
                checkout([$class: 'GitSCM', branches: [[name: '*/jianzai']], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: './comlibgo']], userRemoteConfigs: [[url: 'git@github.com:techzhongyi/comlibgo.git']]])
            }
        }

        stage('@@@-准备各个项目的apis') {
            steps {
                sh 'mkdir -p ./httpCore/__all_apis/user/apis && cp ../user/apis/* ./httpCore/__all_apis/user/apis/'
                sh 'mkdir -p ./httpCore/__all_apis/assets/apis && cp ../assets/apis/* ./httpCore/__all_apis/assets/apis/'
                sh 'mkdir -p ./httpCore/__all_apis/share/apis && cp ../share/apis/* ./httpCore/__all_apis/share/apis/'
                sh 'mkdir -p ./httpCore/__all_apis/check/apis && cp ../check/apis/* ./httpCore/__all_apis/check/apis/'
                sh 'mkdir -p ./httpCore/__all_apis/mete_parse/apis && cp ../mete_parse/apis/* ./httpCore/__all_apis/mete_parse/apis/'
            }
        }

        stage('@@@-创建服务配置文件-config.yaml') {
            // todo 判断config_file文件和本地比较是否有变化
            when {
                anyOf {
                    expression { return !fileExists('./httpCore/config.yaml') }
                    expression {
                        def content = readFile(file: './httpCore/config.yaml')
                        return content != CONFIG_FILE
                    }
                }
            }
            steps {
                writeFile encoding: 'utf-8', file: './httpCore/config.yaml', text: CONFIG_FILE
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
                sh 'cp ./httpCore/Dockerfile .'
                sh "docker build --tag registry.cn-beijing.aliyuncs.com/zy_jz/${config_file.Server.Name}:latest --no-cache ."
            }
        }

        stage('@@@-Push image') {
            when {
                expression { return need_push }
            }
            steps {
                sh 'echo xdkj@6789 | docker login -u 新的空间 --password-stdin registry.cn-beijing.aliyuncs.com'
                sh "docker push registry.cn-beijing.aliyuncs.com/zy_jz/${config_file.Server.Name}:latest"
                echo '上传镜像成功！'
                sh 'rm -rf ./httpCore/__all_apis'
            }
        }
    }
}
