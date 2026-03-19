package utils

import "strings"

// MaskPhone 手机号脱敏：显示前3后4
// +8613812345678 -> +86138****5678
// 13812345678 -> 138****5678
func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}

	// E.164 格式：+国家码+号码
	if strings.HasPrefix(phone, "+") {
		// +8613812345678 (13位) -> +86138****5678
		if len(phone) >= 11 {
			return phone[:6] + "****" + phone[len(phone)-4:]
		}
		if len(phone) > 4 {
			return phone[:3] + "****" + phone[len(phone)-2:]
		}
		return "***"
	}

	// 国内格式：13812345678 (11位) -> 138****5678
	if len(phone) >= 11 {
		return phone[:3] + "****" + phone[len(phone)-4:]
	}

	// 短号码
	if len(phone) > 7 {
		return phone[:3] + "****" + phone[len(phone)-4:]
	}

	return "***"
}

// MaskName 姓名脱敏：隐藏姓氏
// 张三 -> *三
// 李四五 -> *四五
func MaskName(name string) string {
	if name == "" {
		return ""
	}

	runes := []rune(name)
	if len(runes) <= 1 {
		return "*"
	}

	// 隐藏第一个字（姓氏）
	return "*" + string(runes[1:])
}

// MaskIDCard 身份证脱敏：显示前4后4
// 110101199001011234 -> 1101**********1234
func MaskIDCard(id string) string {
	if id == "" {
		return ""
	}

	if len(id) == 18 {
		return id[:4] + "**********" + id[14:]
	}

	// 非标准长度，至少8位才脱敏
	if len(id) >= 8 {
		return id[:4] + "****" + id[len(id)-4:]
	}

	return "***masked***"
}
