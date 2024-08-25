#!/bin/bash

# 一键安装脚本
#bash <(curl -s https://raw.githubusercontent.com/yumenaka/comi/master/get_comigo.sh)

if command -v curl &> /dev/null; then
    echo "curl found"
else
    echo "curl not found, please install curl"
    exit 1
fi

if command -v tar &> /dev/null; then
    echo "check tar"
else
    echo "tar not found, please install tar"
    exit 1
fi

# 最新版本tag
latest_tag=$(curl --silent "https://api.github.com/repos/yumenaka/comi/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# 使用正则表达式提取基础版本号
if [[ $latest_tag =~ ^v([0-9]+\.[0-9]+\.[0-9]+) ]]; then
    Version=${BASH_REMATCH[1]}
else
    echo "Error: unable to parse version number"
    exit 1
fi
# 根据操作系统和处理器架构，选择下载文件
if [[ "$(uname -s)" == "Linux" ]]; then
  if [[ "$(uname -m)" == "x86_64" ]]; then
    file_name="comi_v${Version}_Linux_x86_64.tar.gz" # x86 64位
  elif [[ "$(uname -m)" == "armv7l" ]]; then
    file_name="comi_v${Version}_Linux-armv7.tar.gz"  # ARM 32位
  elif [[ "$(uname -m)" == "arm64" || "$(uname -m)" == "aarch64" ]]; then
    file_name="comi_v${Version}_Linux-armv8.tar.gz"   # ARM 64位
  else
    echo "Unsupported architecture: $(uname -m)"
    exit 1
  fi
elif [[ "$(uname -s)" == "Darwin" ]]; then
  if [[ "$(uname -m)" == "x86_64" ]]; then
    file_name="comi_${latest_tag}_MacOS_x86_64.tar.gz" # MacOS x86 64位
  elif [[ "$(uname -m)" == "arm64" ]]; then
    file_name="comi_${latest_tag}_MacOS_arm64.tar.gz" # MacOS ARM 64位
  else
    echo "Unsupported architecture: $(uname -m)"
    exit 1
  fi
else
  echo "Unsupported platform: $(uname -s)"
  exit 1
fi


# 下载文件并解压
url="https://github.com/yumenaka/comigo/releases/download/${latest_tag}/${file_name}"
echo "Downloading $url"
curl -L -O $url
tar xvf $file_name

# 清理下载文件
echo "Cleaning up $file_name"
rm $file_name

# 添加执行权限并移动到 bin 目录
echo "Adding execute permission and moving to /usr/local/bin/"
chmod +x comi
sudo mv comi /usr/local/bin/

# 获取系统语言
system_language=$(locale | grep -oP 'LANG=\K\w+')

# 提示用户
if [ "$system_language" == "en_US" ]; then
    echo -e "\033[34mComigo is installed, you can run the comi command in the comics directory to scan for comics\033[0m"
elif [ "$system_language" == "zh_CN" ]; then
    echo -e "\033[34mcomigo 安装完毕，可以在漫画目录下执行 comi 命令扫描漫画了\033[0m"
elif [ "$system_language" == "ja_JP" ]; then
    echo -e "\033[34mComigoがインストールされています。 これで、ディレクトリでcomiコマンドを実行し、コミックをスキャンできるようになりました\033[0m"
else
    echo -e "\033[34mComigo is installed, you can run the comi command in the comics directory to scan for comics\033[0m"
fi


