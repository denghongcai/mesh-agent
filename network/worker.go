package network

import (
	"net"
	"log"
	"bufio"
	"sync"
)

type Worker struct {
	conn *net.TCPConn
	addr string
	codecConfig *CodecConfig
	closeChan chan bool
	closed sync.Once
}

func NewWorker(conn *net.TCPConn, config *CodecConfig) *Worker {
	return &Worker{conn: conn, codecConfig:config,closeChan:make(chan bool)}
}

func (w *Worker) Close() {
	w.closed.Do(w.close)
}

func (w *Worker) close() {
	w.conn.Close()
	close(w.closeChan)
}

func (w *Worker) Run() {
	go w.readHandler()
	go w.writeHandler()
	<- w.closeChan
}

func (w *Worker) readHandler() {
	defer w.conn.Close()
	for {
		reader := bufio.NewReader(w.conn)
		// Do Decode
		req, err := w.codecConfig.Decoder.DecodeRequest(reader)
		if err != nil {
			log.Println(err)
			w.Close()
			break
		}
		select {
		case w.codecConfig.ReadChan <- req:
		case <- w.closeChan:
			break
		}
	}
}

func (w *Worker) writeHandler() {
	defer w.conn.Close()
	for {
		select {
		case msg := <-w.codecConfig.WriteChan:
			res, err := w.codecConfig.Encoder.EncodeResponse(msg)
			if err != nil {
				log.Println(err)
				break
			}
			_, err = w.conn.Write(res)
			nerr, ok := err.(*net.OpError)
			if ok && nerr.Temporary() {
				// should retry
				w.Close()
				break
			}
			if err != nil {
				// should retry
				w.Close()
				break
			}
		case <- w.closeChan:
			break
		}

	}
}