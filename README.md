<div align="center">
  <img src="https://raw.githubusercontent.com/YanxinTang/clipboard-online/master/images/clipboard-icon.png" style="display: inline-block; vertical-align: middle;">
  <h1 style="display: inline-block; vertical-align: middle;">clipboard-online</h1>
</div>

![GitHub release (latest by date)](https://img.shields.io/github/v/release/YanxinTang/clipboard-online)

clipboard-online is an application to share cilpboard text between Windows and iOS

## Documentation

【[中文](https://github.com/YanxinTang/clipboard-online/blob/master/README_zh.md)】【[English](https://github.com/YanxinTang/clipboard-online/blob/master/README.md)】

## Download

1. Directly download

    You can download latest release exe from [here](https://github.com/YanxinTang/clipboard-online/releases)

2. From source code(only windows now)

    Before you build, make sure you have installed golang. If not, maybe you need [this](https://golang.org/dl/)

    - `git clone git@github.com:YanxinTang/clipboard-online.git`
    - `cd clipboard-online`
    - `go get github.com/akavel/rsrc`
    - `./build.sh`
    - You can find release bin at `release` directory

## Usage

1. Run `clipboard-online` on your windows
2. Setup shortcuts on you iPhone/iPad (Open link from safari)
    - Copy: [https://www.icloud.com/shortcuts/242c55e0895e4235875bc71f1f010199](https://www.icloud.com/shortcuts/242c55e0895e4235875bc71f1f010199)
    - Paste: [https://www.icloud.com/shortcuts/6a46febf2f0c4ef4b00bbc41f03ccd2f](https://www.icloud.com/shortcuts/6a46febf2f0c4ef4b00bbc41f03ccd2f)
3. Have fun...😊

## Configuration

`clipboard-online.exe` will create two file which are `config.json` and `log.txt` in the execute path when first running

You can make customization by editing `config.json`

### `config.json`

- `port`
  - type: `string`
  - default: `"8086"`

- `logLevel`
  - type: `string`
  - default: `"warning"`
  - values: `"panic"`, `"fatal"`, `"error"`, `"warning"`, `"info"`, `"debug"`, `"trace"`

## API

The default http server will listen `8086` port and you can't chanage that since hardcoded.

### 1. Get windows clipboard

> Request

- URL: `/`
- Method: `GET`

> Reponse

- Body: `<clipboard text>`

### 2. Set windows clipboard

> Request

- URL: `/`
- Method: `POST`
- Body: `text you want to set`

> Reponse

Reponse body is empty. If set successfully, status code will be `200`
