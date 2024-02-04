# pm (password manager) - 简单实用的密码管理工具

## Build

```bash
make
```

## Install
pm有一些依赖,请自行安装:
```
sudo apt-get install git xclip
```

## Usage

1. 在任意的git-server上创建一个私有仓库(github, gitlab, gitee等),并复制仓库的ssh地址:例如:`git@github.com/yangchnet/pass.git`

2. 在本地初始化pm,并指定git的用户名,邮箱,仓库地址,以及私钥的路径
```bash
pm init --git-email example@qq.com --git-username yangchnet --git-url git@github.com:yangchnet/pass.git --private-key-path /path/to/private/key
```
如果还没有私钥,可不传`--private-key-path`参数,pm会自动创建一个新的私钥

3. 将刚才创建的私钥对应的公钥添加到git的ssh-key中,以便pm可以自动push和pull(如果之前配置过ssh密钥,可参考第二步直接使用)

![20240204190855](https://raw.githubusercontent.com/lich-Img/blogImg/master/img/20240204190855.png)
![20240204191010](https://raw.githubusercontent.com/lich-Img/blogImg/master/img/20240204191010.png)

4. 拉取远程仓库的数据
```bash
pm pull
```

5. 新建一个密码
```bash
pm new --name github --account yangchnet
```

可选参数:`--note`, `--url`

pm会自动为你生成一个密码,并将密码复制到你的剪贴板中,密码被保存在sqlite数据库中.

6. push到远程仓库保存
```bash
pm push
```

## TODO
[ ] 数据导入

[ ] windows环境测试

[ ] 支持自定义密码参数(长度,包含字符等)

[ ] OSS作为remote

[ ] ...