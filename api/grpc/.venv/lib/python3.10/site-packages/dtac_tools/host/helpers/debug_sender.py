import socket

class DebugSender:
    def __init__(self, port):
        self.port = port
        self.address = ('127.0.0.1', port)
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    def write(self, message):
        if not isinstance(message, bytes):
            message = message.encode('utf-8')
        self.sock.sendto(message, self.address)

    def close(self):
        self.sock.close()

if __name__ == "__main__":
    udp_sender = DebugSender(9999)  # replace 9999 with your desired port
    udp_sender.write("Hello, UDP!")
    udp_sender.close()
