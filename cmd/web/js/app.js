import ePub from 'epubjs';

const message = type => id => data => {
    return {
        type: type,
        sender_id: id,
        data: data,
    }
}

class Orchestrator {
    constructor(socket, rendition) {
        this.socket = socket;
        this.rendition = rendition;
        this.handleOpen = this.handleOpen.bind(this);
        this.handleMessage = this.handleMessage.bind(this);
        this.handleClose = this.handleClose.bind(this);
        this.handleError = this.handleError.bind(this);
        this.handleRelocated = this.handleRelocated.bind(this);
        this.handleHightlighted = this.handleHightlighted.bind(this);
        this.handleSendSelected = this.handleSendSelected.bind(this);
        this.handleSendRelocated = this.handleSendRelocated.bind(this);
    }

    init() {
        this.socket.onopen = this.handleOpen;
        this.socket.onmessage = this.handleMessage;
        this.socket.onclose = this.handleClose;
        this.socket.onerror = this.handleError;
        this.rendition.on('selected', this.handleSendSelected);
        this.rendition.on('relocated', this.handleSendRelocated);
    }

    handleOpen() {
        console.log('connection opened');
    }

    handleClose(event) {
        if (event.wasClean) {
            console.log(`connection closed, code=${event.code} reason=${event.reason}`);
        } else {
            console.log(event);
            console.error('connection died');
        }
    }

    handleError(error) {
        console.error(error.message);
    }

    handleMessage(event) {
        const m = JSON.parse(event.data)
        switch(m.type) {
            case 'relocated':
                this.handleRelocated(m.data.cfi);
                break;
            case 'highlighted':
                this.handleHightlighted(m.data.cfiRange);
                break;
            default:
                console.log(`unrecognized message type of ${message.type}`);
        }
    }

    handleRelocated(cfi) {
        if (typeof cfi !== 'string') {
            this.handleError(new Error(`handleRelocated expected a string got ${typeof cfi}`));
        }
        const currentLocation = this.rendition.currentLocation();
        console.log(cfi);
        if (currentLocation.start.cfi != cfi) {
            this.rendition.display(cfi);
        }
    }
    
    handleHightlighted(cfiRange) {
        this.rendition.annotations.highlight(
            cfiRange,
            {},
            e => console.log(e)
        );
    }

    handleSendSelected(cfiRange) {
        const m = message('highlighted')({cfiRange});
        this.socket.send(JSON.stringify(m));
    }

    handleSendRelocated(location) {
        const cfi = location.start.cfi;
        const m = message('relocated')({cfi});
        this.socket.send(JSON.stringify(m));
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const room = window.location.pathname.split("/").pop();
    const socket = new WebSocket(`ws://localhost:8080/ws/${room}`);
    const bookUri = '/public/childrens-literature.epub';
    const book = ePub(bookUri);
    const rendition = book.renderTo('viewer', {
        width: '100%',
        height: 400,
    })
    rendition.display();
    const o = new Orchestrator(socket, rendition);
    o.init()

    const prev = document.getElementById('prev');
    const next = document.getElementById('next');

    prev.addEventListener('click', function() {
        rendition.prev();
    });

    next.addEventListener('click', function() {
        rendition.next();
    })
})
