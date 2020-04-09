let ws;

$("#connect_ws").click(function() {
    let ws_host = $('#ws_host').val();
    let ws_user_id = $('#ws_user_id').val();
    let ws_device = $('#ws_device').val();

    ws = new WebSocket(ws_host + "&user_id=" + ws_user_id + "&device=" + ws_device)

    ws.onopen = function() {};

    ws.onmessage = function (evt) {
        showMessage("chat", ws_user_id, evt.data)
    };

    ws.onclose = function() { 
    alert("Connection is closed..."); 
    };

    $("#connect_ws").hide();
    $("#disconnect_ws").show();
    $("#send_message").show();
});

$("#disconnect_ws").click(function() {
    $("#disconnect_ws").hide();
    $("#send_message").hide();
    $("#connect_ws").show();

    ws.close();
});

$("#form_message").submit(function(e) {
    e.preventDefault();

    let msg = $('#ws_message').val();
    ws.send(msg);
});

function showMessage(type, user, msg) {
    if (type == "chat") {
        $("#chat-list").append(`<li>${user}: ${msg}</li>`);
    } else {
        $("#chat-list").append(`<li>${user}: <b>${type}</b></li>`);
        console.log(type, user, msg)
    }
}