# Go static responder

Send static HTTP response fast and log visits to ClickHouse 

## Features

- Response very fast to `/` URI with static string
- Log every visit to Yandex ClickHouse
- Configurable by environment variables
- Shows buffer size on `/load`

## Configuration

- `PORT` - HTTP listen port (8080 by default)
- `RESPONSE` - string which response to `GET /`
- `BUFFER` - buffer size (in requests) between HTTP server and DB writer (100,000 bн default)
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
