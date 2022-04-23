class Client extends EventTarget {
    constructor(uri) {
        super();
        this.socket = new WebSocket(uri);
        this.socket.onopen = this.handleOpen;
        this.socket.onclose = this.handleClose;
    }

    handleOpen() {
        console.log('connetion opened');
    }

    handleClose(event) {
        if (event.wasClean) {
            console.log(`connection closed, code=${event.code} reason=${event.reason}`);
        } else {
            console.log(event);
            console.error('connection died');
        }
    }

    handleJoined(event) {
        this.id = event.data.client_id;
        this.dispatchEvent(new Event('joined'));
    }

    handleEvent(event) {
        const e = new CustomEvent(event.type)
    }

    sendEvent(event) {

    }
}