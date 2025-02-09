let ws = new WebSocket("ws://" + window.location.host + "/ws");
ws.onmessage = function(event) {
    let chat = document.getElementById("chat");
    let message = document.createElement("div");
    message.textContent = event.data;
    chat.appendChild(message);
    chat.scrollTop = chat.scrollHeight;
};

let input = document.getElementById("input");
input.addEventListener("keyup", function(event) {
    if (event.keyCode === 13) {
        ws.send(input.value);
        input.value = "";
    }
});
