package gateway

import (
	"encoding/json"

	core "onebot/core/xianqu"
	middle "onebot/middleware"
	ser "onebot/server"
)

// broadcast 将不同bot的context广播到该bot的所有连接
func broadcast(ctx *core.Context) {
	INFO("[上报] %v", ctx)
	// OneBot标准中的报文string格式
	data, _ := json.Marshal(ctx.Response)
	// OneBot标准中的报文array格式
	newCtx := core.CopyCtx(ctx)
	middle.ResponseToArray(newCtx)
	data2, _ := json.Marshal(newCtx.Response)
	// 遍历所有的HttpPost连接的服务端
	for _, s := range Connects[ctx.Bot].HttpServers {
		if s.Conn.SendStatus == ser.SOK {
			if s.PostMessageFormat == "string" {
				s.SendChan <- data
				continue
			}
			s.SendChan <- data2
		}
	}
	// 遍历所有的反向WS连接的服务端
	for _, s := range Connects[ctx.Bot].WebSocketClients {
		if s.Conn.SendStatus == ser.SOK {
			if s.PostMessageFormat == "string" {
				s.SendChan <- data
				continue
			}
			s.SendChan <- data2
		}
	}

	Connects[ctx.Bot].WebSocketServers.Broadcast(data)
}

// callapi 将context分发到core的各个API处
func callapi(ctx *core.Context) {
	INFO("[响应] %v", ctx)
	switch ctx.Request["action"].(string) {
	case "send_private_msg":
		core.ApiSendPrivateMsg(ctx)
	case "send_group_msg":
		core.ApiSendGroupMsg(ctx)
	case "send_msg":
		core.ApiSendMsg(ctx)
	case "delete_msg":
		core.ApiDeleteMsg(ctx)
	case "get_msg":
		core.ApiGetMsg(ctx)
	case "get_forward_msg":
		core.ApiGetForwardMsg(ctx)
	case "send_like":
		core.ApiSendLike(ctx)
	case "set_group_kick":
		core.ApiSetGroupKick(ctx)
	case "set_group_ban":
		core.ApiSetGroupBan(ctx)
	case "set_group_anonymous_ban":
		core.ApiSetGroupAnonymousBan(ctx)
	case "set_group_whole_ban":
		core.ApiSetGroupWholeBan(ctx)
	case "set_group_admin":
		core.ApiSetGroupAdmin(ctx)
	case "set_group_anonymous":
		core.ApiSetGroupAnonymous(ctx)
	case "set_group_card":
		core.ApiSetGroupCard(ctx)
	case "set_group_name":
		core.ApiSetGroupName(ctx)
	case "set_group_leave":
		core.ApiSetGroupLeave(ctx)
	case "set_group_special_title":
		core.ApiSetGroupSpecialTitle(ctx)
	case "set_friend_add_request":
		core.ApiSetFriendAddRequest(ctx)
	case "set_group_add_request":
		core.ApiSetGroupAddRequest(ctx)
	case "get_login_info":
		core.ApiGetLoginInfo(ctx)
	case "get_stranger_info":
		core.ApiGetStrangerInfo(ctx)
	case "get_friend_list":
		core.ApiGetFriendList(ctx)
	case "get_group_info":
		core.ApiGetGroupInfo(ctx)
	case "get_group_list":
		core.ApiGetGroupList(ctx)
	case "get_group_member_info":
		core.ApiGetGroupMemberInfo(ctx)
	case "get_group_member_list":
		core.ApiGetGroupMemberList(ctx)
	case "get_group_honor_info":
		core.ApiGetGroupHonorInfo(ctx)
	case "get_cookies":
		core.ApiGetCookies(ctx)
	case "get_csrf_token":
		core.ApiGetCsrfToken(ctx)
	case "get_credentials":
		core.ApiGetCredentials(ctx)
	case "get_record":
		core.ApiGetRecord(ctx)
	case "get_image":
		core.ApiGetImage(ctx)
	case "can_send_image":
		core.ApiCanSendImage(ctx)
	case "can_send_record":
		core.ApiCanSendRecord(ctx)
	case "get_status":
		core.ApiGetStatus(ctx)
	case "get_version_info":
		core.ApiGetVersionInfo(ctx)
	case "set_restart":
		core.ApiSetRestart(ctx)
	case "clean_cache":
		core.ApiCleanCache(ctx)
	}
}
