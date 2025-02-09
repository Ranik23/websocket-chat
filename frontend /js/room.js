let ws, clientsWs;

document.getElementById("joinRoomBtn").addEventListener("click", function () {
    const name = document.getElementById("nameInput").value.trim();
    const room = document.getElementById("roomInput").value.trim();
    if (!name || !room) return;

    ws = new WebSocket("ws://" + window.location.host + "/ws?room=" + encodeURIComponent(room) + "&name=" + encodeURIComponent(name));

    ws.onopen = function () {
        console.log("Подключено к комнате: " + room);
        document.getElementById("room-selection").style.display = "none";
        document.getElementById("chat-section").style.display = "block";
        document.getElementById("roomName").textContent = room;
        connectClientsWebSocket(room);
    };

    ws.onmessage = function (event) {
        try {
            const data = JSON.parse(event.data);
            const sender = data.sender || 'Неизвестный';
            const text = data.text || '';

            const chat = document.getElementById("chat");
            const messageDiv = document.createElement("div");
            messageDiv.className = "message";
            messageDiv.innerHTML = `<strong>${sender}:</strong> ${text}`;

            chat.appendChild(messageDiv);
            chat.scrollTop = chat.scrollHeight;
        } catch (err) {
            console.error("Ошибка при разборе JSON:", err);
        }
    };

    ws.onerror = function (error) {
        console.error("Ошибка WebSocket:", error);
    };

    ws.onclose = function () {
        console.log("Соединение закрыто");
    };
});

document.getElementById("input").addEventListener("keyup", function (event) {
    if (event.keyCode === 13 && ws && ws.readyState === WebSocket.OPEN) {
        const name = document.getElementById("nameInput").value.trim();
        const message = { sender: name, text: document.getElementById("input").value };
        ws.send(JSON.stringify(message));
        document.getElementById("input").value = "";
    }
});

function connectClientsWebSocket(room) {
    clientsWs = new WebSocket("ws://" + window.location.host + "/wscount?room=" + encodeURIComponent(room));

    clientsWs.onmessage = function (event) {
        try {
            const data = JSON.parse(event.data);
            const clientsList = document.getElementById("clientsList");
            clientsList.innerHTML = "";

            if (data.clients && data.clients.length > 0) {
                data.clients.forEach(clientName => {
                    const li = document.createElement("li");
                    li.textContent = clientName;
                    clientsList.appendChild(li);
                });
            } else {
                clientsList.innerHTML = "<li>Нет игроков</li>";
            }
        } catch (err) {
            console.error("Ошибка при разборе JSON от сервера:", err);
        }
    };

    clientsWs.onerror = function (error) {
        console.error("Ошибка WebSocket (clients):", error);
    };

    clientsWs.onclose = function () {
        console.log("Соединение со списком клиентов закрыто");
    };
}
