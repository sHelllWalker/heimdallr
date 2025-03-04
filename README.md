# heimdallr
Simple go application for integrating Sentry with Mattermost and Telegram

## Configuration:
[See all enviroments](https://github.com/sHelllWalker/heimdallr/blob/main/internal/config/config.go)

## How to use
1. Create MM webhook and set it in WEBHOOK_URL env
2. If you need TG - create bot via BotFather and set its token in TOKEN env
3. Create sentry webhook with url: `{{your_domain}}/broadcast?channel={{set here target channel name}}&chatId={{set here needed tg chat id}}`

If the basic templates are not enough for you, specify the paths to the files of the templates you are overriding for each messenger and event type from the sentry (*TemplatePath envs)
