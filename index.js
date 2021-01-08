new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        rollNo: null, // Our username
        totalWord: null, // Our charCount
        totalChar: null, // Our charCount
        wpm: null, // Our charCount
        joined: false // True if email and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://127.0.0.1:8000/ws');
        
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        rollNo: document.getElementById("rollno").innerHTML,
                        message: $('<p>').html(this.newMsg).text(),
						totalWord: document.getElementById("totalWordId").innerHTML, 					
                      	totalChar: document.getElementById("totalCharId").innerHTML ,					
                      	wpm: document.getElementById("wpmId").innerHTML 					
                     }
                ));
               // this.newMsg = ''; // Reset newMsg
            }
        },

     

         
    }
});