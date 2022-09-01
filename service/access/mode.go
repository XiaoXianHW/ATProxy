package access

import "log"

const (
	DefaultMode = iota
	AllowMode
	BlockMode
)

func ParseAccessMode(mode string) int {
	switch mode {
	case "allow", "whitelist":
		return AllowMode
	case "block", "blacklist":
		return BlockMode
	case "":
		return DefaultMode
	default:
		log.Panicf("未知的模式: %q", mode)
		return 0
	}
}
