(function () {
    var page = {
        controller: function () {
          page.socketConnect();
        },
        socketConnect: function () {
            var ws = new WebSocket("ws://127.0.0.1:9999");
            ws.onopen = function (ev) {
                console.log("open", ev);
            };
            ws.onmessage = function (ev) {
                console.log("message", ev);
            };
            ws.onclose = function (ev) {
                console.log("close",ev);
            };
            ws.onerror = function (ev) {
                console.log("error", ev);
            }
        },

    };
    page.controller();
})()