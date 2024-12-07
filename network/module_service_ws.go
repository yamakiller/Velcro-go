package network

import (
	"context"
	"net/http"
	sync "sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/yamakiller/velcro-go/gofunc"
	"github.com/yamakiller/velcro-go/utils/circbuf"
	"github.com/yamakiller/velcro-go/vlog"
)

func newWSNetworkServerModule(system *NetworkSystem) *wsNetworkServerModule {
	return &wsNetworkServerModule{
		system: system,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type wsNetworkServerModule struct {
	system    *NetworkSystem
	server    *http.Server
	waitGroup sync.WaitGroup
	upgrader  *websocket.Upgrader
}

func (t *wsNetworkServerModule) Open(addr string) error {
	t.server = &http.Server{Addr: addr, Handler: t}
	gofunc.GoFunc(context.Background(), func() {
		vlog.Fatalf("VELCRO: network server listen failed, addr=%s error=%s", addr, t.server.ListenAndServe())
	})
	vlog.Infof("VELCRO: network server listen at addr=%s", addr)
	return nil
}

func (t *wsNetworkServerModule) Stop() {

}

func (tnc *wsNetworkServerModule) Network() string {
	return "wsserver"
}
func (t *wsNetworkServerModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := t.upgrader.Upgrade(w, r, nil)
	if err != nil {
		vlog.Errorf("ws upgrade: %v", err.Error())
		return
	}
	defer conn.Close()
	//新建客户端
	if err = t.spawn(conn); err != nil {
		return
	}
}

func (t *wsNetworkServerModule) spawn(conn *websocket.Conn) error {
	id := t.system.handlers.NextId()

	ctx := clientContext{system: t.system, state: stateAccept}
	handler := &wsClientHandler{
		conn:      conn,
		sendbox:   circbuf.NewLinkBuffer(4096),
		sendcond:  sync.NewCond(&sync.Mutex{}),
		keepalive: uint32(t.system.Config.Kleepalive),
		invoker:   &ctx,
		mailbox:   make(chan interface{}, 1),
		stopper:   make(chan struct{}),
		refdone:   &t.waitGroup,
	}

	cid, ok := t.system.handlers.Push(handler, id)
	if !ok {
		handler.Close()
		// 释放资源
		close(handler.mailbox)
		close(handler.stopper)
		handler.sendbox.Close()
		handler.sendbox = nil
		handler.sendcond = nil

		return errors.Errorf("client-id %s existed", cid.ToString())
	}

	ctx.self = cid
	ctx.incarnateClient()

	handler.start()
	return nil
}
