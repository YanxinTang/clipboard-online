<div align="center">
  <img src="https://raw.githubusercontent.com/YanxinTang/clipboard-online/master/images/clipboard-icon.png" style="display: inline-block; vertical-align: middle;">
  <h1 style="display: inline-block; vertical-align: middle;">clipboard-online</h1>
</div>

![GitHub release (latest by date)](https://img.shields.io/github/v/release/YanxinTang/clipboard-online)

clipboard-online 是一款可以帮你在 💻Windows 和 📱iOS 之间分享剪切板的应用

## 文档

【[English](https://github.com/YanxinTang/clipboard-online/blob/master/README.md)】【[中文](https://github.com/YanxinTang/clipboard-online/blob/master/README_zh.md)】

## 下载

1. 直接下载

    在[这里](https://github.com/YanxinTang/clipboard-online/releases)下载发布的 .exe 文件

2. 源码编译(只在 Windows 下可用，其他平台未知)

    构建之前，请确保你已经安装了 golang. 如果没有，可能你需要[这个](https://golang.org/dl/)

    - `git clone git@github.com:YanxinTang/clipboard-online.git`
    - `cd clipboard-online`
    - `go get github.com/akavel/rsrc`
    - 构建
      - Bash: `./build.sh`
      - PowerShell：`.\build.ps1`
    - 你可以在 `release` 目录下找到可执行文件

## 使用

### iOS 用户

1. 在 Windows 上运行 `clipboard-online`
2. iPhone 或 iPad 上安装快捷指令 （在 safari 中打开链接）
    - Copy:
      - iCloud: [https://www.icloud.com/shortcuts/f463a1e431c94c60b8a5c65305eb819f](https://www.icloud.com/shortcuts/f463a1e431c94c60b8a5c65305eb819f)
      - ![复制](./images/copy.png)
    - Paste:
      - iCloud: [https://www.icloud.com/shortcuts/90e7a2af70df4707a17dece8c263afc5](https://www.icloud.com/shortcuts/90e7a2af70df4707a17dece8c263afc5)
      - ![粘贴](./images/paste.png)

3. 设置 ip 地址和 authkey （默认是空字符串）
4. 玩的开心...😊

### Android 用户

1. HTTP Shortcuts: [https://meta.appinn.net/t/topic/20322](https://meta.appinn.net/t/topic/20322)

2. Tasker: [https://github.com/YanxinTang/clipboard-online/issues/26](https://github.com/YanxinTang/clipboard-online/issues/26)

非常感谢 [@xiozf](https://github.com/xiozf)，[@Befod](https://github.com/Befod) 及其他贡献者，谢谢你们为 Android 平台提供的支持！

## 配置

`clipboard-online.exe` 将在运行路径下面创建两个文件： `config.json` and `log.txt`

你可以通过修改 `config.json` 来自定义配置

### `config.json`

- `port`
  - 类型: `string`
  - 默认: `"8086"`

- `logLevel`
  - 类型: `string`
  - 默认: `"warning"`
  - 可选: `"panic"`, `"fatal"`, `"error"`, `"warning"`, `"info"`, `"debug"`, `"trace"`

- `authkey`
  - type: `string`
  - default: `''`

- `authkeyExpiredTimeout`
  - type: `int64`
  - default: `30`

- `tempDir`
  - type: `string`
  - default: `./temp`

- `reserveHistory`
  - type: `Boolean`
  - default: `false`

- `notify`
  - type: `object`
  - children:
    - `copy`
      - type: `Bollean`
      - default: `false`
    - `paste`
      - type: `Boolean`
      - default: `false`

## API

### 公共 headers

#### 必选

- `X-API-Version`: indicates version of api

#### 可选

- `X-Client-Name`: indicates name of device
- `X-Auth`: hashed authkey. Value from `md5(config.authkey + timestamp/30)`

### 1. 获取 Windows 剪切板

> Request

- URL: `/`
- Method: `GET`

> Reponse

- Body: `json`

```json
// 200 ok

{
  "type": "text",
  "data": "clipboard text on the server"
}

{
  "type": "file",
  "data": [
    {
      "name": "filename",
      "content": "base64 string of file bytes"
    }
    ...
  ]
}

```

### 2. 设置 Windows 剪切板

> Request

- URL: `/`
- Method: `POST`
- Headers:
  - `X-Content-Type`: indicates type of request body content
    - `required`
    - values: `text`, `file`, `media`

- Body: `json`

For text:

```json
{
  "data": "text you want to set"
}
```

For file:

```json
{
  "data": [
    {
      "name": "filename",
      "base64": "base64 string of file bytes"
    }
  ]
}
```

响应的 body 为空。如果剪切板设置成功，状态码将返回 `200`
