# munchy

SlackBot that tells you what's for lunch at TU Berlin today. Uses [go-eat](https://github.com/pfandzelter/go-eat) as a backend. Runs on AWS Lambda. 
```
$ aws configure
$ make
```

You will need AWS access keys and an AWS region where you'd like to deploy this. Also, you need an URL for Slack Webhooks to go to.

![Gopher](https://random.pfandzelter.com/icon.png)