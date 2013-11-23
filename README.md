# circle2slack

Proxy [CircleCI](http://circleci.com) webhooks to [Slack](http://slack.com)

    heroku create -b https://github.com/kr/heroku-buildpack-go.git
    heroku config:add SLACK_BOTNAME=circleci SLACK_CHANNEL=#code SLACK_ORGANIZATION=org SLACK_TOKEN=xxx

`circle.yml`:

```yml
notify:
  webhooks:
    - url: https://circle2slack-app.herokuapp.com/build/
```
