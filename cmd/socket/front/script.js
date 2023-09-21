const socket = new WebSocket('ws://localhost:8080'); // WebSocket连接地址

// 获取HTML元素
const messageInput = document.getElementById('message-input');
const sendButton = document.getElementById('send-button');
const chatMessages = document.getElementById('chat-messages');

// 处理WebSocket连接建立
socket.addEventListener('open', (event) => {
    console.log('WebSocket connected');
});

// 处理WebSocket消息接收
socket.addEventListener('message', (event) => {
    const message = event.data;
    displayMessage(message);
});

// 处理WebSocket连接关闭
socket.addEventListener('close', (event) => {
    console.log('WebSocket closed');
});

// 发送消息
sendButton.addEventListener('click', () => {
    const message = messageInput.value;
    if (message) {
        socket.send(message);
        messageInput.value = '';
    }
});

// 显示消息在聊天界面
function displayMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.textContent = message;
    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}
