//此文件静态导入，不需要编译
'use strict'


// Websocket 连接和消息处理
// https://www.ruanyifeng.com/blog/2017/05/websocket.html
// https://developer.mozilla.org/zh-CN/docs/Web/API/WebSocket

// 定义WebSocket变量和重连参数
let socket
let reconnectAttempts = 0
const maxReconnectAttempts = 200
const reconnectInterval = 3000 // 每次重连间隔3秒

// 用户ID和令牌，假设已在其他地方定义
const userID = Alpine.store('global').userID
// 假设token是一个有效的令牌 TODO:使用真正的令牌
const token = 'your_token'

// 翻页数据，假设已在其他地方定义
const flip_data = {
	book_id: book.id,
	now_page_num: Alpine.store('flip').nowPageNum,
	need_double_page_mode: false,
}

// 建立WebSocket连接的函数
function connectWebSocket() {
	// 根据当前协议选择ws或wss
	const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
	const wsUrl = wsProtocol + window.location.host + '/api/ws'
	socket = new WebSocket(wsUrl)

	// 连接打开时的回调
	socket.onopen = function () {
		console.log('WebSocket连接已建立')
		reconnectAttempts = 0 // 重置重连次数
	}

	// 收到消息时的回调
	socket.onmessage = function (event) {
		const message = JSON.parse(event.data)
		handleMessage(message) // 调用处理函数
	}

	// 连接关闭时的回调
	socket.onclose = function () {
		console.log('WebSocket连接已关闭')
		attemptReconnect() // 尝试重连
	}

	// 发生错误时的回调
	socket.onerror = function (error) {
		console.log('WebSocket发生错误：', error)
		socket.close() // 关闭连接以触发重连
	}
}

// 处理收到的翻页消息
function handleMessage(message) {
	// console.log("收到消息：", message);
	// console.log("My userID：" + userID);
	// console.log("Remote userID：" + message.user_id);
	// 根据消息类型进行处理
	if (message.type === 'flip_mode_sync_page' && message.user_id !== userID) {
		// 解析翻页数据
		const data = JSON.parse(message.data_string)
		if (Alpine.store('global').syncPageByWS && data.book_id === book.id) {
			//console.log("同步页数：", data);
			// 更新翻页数据
			flip_data.now_page_num = data.now_page_num
			// 更新页面
			jumpPageNum(data.now_page_num)
		}
	} else if (message.type === 'heartbeat') {
		console.log('收到心跳消息')
	} else {
		//console.log("不处理此消息"+message);
	}
}

// 发送翻页数据到服务器
function sendFlipData() {
	flip_data.now_page_num = Alpine.store('flip').nowPageNum
	const flipMsg = {
		type: 'flip_mode_sync_page', // 或 "heartbeat"
		status_code: 200,
		user_id: userID,
		token: token,
		detail: '翻页模式，发送数据',
		data_string: JSON.stringify(flip_data),
	}
	if (socket.readyState === WebSocket.OPEN) {
		socket.send(JSON.stringify(flipMsg))
		//console.log("已发送翻页数据"+JSON.stringify(flipMsg));
		//console.log("已发送翻页数据");
	} else {
		console.log('WebSocket未连接，无法发送消息')
	}
}

// 尝试重连函数
function attemptReconnect() {
	if (reconnectAttempts < maxReconnectAttempts) {
		reconnectAttempts++
		console.log(`第 ${reconnectAttempts} 次重连...`)
		setTimeout(() => {
			connectWebSocket()
		}, reconnectInterval)
	} else {
		console.log('已达到最大重连次数，停止重连')
	}
}

// 页面加载完成后建立WebSocket连接
window.onload = function () {
	connectWebSocket()
}
