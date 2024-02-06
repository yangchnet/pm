# pm (password manager) - 简单实用的密码管理工具

pm 是一个简单实用的密码管理工具,它可以帮助你生成、管理密码，可选的将密码同步到远端存储中，以便你可以在多个主机间同步你的密码.

## Build

```bash
make
```

## Install
如果你安装了go,可以直接使用go install安装:
```bash
go install github.com/yangchnet/pm@latest
```

或直接从release下载二进制文件

## Usage

![pass](https://raw.githubusercontent.com/lich-Img/blogImg/master/img/pass.gif)

1. 初始化pm
```bash
pm init
```

可选参数：
- `--remote` `-r`. 指定使用的远端存储, 可选项：`git`、`empty`，默认empty。当使用`git`时,需要输入git配置信息；当使用`empty`时，pm将不会进行远端存储。
- `--store` `-s`. 指定使用的本地存储, 可选项：`sqlite`、`file`，默认file。当使用`sqlite`时,pm将使用sqlite数据库存储密码；当使用`file`时，pm将使用文本文件进行密码存储。

2. 录入一个新的密码
```bash
pm new passname # passname为你的密码名称, 必须唯一
```

可选参数：
- `--passwd` `-p`. 指定密码，如果不指定，pm将会生成一个密码。
- `--note` `-n`. 指定备注信息。
- `--url` `-u`. 指定url信息。
- `--account` `-a`. 指定账号信息。

当未设置主密码或主密码已失效时，pm将会提示你输入主密码。主密码强制每24小时必须重新输入一次。

3. 获取密码
```bash
pm get passname
```

4. 删除密码
```bash
pm del passname
```


5. push & pull
```bash
pm push # 将本地密码同步到远端

pm pull # 将远端密码同步到本地
```

正常情况下，pm会在你对密码进行了改动后自动将数据同步到远端，并在你get密码前自动同步远端密码到本地。

Enjoy it!