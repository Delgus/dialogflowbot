# Dialogflow-based bot

## Prepare Dialogflow
1) https://developers.google.com/assistant/conversational/df-asdk/dialogflow/project-agent
2) You can use prebuilt agents as Small Talk)))
3) Enable API and create account private key (json) https://cloud.google.com/dialogflow/es/docs/quick/setup

## Bot
[@dialogflowdelgusbot](https://t.me/dialogflowdelgusbot)

## Enviroments
```env
# For dialogflow
CREDENTIALS_JSON={}
PROJECT_ID=example

# For Telegram bot
TG_ACCESS_TOKEN='rrrrrr:uy54i5uy4iu5yi'
TG_WEBHOOK='example.com'

# For Telegram hook
LOG_TG_CHAT_ID=
LOG_TG_ACCESS_TOKEN=
LOG_LEVEL=

# Out port for deploy
DEPLOY_PORT=
```

## Nginx configuration example with https (tg webhook need https) Ubuntu 20

### install nginx
```
sudo apt-get install nginx
```

### register domain
I choosed freenom.com for registration my own domain


### certbot nginx ubuntu
https://certbot.eff.org/lets-encrypt/ubuntufocal-nginx

```
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot --nginx
```
```
server {
        ...

        location / {
                proxy_pass http://127.0.0.1:5100;
        }

        ...
}

```
5100 - out port for deploy. See higher DEPLOY_PORT
