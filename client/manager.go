package client

import (
	"fmt"
	"sync"

	"PIM_Server/log"
)

// Manager 连接管理
type Manager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // platform+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 用户连接处理
	Login       chan *Login        // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

// NewManager New客户端管理器
func NewManager() (manager *Manager) {
	manager = &Manager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *Login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
	return
}

// GetUserKey 获取用户key
func GetUserKey(platform PlatformType, uid int64) (key string) {
	return fmt.Sprintf("%d_%d", platform, uid)
}

/**************************  manager  ***************************************/

// InClient 判断是否存在
func (manager *Manager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]

	return
}

// GetClients 获取连接
func (manager *Manager) GetClients() (clients map[*Client]bool) {

	clients = make(map[*Client]bool)

	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value

		return true
	})

	return
}

// ClientsRange 遍历
func (manager *Manager) ClientsRange(f func(client *Client, value bool) (result bool)) {

	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}

	return
}

// GetClientsLen 获取连接数量
func (manager *Manager) GetClientsLen() (clientsLen int) {

	clientsLen = len(manager.Clients)

	return
}

// AddClients 添加客户端
func (manager *Manager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *Manager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// GetUserClient 获取用户的连接
func (manager *Manager) GetUserClient(platformID PlatformType, uid int64) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	userKey := GetUserKey(platformID, uid)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}

	return
}

// GetUsersLen 获取用户数量
func (manager *Manager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)

	return
}

// AddUsers 添加用户
func (manager *Manager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	manager.Users[key] = client
}

// DelUsers 删除用户
func (manager *Manager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	key := GetUserKey(client.PlatformId, client.Uid)
	if value, ok := manager.Users[key]; ok {
		// 判断是否为相同的用户
		if value.Addr != client.Addr {

			return
		}
		delete(manager.Users, key)
		result = true
	}

	return
}

// GetUserKeys 获取所有用户的key
func (manager *Manager) GetUserKeys() (userKeys []string) {

	userKeys = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}

	return
}

// GetUserList 获取指定平台下的用户uid列表
func (manager *Manager) GetUserList(platformId PlatformType) (userList []int64) {

	userList = make([]int64, 0)

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	for _, v := range manager.Users {
		if v.PlatformId == platformId {
			userList = append(userList, v.Uid)
		}
	}

	fmt.Println("GetUserList len:", len(manager.Users))

	return
}

// GetUserClients 获取用户的key
func (manager *Manager) GetUserClients() (clients []*Client) {

	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}

	return
}

// sendAll 向全部成员(除了自己)发送数据
func (manager *Manager) sendAll(message []byte, ignoreClient *Client) {

	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient {
			conn.SendRawMsg(message)
		}
	}
}

// sendAppIdAll 向全部成员(除了自己)发送数据
func (manager *Manager) sendAppIdAll(message []byte, platformId PlatformType, ignoreClient *Client) {

	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient && conn.PlatformId == platformId {
			conn.SendRawMsg(message)
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *Manager) EventRegister(client *Client) {
	manager.AddClients(client)
	manager.AddUsers(GetUserKey(client.PlatformId, client.Uid), client)

	fmt.Println("EventRegister client conn, addr", client.Addr)
	go client.Run()
	client.Send <- []byte("{\"event\":\"connect\",\"ping_interval\":\"20\",\"ping_timeout\":60}")
}

// EventLogin 用户登录
// todo: 先仅处理本地
func (manager *Manager) EventLogin(login *Login) {

	if login == nil {
		log.Infof("client is nil.")
		return
	}

	cli := login.Client
	// 连接存在，在添加
	if manager.InClient(cli) {
		userKey := login.GetKey()
		manager.AddUsers(userKey, login.Client)
	}

	fmt.Println("EventLogin client auth", cli.Addr, login.PlatformID, login.Uid)

	//go cli.Run()

	//orderId := helper.GetOrderIdTime()
	//SendUserMessageAll(auth.PlatformID, auth.UserId, orderId, models.MessageCmdEnter, "哈喽~")
}

// EventUnregister 用户断开连接
// todo: 先仅处理本地
func (manager *Manager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	deleteResult := manager.DelUsers(client)
	if deleteResult == false {
		// 不是当前连接的客户端

		return
	}

	// 清除redis登录数据
	//userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	//if err == nil {
	//	userOnline.LogOut()
	//	cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	//}

	fmt.Println("EventUnregister client logout", client.Addr, client.PlatformId, client.Uid)
	// 关闭 客户端
	client.Close()

	//if client.UserId != "" {
	//	orderId := helper.GetOrderIdTime()
	//	SendUserMessageAll(client.PlatformID, client.UserId, orderId, models.MessageCmdExit, "用户已经离开~")
	//}
}

// Start 管道处理程序
func (manager *Manager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)

		case login := <-manager.Login:
			// 用户登录
			manager.EventLogin(login)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)

		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}
