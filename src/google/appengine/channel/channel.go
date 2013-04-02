package channel


               socket = channel.open();
                socket.onopen = onOpened;
                socket.onmessage = onMessage;
                socket.onerror = onError;
                socket.onclose = onClose;

type Channel interface {
    Open() (socket Socket, err error)
}

type Socket interface {
    OnOpened()
    OnMessage(msg string)
    onError(err error)
    onClose()
}
