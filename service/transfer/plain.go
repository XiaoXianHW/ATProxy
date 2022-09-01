package transfer

import (
	"github.com/fatih/color"
	"github.com/xtls/xray-core/common/buf"
	"io"
	"log"
	"net"
	"runtime"
)

const (
	FLOW_ORIGIN = iota
	FLOW_LINUX_ZEROCOPY
	FLOW_ZEROCOPY
	FLOW_MULTIPLE
	FLOW_AUTO
)

type writerOnly struct {
	io.Writer
}

func SimpleTransfer(a, b net.Conn, flow int) {
	switch flow {
	case FLOW_ORIGIN:
		defer a.Close()
		defer b.Close()
		go io.Copy(writerOnly{b}, a)
		io.Copy(writerOnly{a}, b)

	case FLOW_ZEROCOPY:
		fallthrough

	case FLOW_LINUX_ZEROCOPY:
		if runtime.GOOS != "linux" {
			log.Panic(color.HiRedString("只有基于 Linux 的系统支持 Linux ZeroCopy，请将您的流程设置为 origin 或 auto。"))
		}
		fallthrough

	case FLOW_AUTO:
		if runtime.GOOS == "linux" {
			defer a.Close()
			defer b.Close()
			go io.Copy(b, a)
			io.Copy(a, b)
			return // TODO: Use MULTIPLE when fail to sendfile or splice
		}
		fallthrough

	case FLOW_MULTIPLE:
		aReader := buf.NewReader(a)
		bReader := buf.NewReader(b)
		aWriter := buf.NewWriter(a)
		bWriter := buf.NewWriter(b)
		defer a.Close()
		defer b.Close()
		go buf.Copy(bReader, aWriter)
		buf.Copy(aReader, bWriter)
	}
}
