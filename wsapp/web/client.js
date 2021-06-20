window.onload = function () {
    let wsSchema;

    if (document.location.protocol === "http:") {
        wsSchema = "ws:";
    } else if (document.location.protocol === "https:") {
        wsSchema = "wss:";
    }

    let ws = new WebSocket(wsSchema + "//" + document.location.host + "/entry");

    let chatbox = document.getElementById("chatbox");

    ws.onerror = () => alert("WEBSOCKET SERVER DOESN'T WORK!")

    ws.onclose = () => alert("CONNECTION CLOSED")

    ws.onmessage = (e) => {
        let msg = JSON.parse(e.data);

        if (msg.type == "connect") {
            return
        }

        chatbox.innerHTML += "client_" + msg.user_id + ": " + msg.body + "<br>";
        chatbox.scrollTop = 9999;
    };

    // отправка сообщений на вебсокет
    let form = document.querySelector('form');
    form.onsubmit = function () {
        if (form[0].value !== '') {
            ws.send(JSON.stringify({ body: form[0].value }));
        }
        form[0].value = '';
        return false;
    };
};