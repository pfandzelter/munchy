# munchy

SlackBot that tells you what's for lunch at TU Berlin today. Uses [go-eat](https://github.com/pfandzelter/go-eat) as a backend. Runs on AWS Lambda.

```sh
aws configure
make
```

You will need AWS access keys and an AWS region where you'd like to deploy this. Also, you need a URL for Slack Webhooks to go to. You will need a DeepL API key to use translation. To disable translation set `deepl_target_lang` to `DE`.

![Gopher](https://random.pfandzelter.com/icon.png)
