var WebSocketClient = require('websocket').client;
var moment = require("moment")


// 同時連線到 ws server 的人數
const clientCount = 10;

for (let i = 0; i < clientCount; i++) {

    let client = new WebSocketClient();

    client.on('connectFailed', function (error) {
        console.log('Connect Error: ' + error.toString());
    });

    client.on('connect', function (connection) {
        console.log('WebSocket client connected');
        connection.on('error', function (error) {
            console.log("Connection Error: " + error.toString());
        });
        connection.on('close', function () {
            console.log('Connection Closed');
        });
        connection.on('message', function (message) {
            if (message.type === 'utf8') {
                console.log("Received: '" + message.utf8Data + "'");
            }
        });

        function sendNumber() {
            if (connection.connected) {
                let date = moment().format('MMMM Do YYYY, h:mm:ss a');
                connection.sendUTF(date.toString());

                let min = 1000
                let max = 4000
                let timeout = Math.random() * (max - min) + min;

                setTimeout(sendNumber, parseInt(timeout));
            }
        }

        sendNumber();
    });

// command: node client.js
    client.connect('ws://localhost:8011/ws');

    console.log("connected")
}



