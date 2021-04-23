# Access-Token-Server-Samples

This repository contains:

- Sample access token servers
- Sample clients that can fetch access tokens from the token servers

## Architecture

### Access token authentication for Agora RTC SDK

![](https://web-cdn.agora.io/docs-files/1608020494998)

### Access token authentication for Agora RTM SDK

Work in progress

## Servers

### Golang

Token server sample based on Golang. You can deploy this server locally, on a virtual machine or cloud.

#### How to run

1. Install dependencies.

```shell
go get
```
> For users in China, GOPROXY is recommended. See https://goproxy.cn/.

2. Run the project.

```shell
go run server.go
```

## Clients

### RTC Web SDK 4.x.

Web client based on RTC Web SDK 4.x. You can use this client to fetch token from the sample server.

#### How to run

Open `index.html`.
