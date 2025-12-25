let socket = null;
let reconnectTimer = null;

function connectChat() {
  const token = localStorage.getItem("token");

  if (!token) {
    console.warn("⚠️ No token, cannot connect chat");
    return;
  }

  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("ℹ️ WebSocket already connected");
    return;
  }

  console.log("🔌 Connecting to chat server...");

  socket = new WebSocket(
    `ws://localhost:8080/ws?token=${encodeURIComponent(token)}`
  );

  socket.onopen = () => {
    console.log("✅ Connected to MangaHub Chat");
    document.getElementById("chat-status").innerText = "Online";

    // Clear reconnect timer nếu có
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  };

  socket.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data);
      displayMessage(msg);
    } catch (err) {
      console.error("❌ Invalid message format:", event.data);
    }
  };

  socket.onclose = (event) => {
    console.warn(
      `❌ Disconnected (code=${event.code}, reason=${
        event.reason || "no reason"
      })`
    );
    document.getElementById("chat-status").innerText = "Offline";

    if (!reconnectTimer) {
      reconnectTimer = setTimeout(() => {
        reconnectTimer = null;
        connectChat();
      }, 3000);
    }
  };

  socket.onerror = (err) => {
    console.error("❌ WebSocket error", err);
    socket.close();
  };
}

function sendMessage() {
  const input = document.getElementById("chat-input");
  const content = input.value.trim();

  if (!socket || socket.readyState !== WebSocket.OPEN) {
    alert("Chat is not connected");
    return;
  }

  if (content) {
    socket.send(content);
    input.value = "";
  }
}

function displayMessage(msg) {
  const container = document.getElementById("chat-messages");

  const msgHtml = `
    <div class="chat-msg">
      <span class="user" style="color:#ff0055;font-weight:bold;">
        ${msg.username}:
      </span>
      <span class="content">${msg.content}</span>
      <small style="color:#666;font-size:0.7em;">
        ${msg.time}
      </small>
    </div>
  `;

  container.insertAdjacentHTML("beforeend", msgHtml);
  container.scrollTop = container.scrollHeight;
}
