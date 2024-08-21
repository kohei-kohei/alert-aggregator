# alert-aggregator

This action roughly aggregates alerts that are being sent to Slack. The following alerts from the respective bots are aggregated:

- Datadog
- Sentry
- PagerDuty
- AWS Chatbot
- digdag-alert

# Usage

This action requires a [Slack App](https://api.slack.com/quickstart) with [conversations.history](https://api.slack.com/methods/conversations.history) permissions, which must be added to the target channel. Please note that `conversations.history` is limited to 999 messages without pagination, so longer periods may miss alerts.

To send results to Slack, the Slack App needs [chat.postMessage](https://api.slack.com/methods/chat.postMessage) permission.

## Environment Variables

`SLACK_BOT_TOKEN` and `GET_CHANNEL_ID` are required.

### `SLACK_BOT_TOKEN`

The `Bot User OAuth Token` for the created slack app.

### `GET_CHANNEL_ID`

The Slack channel ID where you want to aggregate alerts.

How to find the ID: https://slack.com/help/articles/221769328-Locate-your-Slack-URL-or-ID

### `SEND_CHANNEL_ID`

The Slack channel ID where you want to send the aggregated alert results. The Slack App must also be added to this destination channel. If not specified, the results won't be sent anywhere.

## Inputs

### `since`

The start date and time for alert aggregation. It must follow RFC3339 format. The default value is `00:00 8 days ago` based on the server's local time.

### `until`

The end date and time for alert aggregation. It must follow RFC3339 format. The default value is `00:00 yesterday` based on the server's local time.

# Examples

## Just aggregating alerts

By default, it aggregates 1 week of alerts in UTC.

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
      GET_CHANNEL_ID: ${{ secrets.GET_CHANNEL_ID }}
```

You can also specify a time zone.

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
      GET_CHANNEL_ID: ${{ secrets.GET_CHANNEL_ID }}
      TZ: 'Asia/Tokyo'
```


## Aggregating alerts and sending them to a specified channel

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
      GET_CHANNEL_ID: ${{ secrets.GET_CHANNEL_ID }}
      SEND_CHANNEL_ID: ${{ secrets.SEND_CHANNEL_ID }}
```

## Specifying the period to aggregate alerts

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
      GET_CHANNEL_ID: ${{ secrets.GET_CHANNEL_ID }}
    with:
      since: '2024-07-01T00:00:00+09:00'
      until: '2024-08-01T00:00:00+09:00'
```

# License

The scripts and documentation in this project are released under the [MIT License](LICENSE)
