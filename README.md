# clipboard-online

clipboard-online is an application to share cilpboard text between Windows and iOS

## Download

1. Directly download

    Ahho, not complete...

2. From source code(only windows now)

    Before you build, make sure you have installed golang. If not, maybe you need [this](https://golang.org/dl/)
    - `git clone git@github.com:YanxinTang/clipboard-online.git`
    - `cd clipboard-online`
    - `./build/build.bat`
    - You can find release bin at `release` directory

## Usage

1. Run `clipboard-online` on your windows
2. Setup shortcuts on you iPhone/iPad
3. Have fun...

## API

The default http server will listen `8000` port and you can't chanage that since hardcoded.

### 1. Get windows clipboard

> Request

- URL: `/clipboard`
- Method: `GET`

> Reponse

- Body: `<clipboard text>`

### 2. Set windows clipboard

> Request

- URL: `/clipboard`
- Method: `POST`
- Body: `text you want to set`

> Reponse

Reponse body is empty. If set successfully, status code will be `200`
