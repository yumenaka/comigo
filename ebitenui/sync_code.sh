#!/bin/bash
# 本地代码同步到远程服务器
echo "sync code start"
rsync  -avz -e ssh\
  --exclude '.air.toml'\
  --exclude 'sync_code.sh'\
  --exclude 'bin'\
  --exclude 'test'\
  --exclude 'upload'\
  --exclude 'temp' \
  --exclude '.git'\
  --exclude 'node_modules'\
  --exclude '.vscode'\
  --exclude '.vite' \
  --exclude '.idea'\
  --exclude '.dart_tool'\
  --exclude 'build'\
  --exclude 'tmp'\
  --exclude '.DS_Store'\
  /home/user/comigo/ yume@some-machine.localhost:/home/user/comigo/
echo "sync code done"

# 忽略已经提交的文件 (.gitignore文件无效)
# git update-index --assume-unchanged sync_code.sh
# git update-index --assume-unchanged .air.toml

# 取消忽略
# git update-index --no-assume-unchanged sync_code.sh
# git update-index --no-assume-unchanged .air.toml