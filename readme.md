# NightWriter

[![GitHub stars](https://img.shields.io/github/stars/Titanexx/NightWriter)](https://github.com/Titanexx/NightWriter/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/Titanexx/NightWriter)](https://github.com/Titanexx/NightWriter/network)

<p align="center">
<img width="150" alt="logo" src="https://raw.githubusercontent.com/Titanexx/NightWriter/master/images/logo.png">
</p>

> :warning: **It's an alpha release, use it at your own risk**

# Another creation document tool ?

After many years of searching for a good document tool that is simple, secure and modern, I decided to develop it.
NightWriter offers simple, secured by end-to-end encryption and real-time collaborative editing.

Markdown is the markup language. As NightWriter exports your documents in markdown, you can easily convert them in many formats.
For example, if you want convert into a docx file, you can use my other project [MDdot](https://github.com/Titanexx/MDdot)

## Security

The end to end encryption is based on RSA and AES 256.
When you create your account, a keypair is generated and the stored private key is encrypted by your master password.
Documents are encrypted by AES 256 and the key is encrypted with the public key of the user or the group before to be send to the server.

The server never knows the passwords or even the content of your documents.

The cryptography implementation is done by [Forge](https://github.com/digitalbazaar/forge) which is widely used.

# Open Source Components

NightWriter is built on [Go](https://golang.org/), [Gin](https://github.com/gin-gonic/gin) and [GORM](https://gorm.io/index.html) for the backend and [Vue3](https://v3.vuejs.org/), [Tui.Editor](https://github.com/nhn/tui.editor) and [tailwindcss](https://tailwindcss.com/) for the front.

# How to test it locally ?
## The manual way
### Prequesites

- Go (in 1.16.x)
- Node (in 16.x) & npm (in 7.x)

### Step

Clone this repository and execute:
```
cd NightWriter
docker-compose up -d
cd front
npm install .
npm run build
cd ..
cp .env.example .env
go run .
```

Enjoy !

# How to install on a server ?

It's coming.

# Licenses

NightWriter is dual-licensed - Copyright 2021 Titanex

If you want commercialize service which include NightWriter, you need a commercial, non-free license.
If you want a direct support or specific development, you need a non-free license.
Otherwise, NightWriter can be used without charge under the terms of GPLv3.

# Support

## How support ?

You can contribute to the project. There is still a lot to do and I have many ideas for the future.
You can support by buying :coffee:, :tea: or :beers: [https://www.buymeacoffee.com/titanex](https://www.buymeacoffee.com/titanex)

Other support ways are coming.

## Why support the tool ?

As I used many awesome open source project, I want support some of them.
As the project requires a significant investment of time, even a simple coffee will help me to continue the development.

