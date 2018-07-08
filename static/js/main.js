let username = ""

$.ajax({
    url: "https://randomuser.me/api/",
    async: false,
    success: function (result) {
        username = result.results[0].login.username
    }
});

console.log("Your username is: " + username)

let loc = window.location;
let uri = 'ws://' + loc.host + "/socket/" + username;
console.log("Connecting to " + uri)

ws = new WebSocket(uri);
ws.onopen = function () {
    console.log('Connected');
};

let i = -1;

let app = new Vue({
    el: "#app",
    data: {
        messages: [],
        message: "",
    },
    methods: {
        send: function (event) {
            let msg = {
                user: username,
                body: this.message,
            }

            this.$set(this.messages, ++i, msg)
            ws.send(JSON.stringify(msg));
            this.message = ''
        }
    }
});


app.$set(app.messages, ++i, {
    user: "DEBUG",
    body: "Your username is " + username,
})

ws.onmessage = function (evt) {
    console.log(evt);
    let data = JSON.parse(evt.data)
    app.$set(app.messages, ++i, data)
    document.getElementById( 'bottom' ).scrollIntoView()
};