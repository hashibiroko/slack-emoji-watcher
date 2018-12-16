# slack-emoji-watcher

## Install

require Golang environment and setup GOROOT.

```
$ go get github.com/hashibiroko/slack-emoji-watcher
```

## Usage

Please invite a bot in advance to the channel you want to notify.

#### Example 1:

```
$ slack-emoji-watcher -token=xxxxxx-xxxxxxxxx -channel=watcher
```

#### Example 2: setting environment

```
$ export SLACK_BOT_TOKEN="xxxxxx-xxxxxxxxx"
$ export SLACK_CHANNEL_NAME="watcher"
$ slack-emoji-watcher
```

### Flags

| name | description | default | require | environment |
| :--- | :---------- | :-----: | :-----: | :---------- |
| token | Set your slack bot token |  | true | SLACK_BOT_TOKEN |
| channel | Set your slack notification channel | random |  | SLACK_CHANNEL_NAME |
