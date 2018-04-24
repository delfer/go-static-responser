# Go Static Responser
[![Docker Stars](https://img.shields.io/docker/stars/delfer/go-static-responser.svg)](https://hub.docker.com/r/delfer/go-static-responser/) [![Docker Pulls](https://img.shields.io/docker/pulls/delfer/go-static-responser.svg)](https://hub.docker.com/r/delfer/go-static-responser/) [![Docker Automated build](https://img.shields.io/docker/automated/delfer/go-static-responser.svg)](https://hub.docker.com/r/delfer/go-static-responser/) [![Docker Build Status](https://img.shields.io/docker/build/delfer/go-static-responser.svg)](https://hub.docker.com/r/delfer/go-static-responser/) [![MicroBadger Layers](https://img.shields.io/microbadger/layers/delfer/go-static-responser.svg)](https://hub.docker.com/r/delfer/go-static-responser/) [![MicroBadger Size](https://img.shields.io/microbadger/image-size/delfer/go-static-responser.svg)](https://hub.docker.com/r/delfer/go-static-responser/)

Send static HTTP response fast and log visits to ClickHouse 

## Features

- Response very fast to `/` URI with static string
- Log every visit to Yandex ClickHouse
- Configurable by environment variables
- Shows buffer size on `/load`

## Configuration

Environment variables:
- `PORT` - HTTP listen port (8080 by default)
- `RESPONSE` - string which response to `GET /`
- `BUFFER` - buffer size (in requests) between HTTP server and DB writer (100,000 bн default)
- `DISABLE_CH` - set `true` to disable writing to ClickHouse
- ClickHouse connection
  - `CH_HOST` - host (127.0.0.1 by default)
  - `CH_PORT` - port (9000 by default)
  - `CH_DEBUG` - debug enabled true/false (false by default)
  - `CH_USER` - user (empty=default by default)
  - `CH_PASSWORD` - password (nothing by default)
  - `CH_DB` - database (empty=default by default)

## Usage

```
docker run -d --restart always \
    -e RESPONSE="Hello!" \
    -e CH_HOST=10.0.0.1 \
    -e CH_PASSWORD="password" \
    -p 8080:8080 delfer/go-static-responser
```
Open http://10.0.0.1/ in you browser or by `curl` to make new visit,
open http://10.0.0.1/load to get current buffer usage

## License

MIT
