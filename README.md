# Dialogflow-based bot

## Demo

### Telegram

[@dialogflowdelgusbot](https://t.me/dialogflowdelgusbot)

### VKBot

[vkbot](https://vk.com/im?sel=-205333881)

### Web (websocket-based)

[web chat](https://delgus.tk)

## How use?


### Prepare Dialogflow
1) https://developers.google.com/assistant/conversational/df-asdk/dialogflow/project-agent
2) You can use prebuilt agents as Small Talk)))
3) Enable API and create account private key (json) https://cloud.google.com/dialogflow/es/docs/quick/setup

### register and configurate domain
I choosed freenom.com for registration my own domain

### configurate server (ex. certbot, nginx, ubuntu)

```
sudo apt-get install nginx
```

https://certbot.eff.org/lets-encrypt/ubuntufocal-nginx

```
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
sudo certbot --nginx
```

### configurate proxy

```
server {
        ...

        location / {
                proxy_pass http://127.0.0.1:5100;
                # for websockets in aws (important!)
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "Upgrade";
                proxy_set_header Host $host;
        }

        ...
}

```

### enviroments

dialogflow.env

```env
# dilalogflow account
CREDENTIALS_JSON={"type": "service_account" ...}
PROJECT_ID=uwytiwuyet-g6566

# telegram
TG_ACCESS_TOKEN=2000000000:YRUIYIUYIUYIUWYTIUTWUIYTWIUTW
TG_WEBHOOK=https://mydomain.tk/tg

# vk
VK_ACCESS_TOKEN=e5b078c3285f44f08fed8e70f4d2eb3e9b3d38d7aa5e91e4b9deff00194e49e31621e6bfc518bc4d2b275
VK_WEBHOOK=http://mmydomain.tk/vk
VK_CONFIRM_KEY=fc0d8070

# websockets
WS_URL=wss://delgus.tk/entry

```

### run with docker

```
docker run --name dfb -p 5100:80 --env-file dialogflow.env --restart unless-stopped -d delgus/dialogflowbot
```



