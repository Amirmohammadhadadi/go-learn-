// ============================================================================
// FILE: websocket_guide.go
// TITLE: راهنمای کامل WebSocket با gorilla/websocket
// HOW TO RUN: go run websocket_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - WebSocket چیست؟
// ============================================================================
//
// WebSocket یک پروتکل ارتباطی دوطرفه (full-duplex) است که یک کانال ارتباطی
// persistent بین کلاینت و سرور ایجاد می‌کند.
//
// تفاوت با HTTP:
// 1. HTTP: درخواست-پاسخ (یک‌طرفه، هر بار اتصال جدید)
// 2. WebSocket: دوطرفه، persistent، کم‌تاخیر
//
// کاربردهای WebSocket:
// 1. چت و پیام‌رسانی实时
// 2. بازی‌های آنلاین
// 3. نمایش داده‌های زنده (قیمت سهام، نمرات ورزشی)
// 4. همکاری实时 (مثل Google Docs)
// 5. نوتیفیکیشن‌های server push
//
// قانون طلایی:
// "از WebSocket برای ارتباطات实时 دوطرفه استفاده کن.
//  همیشه Ping/Pong را پیاده‌سازی کن تا اتصالات مرده را تشخیص دهی.
//  از Mutex برای محافظت از map کلاینت‌ها استفاده کن."
// ============================================================================

package __websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ============================================================================
// بخش 1: پیکربندی WebSocket
// ============================================================================

// Upgrader برای ارتقاء اتصال HTTP به WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // سایز بافر خواندن
	WriteBufferSize: 1024, // سایز بافر نوشتن
	CheckOrigin: func(r *http.Request) bool {
		// در production، originها را بررسی کن
		return true
	},
	// Subprotocols: []string{"chat"}, // پروتکل‌های پشتیبانی شده
}

// ============================================================================
// بخش 2: انواع پیام‌ها
// ============================================================================

// MessageType انواع پیام‌های WebSocket
type MessageType string

const (
	TypeChat     MessageType = "chat"
	TypeJoin     MessageType = "join"
	TypeLeave    MessageType = "leave"
	TypeTyping   MessageType = "typing"
	TypePing     MessageType = "ping"
	TypePong     MessageType = "pong"
	TypeError    MessageType = "error"
	TypeUserList MessageType = "user_list"
)

// Message ساختار پیام WebSocket
type Message struct {
	Type      MessageType `json:"type"`
	Sender    string      `json:"sender,omitempty"`
	Content   string      `json:"content,omitempty"`
	Room      string      `json:"room,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Users     []string    `json:"users,omitempty"`
}

// ============================================================================
// بخش 3: کلاینت WebSocket
// ============================================================================

// Client نماینده یک اتصال WebSocket
type Client struct {
	ID   string
	Room string
	Name string
	Conn *websocket.Conn
	Send chan Message
	Hub  *Hub
	mu   sync.Mutex
}

// NewClient ایجاد کلاینت جدید
func NewClient(conn *websocket.Conn, hub *Hub, room, name string) *Client {
	return &Client{
		ID:   generateID(),
		Room: room,
		Name: name,
		Conn: conn,
		Send: make(chan Message, 256),
		Hub:  hub,
	}
}

// readPump خواندن پیام‌ها از WebSocket
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// تنظیم Pong handler برای تشخیص زنده بودن اتصال
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error: %v", err)
			}
			break
		}

		msg.Sender = c.Name
		msg.Timestamp = time.Now()
		c.Hub.broadcast <- msg
	}
}

// writePump نوشتن پیام‌ها به WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(msg); err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		case <-ticker.C:
			// ارسال Ping برای حفظ اتصال
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ============================================================================
// بخش 4: Hub (مرکز مدیریت اتصالات)
// ============================================================================

// Hub مدیریت همه کلاینت‌ها و اتاق‌ها
type Hub struct {
	// Rooms: map[roomName]map[clientID]*Client
	rooms      map[string]map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mu         sync.RWMutex
}

// NewHub ایجاد Hub جدید
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message, 256),
	}
}

// Run اجرای Hub (گوروتین اصلی)
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// اضافه کردن به اتاق
			if _, ok := h.rooms[client.Room]; !ok {
				h.rooms[client.Room] = make(map[string]*Client)
			}
			h.rooms[client.Room][client.ID] = client
			h.mu.Unlock()

			// ارسال لیست کاربران به اتاق
			h.sendUserList(client.Room)

			// اعلام ورود کاربر
			h.broadcast <- Message{
				Type:      TypeJoin,
				Sender:    client.Name,
				Content:   fmt.Sprintf("%s joined the room", client.Name),
				Room:      client.Room,
				Timestamp: time.Now(),
			}
			log.Printf("Client %s joined room %s", client.Name, client.Room)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.Room]; ok {
				if _, ok := clients[client.ID]; ok {
					delete(clients, client.ID)
					if len(clients) == 0 {
						delete(h.rooms, client.Room)
					}
				}
			}
			h.mu.Unlock()
			close(client.Send)

			// ارسال لیست کاربران به اتاق
			h.sendUserList(client.Room)

			// اعلام خروج کاربر
			h.broadcast <- Message{
				Type:      TypeLeave,
				Sender:    client.Name,
				Content:   fmt.Sprintf("%s left the room", client.Name),
				Room:      client.Room,
				Timestamp: time.Now(),
			}
			log.Printf("Client %s left room %s", client.Name, client.Room)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients, ok := h.rooms[msg.Room]
			if !ok {
				h.mu.RUnlock()
				continue
			}
			// ارسال به همه کلاینت‌های اتاق
			for _, client := range clients {
				select {
				case client.Send <- msg:
				default:
					// اگر بافر پر است، کلاینت را disconnect کن
					close(client.Send)
					delete(clients, client.ID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// sendUserList ارسال لیست کاربران اتاق
func (h *Hub) sendUserList(room string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[room]
	if !ok {
		return
	}

	users := make([]string, 0, len(clients))
	for _, client := range clients {
		users = append(users, client.Name)
	}

	msg := Message{
		Type:      TypeUserList,
		Users:     users,
		Room:      room,
		Timestamp: time.Now(),
	}

	for _, client := range clients {
		select {
		case client.Send <- msg:
		default:
		}
	}
}

// ============================================================================
// بخش 5: HTTP Handlers
// ============================================================================

var hub = NewHub()

// WebSocketHandler مدیریت اتصالات WebSocket
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// استخراج پارامترها
	room := r.URL.Query().Get("room")
	name := r.URL.Query().Get("name")

	if room == "" {
		room = "default"
	}
	if name == "" {
		name = fmt.Sprintf("Anonymous-%d", time.Now().UnixNano())
	}

	// ارتقاء اتصال
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	// ایجاد کلاینت
	client := NewClient(conn, hub, room, name)

	// ثبت در Hub
	hub.register <- client

	// شروع گوروتین‌های خواندن و نوشتن
	go client.writePump()
	go client.readPump()
}

// ============================================================================
// بخش 6: Health Check و Status
// ============================================================================

func healthHandler(w http.ResponseWriter, r *http.Request) {
	hub.mu.RLock()
	roomCount := len(hub.rooms)
	clientCount := 0
	for _, clients := range hub.rooms {
		clientCount += len(clients)
	}
	hub.mu.RUnlock()

	response := map[string]interface{}{
		"status":    "ok",
		"rooms":     roomCount,
		"clients":   clientCount,
		"timestamp": time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// بخش 7: HTML Client (برای تست)
// ============================================================================

func serveHTML(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Chat</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background: #f0f0f0;
        }
        .chat-container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .chat-header {
            background: #4CAF50;
            color: white;
            padding: 15px;
            text-align: center;
        }
        .chat-messages {
            height: 400px;
            overflow-y: auto;
            padding: 15px;
            background: #f9f9f9;
        }
        .message {
            margin-bottom: 10px;
            padding: 8px 12px;
            border-radius: 8px;
            word-wrap: break-word;
        }
        .message.system {
            background: #e0e0e0;
            color: #666;
            text-align: center;
            font-style: italic;
        }
        .message.join {
            background: #d4edda;
            color: #155724;
        }
        .message.leave {
            background: #f8d7da;
            color: #721c24;
        }
        .message.chat {
            background: #007bff;
            color: white;
            margin-left: 20px;
        }
        .message.chat .sender {
            font-weight: bold;
            margin-bottom: 5px;
        }
        .chat-input {
            display: flex;
            padding: 15px;
            background: white;
            border-top: 1px solid #ddd;
        }
        .chat-input input {
            flex: 1;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-right: 10px;
        }
        .chat-input button {
            padding: 10px 20px;
            background: #4CAF50;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
        .chat-input button:hover {
            background: #45a049;
        }
        .user-list {
            float: right;
            width: 200px;
            background: #f0f0f0;
            padding: 10px;
            border-left: 1px solid #ddd;
            min-height: 400px;
        }
        .user-list h4 {
            margin: 0 0 10px 0;
        }
        .user-list ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
        .user-list li {
            padding: 5px;
            border-bottom: 1px solid #ddd;
        }
        .status {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
        }
    </style>
</head>
<body>
    <div class="chat-container">
        <div class="chat-header">
            <h2>WebSocket Chat Room</h2>
            <div class="status" id="status">Connecting...</div>
        </div>
        <div style="display: flex;">
            <div style="flex: 1;">
                <div class="chat-messages" id="messages"></div>
                <div class="chat-input">
                    <input type="text" id="messageInput" placeholder="Type a message..." />
                    <button onclick="sendMessage()">Send</button>
                </div>
            </div>
            <div class="user-list">
                <h4>Users Online</h4>
                <ul id="userList"></ul>
            </div>
        </div>
    </div>

    <script>
        let ws;
        let userName = prompt("Enter your name:", "User-" + Math.floor(Math.random() * 1000));
        let roomName = prompt("Enter room name (or leave empty for default):", "default");

        function connect() {
            const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
            const wsUrl = protocol + "//" + window.location.host + "/ws?room=" + (roomName || "default") + "&name=" + userName;
            ws = new WebSocket(wsUrl);

            ws.onopen = function() {
                document.getElementById("status").innerHTML = "Connected ✅";
                document.getElementById("status").style.color = "green";
                document.getElementById("messageInput").disabled = false;
            };

            ws.onmessage = function(event) {
                const msg = JSON.parse(event.data);
                displayMessage(msg);
            };

            ws.onclose = function() {
                document.getElementById("status").innerHTML = "Disconnected ❌";
                document.getElementById("status").style.color = "red";
                document.getElementById("messageInput").disabled = true;
                // Attempt to reconnect after 3 seconds
                setTimeout(connect, 3000);
            };

            ws.onerror = function(error) {
                console.error("WebSocket error:", error);
            };
        }

        function displayMessage(msg) {
            const messagesDiv = document.getElementById("messages");
            const messageDiv = document.createElement("div");
            messageDiv.className = "message " + msg.type;

            let content = "";
            if (msg.type === "chat") {
                content = '<div class="sender">' + escapeHtml(msg.sender) + ':</div>' +
                         '<div>' + escapeHtml(msg.content) + '</div>' +
                         '<div style="font-size: 10px; opacity: 0.7;">' + new Date(msg.timestamp).toLocaleTimeString() + '</div>';
            } else if (msg.type === "join") {
                content = '🔵 ' + escapeHtml(msg.content);
            } else if (msg.type === "leave") {
                content = '🔴 ' + escapeHtml(msg.content);
            } else if (msg.type === "user_list") {
                updateUserList(msg.users);
                return;
            } else {
                content = escapeHtml(msg.content);
            }

            messageDiv.innerHTML = content;
            messagesDiv.appendChild(messageDiv);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }

        function updateUserList(users) {
            const userList = document.getElementById("userList");
            userList.innerHTML = "";
            users.forEach(user => {
                const li = document.createElement("li");
                li.textContent = user + (user === userName ? " (you)" : "");
                userList.appendChild(li);
            });
        }

        function sendMessage() {
            const input = document.getElementById("messageInput");
            const content = input.value.trim();
            if (content && ws && ws.readyState === WebSocket.OPEN) {
                ws.send(JSON.stringify({
                    type: "chat",
                    content: content
                }));
                input.value = "";
            }
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        document.getElementById("messageInput").addEventListener("keypress", function(e) {
            if (e.key === "Enter") {
                sendMessage();
            }
        });

        connect();
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ============================================================================
// بخش 8: Broadcast API (ارسال پیام به همه کلاینت‌ها از طریق HTTP)
// ============================================================================

type BroadcastRequest struct {
	Room    string `json:"room"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

func broadcastHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BroadcastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	msg := Message{
		Type:      TypeChat,
		Sender:    req.Sender,
		Content:   req.Message,
		Room:      req.Room,
		Timestamp: time.Now(),
	}

	hub.broadcast <- msg

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "sent",
		"message": msg,
	})
}

// ============================================================================
// بخش 9: Chat Room Manager (مدیریت اتاق‌ها)
// ============================================================================

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

type Room struct {
	Name      string
	Clients   map[string]*Client
	CreatedAt time.Time
}

var roomManager = &RoomManager{
	rooms: make(map[string]*Room),
}

func (rm *RoomManager) GetRooms() []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rooms := make([]string, 0, len(rm.rooms))
	for name := range rm.rooms {
		rooms = append(rooms, name)
	}
	return rooms
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms := roomManager.GetRooms()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rooms": rooms,
		"count": len(rooms),
	})
}

// ============================================================================
// بخش 10: Stats و Monitoring
// ============================================================================

func statsHandler(w http.ResponseWriter, r *http.Request) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["total_rooms"] = len(hub.rooms)
	stats["total_clients"] = 0

	roomStats := make(map[string]int)
	for roomName, clients := range hub.rooms {
		count := len(clients)
		roomStats[roomName] = count
		stats["total_clients"] = stats["total_clients"].(int) + count
	}
	stats["rooms"] = roomStats

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ============================================================================
// بخش 11: Graceful Shutdown
// ============================================================================

func gracefulShutdown() {
	// در یک برنامه واقعی، از signal برای shutdown استفاده کن
	log.Println("Shutting down gracefully...")

	// بستن همه اتصالات
	hub.mu.RLock()
	for _, clients := range hub.rooms {
		for _, client := range clients {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte("Server shutting down"))
			client.Conn.Close()
		}
	}
	hub.mu.RUnlock()

	log.Println("Shutdown complete")
}

// ============================================================================
// بخش 12: Helper Functions
// ============================================================================

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 WEBSOCKET GUIDE WITH GORILLA/WEBSOCKET")
	fmt.Println("Real-time bidirectional communication")
	fmt.Println(stringsRepeat("=", 80))

	// راه‌اندازی Hub در گوروتین جدا
	go hub.Run()

	// تنظیم routes
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/ws", WebSocketHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/rooms", roomsHandler)
	http.HandleFunc("/broadcast", broadcastHandler)

	// سرور
	port := ":8080"
	log.Printf("WebSocket server starting on http://localhost%s", port)
	log.Printf("Open http://localhost%s in multiple browsers to test chat", port)
	log.Printf("Endpoints:")
	log.Printf("  GET  /           - Chat client")
	log.Printf("  GET  /ws         - WebSocket endpoint")
	log.Printf("  GET  /health     - Health check")
	log.Printf("  GET  /stats      - Server statistics")
	log.Printf("  GET  /rooms      - List rooms")
	log.Printf("  POST /broadcast  - Broadcast message via HTTP")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

// ============================================================================
// بخش 14: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 WEBSOCKET BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ 1. Ping/Pong Mechanism                                        │
│    - Send periodic ping messages                               │
│    - Handle pong to detect dead connections                    │
│    - Close connections that don't respond                      │
│                                                                 │
│ 2. Reconnection Logic                                          │
│    - Implement exponential backoff                             │
│    - Handle reconnection on client side                        │
│    - Preserve state if possible                                │
│                                                                 │
│ 3. Message Buffering                                           │
│    - Use buffered channels for messages                        │
│    - Set appropriate buffer sizes                              │
│    - Drop messages if buffer is full                           │
│                                                                 │
│ 4. Authentication                                              │
│    - Authenticate before upgrading to WebSocket                │
│    - Use tokens in query string or headers                     │
│    - Re-validate periodically                                  │
│                                                                 │
│ 5. Rate Limiting                                               │
│    - Limit messages per second                                 │
│    - Protect against DoS attacks                               │
│    - Implement per-client limits                               │
│                                                                 │
│ 6. Message Format                                              │
│    - Use JSON for structured messages                          │
│    - Include message type and version                          │
│    - Add timestamps and IDs                                    │
│                                                                 │
│ 7. Error Handling                                              │
│    - Handle unexpected close errors                            │
│    - Log errors but don't expose to clients                    │
│    - Implement graceful degradation                            │
│                                                                 │
│ 8. Scaling                                                     │
│    - Use message brokers (Redis, Kafka) for multi-server       │
│    - Implement sticky sessions if needed                       │
│    - Consider WebSocket gateways                               │
│                                                                 │
│ 9. Security                                                    │
│    - Validate Origin headers                                   │
│    - Use WSS (WebSocket Secure) in production                  │
│    - Implement message size limits                             │
│    - Sanitize user input                                       │
│                                                                 │
│ 10. Monitoring                                                 │
│     - Track active connections                                 │
│     - Monitor message rates                                    │
│     - Log errors and reconnections                             │
└─────────────────────────────────────────────────────────────────┘
`)
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

/*
# اجرای سرور
go run websocket_guide.go

# باز کردن مرورگر در آدرس‌های مختلف
http://localhost:8080

# تست WebSocket با wscat (ابزار خط فرمان)
npm install -g wscat
wscat -c "ws://localhost:8080/ws?room=test&name=Ali"
*/
/*
خلاصه WebSocket مفاهیم کلیدی
مفهوم	توضیح
Upgrader	ارتقاء اتصال HTTP به WebSocket
Client	نماینده یک اتصال WebSocket
Hub	مدیریت همه اتصالات و اتاق‌ها
Message	ساختار استاندارد پیام
Ping/Pong	مکانیزم بررسی زنده بودن اتصال
Room	گروه‌بندی کلاینت‌ها
انواع پیام‌ها
Type	توضیح
chat	پیام چت معمولی
join	کاربر وارد اتاق شد
leave	کاربر اتاق را ترک کرد
typing	کاربر در حال تایپ است
ping/pong	بررسی اتصال
user_list	لیست کاربران آنلاین

*/
