# Dingo Tiles - Ask for Help: Slackbot

This is the source code for the `dingobot` in our http://slack.dingotiles.com/ outback land of Support & Sales help.

In any channel type `/download`:

![](http://cl.ly/3y0b3a113n0L/download/Image%202016-02-18%20at%203.08.48%20pm.png)

## Deployment

```
cf push --no-start
cf set-env slackbot DOWNLOAD_SLACK_TOKEN token
cf set-env slackbot DINGOTILES_ASKFORHELP_IN_URL https://hooks.slack.com/services/xxx/yyy/zzz
cf set-env slackbot AWS_ACCESS_KEY access
cf set-env slackbot AWS_SECRET_KEY secret
cf restart slackbot
```
