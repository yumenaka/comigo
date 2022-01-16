# comigo-webUI

## Project setup
```
yarn install
```

### Compiles and hot-reloads for development
```
yarn serve
```

### Compiles and minifies for production
```
yarn build
```

### Lints and fixes files
```
yarn lint
```

### 升级依赖
````bash
yarn upgrade-interactive --latest
````

### Customize configuration

See [Configuration Reference](https://cli.vuejs.org/config/).

### Node.js更新到17版本后应用启动报错

````bash

export NODE_OPTIONS=--openssl-legacy-provider
#Windows新建系统变量 NODE_OPTIONS ，内容 --openssl-legacy-provider
$Env:NODE_OPTIONS = "--openssl-legacy-provider"
````

### USE vite（自动转换无法热更新，以后再考虑迁移）

````bash
# Convert  project to vite by "webpack-to-vite"
## https://github.com/originjs/webpack-to-vite/blob/main/README-zh.md
yarn
# 启动开发服务器
yarn serve-vite

# 为生产环境构建
yarn build-vite

# 目标目录不在工程项目下的时候，需要添加指示，才会清空目标目录
yarn build-vite --emptyOutDir

#本地预览生产构建
yarn preview-vite
````

