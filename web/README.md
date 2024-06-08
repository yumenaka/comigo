# Comigo WebUI


````bash
# 初始化
yarn

# 启动开发服务器
yarn serve

# 为生产环境构建
yarn  build

# 目标目录不在工程项目下的时候，添加指示，清空目标目录
yarn build-vite --emptyOutDir

#本地预览生产构建
yarn preview-vite

#编译后Windows报错(text/plain;)
regegit命令打开注册表，查看HKEY_CLASSES_ROOT->.js，发现里面有一个Content-Type的配置, 将 text/plain 改成 application/javascript,然后重启本地服务器即可正常.


````

### Node.js更新到17版本后应用启动报错

````bash
export NODE_OPTIONS=--openssl-legacy-provider
#Windows新建系统变量 NODE_OPTIONS ，内容 --openssl-legacy-provider
$Env:NODE_OPTIONS = "--openssl-legacy-provider"
````

## Project setup

```bash
yarn install
```

### Compiles and hot-reloads for development

```bash  
yarn serve
```

### Compiles and minifies for production

```bash
yarn build
```

### Lints and fixes files

```bash
yarn lint
```

### 升级依赖

````bash
yarn upgrade-interactive --latest
````
