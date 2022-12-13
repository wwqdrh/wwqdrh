from http.server import BaseHTTPRequestHandler, HTTPServer
import json

HOST = ("0.0.0.0", 8000)


class Engine(BaseHTTPRequestHandler):
    def do_GET(self):
        # do some mock auth
        if self.path == "/api1":
            self.api1()
        elif self.path == "/api2":
            self.api2()
        else:
            self.send_response_only(404)

    def api1(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.end_headers()
        self.wfile.write(json.dumps({"path": "/api1", "user": self.headers["user"]}).encode("utf8"))
    
    def api2(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.end_headers()
        self.wfile.write(json.dumps({"path": "/api2"}).encode("utf8"))

if __name__ == "__main__":
    server = HTTPServer(HOST, Engine)
    print(f"server on {HOST}")
    server.serve_forever()
