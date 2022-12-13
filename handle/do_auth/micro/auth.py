from http.server import BaseHTTPRequestHandler, HTTPServer
import json

HOST = ("0.0.0.0", 8000)


class Engine(BaseHTTPRequestHandler):
    def do_GET(self):
        # do some mock auth
        if self.path == "/api1":
            self.send_response(200)
            self.send_header("Content-Type", "application/json; charset=utf-8")
            self.end_headers()
            self.wfile.write(json.dumps({"statut": 0, "path": self.path}).encode("utf8"))
        else:
            self.send_response(500)
            self.send_header("Content-Type", "application/json; charset=utf-8")
            self.end_headers()
            self.wfile.write(json.dumps({"statut": 1, "path": self.path}).encode("utf8"))


if __name__ == "__main__":
    server = HTTPServer(HOST, Engine)
    print(f"server on {HOST}")
    server.serve_forever()
