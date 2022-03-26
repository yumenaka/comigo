# comigo-webUI



### 目前已经迁移到了vite

````bash
# 初始化
yarn

# 启动开发服务器
yarn vite
# 或
yarn serve-vite

# 为生产环境构建
yarn vite build
# 或
yarn build-vite

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

### 

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

