package access

import (
	"fmt"
	"github.com/XiaoXianHW/ATProxy/common/set"
	"github.com/XiaoXianHW/ATProxy/config"
)

func GetTargetList(listName string) (*set.StringSet, error) {
	set, ok := config.Lists[listName]
	if ok {
		return set, nil
	}
	return nil, fmt.Errorf("名单 %q 不存在", listName)
}
