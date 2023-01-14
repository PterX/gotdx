package internal

import (
	"gitee.com/quant1x/gotdx/proto/v2"
	"log"
)

// Command 命令字
func Command(pool *ConnPool, factory v2.Factory) v2.Unmarshaler {
	conn := pool.GetConn()

	cli := conn.(*Client)
	req, resp, err := factory()
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Do(req, resp)
	if err != nil {
		log.Println(err)
		_ = cli.Close()
		return nil
	}
	//log.Println(resp)
	pool.ReturnConn(conn)
	return resp
}
