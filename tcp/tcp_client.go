package tcp

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"extruder_web_gui/config"
	"extruder_web_gui/data"
	"log"
	"net"
	"strings"
	"sync"
)

type ConnectionManager struct {
	mu      sync.Mutex
	conn    net.Conn
	address string
}

var manager *ConnectionManager
var once sync.Once

// Return singleton instance of connection manager
func GetConnectionManager(address string) *ConnectionManager {
	once.Do(func() {
		manager = &ConnectionManager{
			address: address,
		}
	})
	return manager
}

// Establish TCP connection
func (cm *ConnectionManager) GetConnection() net.Conn {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conn == nil {
		var err error
		cm.conn, err = net.Dial("tcp", cm.address)
		if err != nil {
			log.Fatalf("Failed to establish TCP connection: %v", err)
		}
		log.Println("TCP connection established:", cm.address)
	}
	return cm.conn
}

// Function to close tcp connection
func (cm *ConnectionManager) CloseConnection() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.conn != nil {
		cm.conn.Close()
		cm.conn = nil
		log.Println("TCP connection closed")
	}
}

// Handle incoming data
func (cm *ConnectionManager) ReadLoop(processFunc func(string)) {
	conn := cm.GetConnection()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from TCP connection:", err)
			break
		}
		processFunc(strings.TrimSpace(line))
	}
}

func TCPDataHandler(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Error connectiong to %s: %v", address, err)
		return
	}
	defer conn.Close()
	log.Println("Connected to TCP-Server:", address)

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading data:", err)
			break
		}
		ProcessTCPData(strings.TrimSpace(line))
	}
}

func ProcessTCPData(line string) {
	msg, err := hex.DecodeString((line))
	//fmt.Println(line)
	if err != nil {
		log.Fatalf("Hex string decoding failed: %v", err)

	}

	if len(msg) != 8 {
		log.Printf("Expected 8 bytes, got %d bytes", len(msg))
		return // Skip this message
	}
	data.GetValueFromMsg(msg)
}

func SendTCPData(id byte, val uint32) {
	manager := GetConnectionManager(config.Cfg.TCPAddress)
	conn := manager.GetConnection()

	messageBytes := make([]byte, 8)
	messageBytes[7] = id
	binary.LittleEndian.PutUint32(messageBytes[0:], val)

	_, err := conn.Write(messageBytes)
	if err != nil {
		log.Fatalf("Failed to write to TCP connection: %v", err)
	}
	log.Printf("Message sent via TCP: %x\n", reverseBytes(messageBytes))
}

// Helper function to reverse byte order
func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
