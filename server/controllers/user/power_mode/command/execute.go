package command

// POST: 执行能耗模式
// GET: 获取能耗模式执行历史

// SendCommandPower 发送能耗命令。
// 通过 client-id 拿到 Client。如果不存在，则报错。也意味着断开连接的客户端不执行。
// 向 Client 发送命令。
// 插入执行记录。
//func (c *Client) SendCommandPower(db *gorm.DB, data int) (int64, error) {
//	client := GlobalSessionManager.GetClient(c.ID)
//	if client == nil {
//		return 0, NewErrClientNotFound(c.ID)
//	}
//	command := NewEventCommandPower(data)
//	now := time.Now()
//	client.SendToSessionChannel(command)
//	// 发送命令后，将刚才发送的内容记录到数据表中。
//	return c.InsertNewCommandExecution(db, command.Code, command.MarshalData(), &now)
//}
