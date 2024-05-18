#!/bin/bash

rsync \
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
  -e 'ssh -p 22' -avP yume@some-machine.localhost:/home/user/comigo/ /home/user/comigo/