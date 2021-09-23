#!/usr/bin/env bash


URL="http://172.16.0.101/frigate"

# name will be generate like 'frigate_v${VERSION}_${OS}${EXT}'
VERSION="1.0.2"
EXT=".tar.gz"
OS=""

function install() {
    name="frigate_v${VERSION}_${OS}"
    path="${URL}/${name}${EXT}"

    rm ${name}${EXT}
    rm ${name}

    echo "download ${path} ..."
    curl -O ${path}

    tar xvf ${name}${EXT}
    rm ${name}${EXT}

    echo "move frigate to /usr/local/bin"
    mv ${name} /usr/local/bin/frigate

    if [$? != 0]; then
        mv ${name} frigate
        echo -e "\033[31m Please add frigate to path by yourself"
    else
        echo -e "\033[32mInstall successfully, Please use the command below to initial your test directory: \033[0m "
        echo -e "\033[32m\t mkdir test && cd test\033[0m"
        echo -e "\033[32m\t frigate init \033[0m"
    fi
}

function getLinuxReleaseType(){
    release=$(cat /etc/*release*)
    case ${release} in
        *"CentOS Linux release 7"*)
            echo "CentOS7"
            OS="centos7"
            ;;
        *)
            echo "unsupported Linux release: ${release}"
            exit 1
            ;;
    esac
}

function getOSType(){
    case "$(uname)" in
        "Darwin")
            echo "Darwin"
            OS="darwin"
            ;;
        "Linux")
            echo "Linux"
            getLinuxReleaseType
            ;;
        *)
            echo "Unsupported now"
            ;;
    esac
}

getOSType
install




