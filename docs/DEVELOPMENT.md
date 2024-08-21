# Development

## Setup

```shell
docker build . -t alert-aggregator
```

## Run

Please put appropriate values in the environment variables.

```shell
docker run --rm -e SLACK_BOT_TOKEN="xxxx" -e GET_CHANNEL_ID="xxxx" -e SEND_CHANNEL_ID="" -e Since="" -e Until="" alert-aggregator
```
