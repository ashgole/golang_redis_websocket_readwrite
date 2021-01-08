new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent2: '', // A running list of chat messages displayed on the screen
        rollNo: null, // Our username
        totalWord: null, // Our charCount
        totalChar: null, // Our charCount
        wpm: null, // Our charCount
        joined: false // True if email and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://127.0.0.1:8000/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data); 
				
		    self.chatContent2 ="<table border='1'>"
			+"<th>Roll No</th><th>message</th><th>Total Words</th><th>Total Characters</th><th>Words/Minute</th>"
			+"<tr>"
                    +"<td>"+msg.rollNo+"</td>"
                    +"<td>"+msg.message+"</td>"
                    +"<td>"+msg.totalWord+"</td>"
                    +"<td>"+msg.totalChar+"</td>"
                    +"<td>"+msg.wpm+"</td>"
                + "</tr></table>" 
				; // Parse emojis; // Parse emojis
				
			var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },

 
 
});