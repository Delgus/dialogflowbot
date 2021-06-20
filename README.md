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

## Nginx configuration example with https
```
server {
        listen 80;
        listen 443 ssl http2;
        server_name dialogflow.example.com;
        ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;
        ssl_trusted_certificate /etc/letsencrypt/live/example.com/chain.pem;

        ssl_stapling on;
        ssl_stapling_verify on;
        resolver 127.0.0.1 8.8.8.8;

        location / {
                proxy_pass http://127.0.0.1:5100;
        }

        location /.well-known {
                  root /var/www/letsencrypt;
        }
}

```
5100 - out port for deploy. See higher DEPLOY_PORT
