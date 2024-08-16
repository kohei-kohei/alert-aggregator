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

You need to pass the `Bot User OAuth Token` to `env.SLACK_BOT_TOKEN` and store it in the repository's Secrets.

## Inputs

### `get-channel-id`

This field is required. The Slack channel ID where you want to aggregate alerts.

How to find the ID: https://slack.com/help/articles/221769328-Locate-your-Slack-URL-or-ID

### `send-channel-id`

The Slack channel ID where you want to send the aggregated alert results. The Slack App must also be added to this destination channel. If not specified, the results won't be sent anywhere.

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
    with:
      get-channel-id: 'YOUR_CHANNEL_ID'
```

You can also specify a time zone.

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
      TZ: 'Asia/Tokyo'
    with:
      get-channel-id: 'YOUR_CHANNEL_ID'
```


## Aggregating alerts and sending them to a specified channel

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
    with:
      get-channel-id: 'YOUR_CHANNEL_ID'
      send-channel-id: 'YOUR_CHANNEL_ID'
```

## Specifying the period to aggregate alerts

```yaml
steps:
  - uses: kohei-kohei/alert-aggregator@v0
    env:
      SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
    with:
      get-channel-id: 'YOUR_CHANNEL_ID'
      since: '2024-07-01T00:00:00+09:00'
      until: '2024-08-01T00:00:00+09:00'
```

# License

The scripts and documentation in this project are released under the [MIT License](LICENSE)
