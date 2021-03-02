package middleware

import (
	"fmt"
	"strings"

	core "onebot/core/xianqu"
)

// 将报文中Response的message转换为array格式
func ResponseToArray(ctx *core.Context) {
	message := ctx.Response["message"]
	switch message.(type) {
	case string:
		//
	default:
		return
	}
	ctx.Response["message"] = CqCode2Array(message.(string))
	return
}

// 将报文中Request的message转换为array格式
func RequestToArray(ctx *core.Context) {
	fmt.Println(ctx)
	if !core.Parse(ctx.Request).Get("params").Exist("message") {
		return
	}
	message := core.Parse(ctx.Request).Get("params").Str("message")
	fmt.Println(message)
	if message == "" {
		return
	}
	switch ctx.Request["params"].(type) {
	case map[string]interface{}:
		ctx.Request["params"].(map[string]interface{})["message"] = CqCode2Array(message)
	}
	return
}

func CqCode2Array(text string) []map[string]interface{} {
	elems := SplitCQText(text)
	var (
		array = []map[string]interface{}{}

		isCQcode bool = false
		isFirst  bool = false
		isType   bool = false
		isKey    bool = false
		isValue  bool = false

		cqCodetype_ string
		cqCodeKey   []string
		cqCodeValue []string
	)
	for _, r := range elems {
		var temp = []rune{}
	elemLoop:
		for i := range r {
			switch {
			// TODO CQ码开始标记 []rune("[CQ:") = [91 67 81 58]
			case r[i] == 91 && r[i+1] == 67 && r[i+2] == 81 && r[i+3] == 58:
				isCQcode = true
				isFirst = true
			// TODO 不是CQ码
			case !isCQcode:
				array = append(
					array,
					map[string]interface{}{
						"type": "text",
						"data": map[string]interface{}{
							"text": string(r),
						},
					},
				)
				break elemLoop
			// TODO type字段标记 []rune(":") = [58]
			case isCQcode && isFirst && r[i] == 58:
				isType = true
				isFirst = false
			// TODO 开始装载type字段 []rune(",") = [44]
			case isType && r[i] != 44:
				temp = append(temp, r[i])
			// TODO 结束装载type字段 key字段标记 []rune(",") = [44]
			case isType && r[i] == 44:
				cqCodetype_ = string(temp)
				temp = []rune{}
				isType = false
				isKey = true
			// TODO 开始装载key []rune("=") = [61]
			case isKey && r[i] != 61:
				temp = append(temp, r[i])
			// TODO 结束装载key字段 value字段标记 []rune("=") = [61]
			case isKey && r[i] == 61:
				cqCodeKey = append(cqCodeKey, string(temp))
				temp = []rune{}
				isKey = false
				isValue = true
			// TODO 开始装载value []rune(",") = [44] []rune("]") = [93]
			case isValue && r[i] != 44 && r[i] != 93:
				temp = append(temp, r[i])
			// TODO 结束装载value字段 key字段标记 []rune(",") = [44]
			case isValue && r[i] == 44:
				cqCodeValue = append(cqCodeValue, string(temp))
				temp = []rune{}
				isValue = false
				isKey = true
			// TODO 结束装载value字段 结束CQ码 []rune("]") = [93]
			case isValue && r[i] == 93:
				cqCodeValue = append(cqCodeValue, string(temp))
				temp = []rune{}
				cqCodeMap := map[string]interface{}{}
				cqCodeMap["type"] = cqCodetype_
				keyValue := map[string]interface{}{}
				for i := range cqCodeKey {
					keyValue[cqCodeKey[i]] = cqCodeValue[i]
				}
				cqCodeMap["data"] = keyValue
				array = append(array, cqCodeMap)
				cqCodeKey = []string{}
				cqCodeValue = []string{}
				isValue = false
				isCQcode = false
			default:

				// TODO do nothing
			}
		}
	}
	return array
}

func SplitCQText(cqcode string) [][]rune {
	var (
		elems    [][]rune
		temp     []rune
		isCQcode bool = false
	)
	r := []rune(cqcode)
	for i := range r {
		switch {
		// TODO CQ码开始标记 []rune("[CQ:") = [91 67 81 58]
		case r[i] == 91 && r[i+1] == 67 && r[i+2] == 81 && r[i+3] == 58:
			isCQcode = true
			elems = append(elems, temp)
			// TODO 清空temp，开始装CQ码
			temp = []rune{}
			temp = append(temp, r[i])
		// TODO CQ码一直到出现"]"为止 []rune("]") = [93]
		case isCQcode && r[i] != 93:
			temp = append(temp, r[i])
		// TODO 出现"]"，开始下一个elem
		case isCQcode && r[i] == 93:
			isCQcode = false
			temp = append(temp, r[i])
			elems = append(elems, temp)
			temp = []rune{}
		default:
			temp = append(temp, r[i])
		case i == len(r)-1:
			temp = append(temp, r[i])
			elems = append(elems, temp)
		}
	}
	return elems
}

func escape(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&#44;", ",")
	text = strings.ReplaceAll(text, "&#91;", "[")
	text = strings.ReplaceAll(text, "&#93;", "]")
	return text
}
