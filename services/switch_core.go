package services

// 交换服务核心模块

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/somebottle/localsend-switch/constants"
	"github.com/somebottle/localsend-switch/entities"
	switchdata "github.com/somebottle/localsend-switch/generated/switchdata/v1"
	"google.golang.org/protobuf/proto"
)

// handleTCP 处理单个 TCP 连接
//
// conn: TCP 连接
// dataChan: 传递接收到的交换数据的通道
// connectionVolume: 控制最大连接数的通道
// sigCtx: 中断信号上下文，用于优雅关闭连接
func handleTCP(conn *net.TCPConn, dataChan chan<- *switchdata.ClientInfo, connectionVolume chan struct{}, sigCtx context.Context) {
	// 用来向中断信号监听协程发送退出信号的管道
	handlerDone := make(chan struct{})
	// 处理完成后的清理
	defer func() {
		close(handlerDone)
		<-connectionVolume
		conn.Close()
	}()
	// 监听中断信号
	go func() {
		select {
		case <-sigCtx.Done():
			conn.Close()
		case <-handlerDone:
			// handleTCP 协程退出，这里也顺带退出
			return
		}
	}()
	// 设置连接的一些传输层属性
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(constants.TCPHeartbeatInterval * time.Second)
	// 接收数据
	buf := make([]byte, constants.TCPSocketReadBufferSize)
	for {
		// 设置读取超时，超过心跳时间没有数据就断开连接
		conn.SetReadDeadline(time.Now().Add(constants.TCPHeartbeatInterval * time.Second))
		// 每组数据传输格式: [ 1 字节的数据类型 | 4 字节的大端数据长度 | 数据 ]

		// 1 字节的数据类型
		//
		// 0x01 - ClientInfo 数据
		// 0x02 - 心跳包
		var dataType byte
		if err := binary.Read(conn, binary.BigEndian, &dataType); err != nil {
			// 读取类型失败，可能是连接出错 / 超时
			return
		}
		switch dataType {
		case 0x02:
			// 心跳包，什么都不做，继续等待下一个数据
			continue
		case 0x01:
			// ClientInfo 数据

			// 4 字节的数据长度
			var dataLength uint32
			if err := binary.Read(conn, binary.BigEndian, &dataLength); err != nil {
				// 读取长度失败，可能是连接出错
				return
			}
			if dataLength > constants.TCPSocketReadBufferSize {
				// 数据长度超过缓冲区大小，直接丢弃连接
				return
			}
			// 接下来读取 dataLength 字节的数据
			payload := buf[:dataLength]
			if _, err := io.ReadFull(conn, payload); err != nil {
				// 读取数据失败，可能是连接出错
				return
			}
			// 反序列化数据
			clientInfo := &switchdata.ClientInfo{}
			if err := proto.Unmarshal(payload, clientInfo); err != nil {
				// 反序列化失败，可能是数据格式错误，直接丢弃连接
				return
			}
			// 发送数据到通道
			dataChan <- clientInfo
		default:
			// 未知的数据类型，也是直接丢弃连接
			fmt.Printf("Unknown data type received over TCP: 0x%02X, closing connection\n", dataType)
			return
		}
	}
}

// receiveSwitchDataThroughTCP 通过 TCP 接收来自其他节点的交换数据
//
// servPort: 监听的服务端口
// dataChan: 传递接收到的交换数据的通道
// errChan: 传递错误信息的通道
// sigCtx: 中断信号上下文，用于优雅关闭服务
func receiveSwitchDataThroughTCP(servPort string, dataChan chan<- *switchdata.ClientInfo, errChan chan<- error, sigCtx context.Context) {
	for {
		// 端口转整数
		port, err := strconv.Atoi(servPort)
		if err != nil {
			errChan <- fmt.Errorf("Invalid service port: %v", err)
			return
		}
		// 控制最大连接数
		connectionVolume := make(chan struct{}, constants.MaxTCPConnections)
		exit, err := func() (bool, error) {
			// 启动 TCP 服务
			tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
				Port: port,
			})
			if err != nil {
				return true, err
			}
			defer tcpListener.Close()
			fmt.Printf("TCP Server listening on port %s\n", servPort)
			// 接受连接
			for {
				select {
				case <-sigCtx.Done():
					return true, nil
				default:
					conn, err := tcpListener.AcceptTCP()
					connectionVolume <- struct{}{}
					if err != nil {
						<-connectionVolume
						continue
					}
					// 处理连接
					go handleTCP(conn, dataChan, connectionVolume, sigCtx)
				}
			}
		}()
		if exit {
			// 收到退出信号
			if err != nil {
				errChan <- err
			}
			break
		}

		fmt.Printf("Restarting TCP Server...\nPrevious error: %v\n", err)
		time.Sleep(constants.TCPServerRestartInterval * time.Second)
	}
}

// SetUpSwitchCore 设置并启动交换服务核心模块
func SetUpSwitchCore(peerAddr string, peerPort string, servPort string, sigCtx context.Context, udpPacketChan <-chan entities.UDPPacketData, errChan chan<- error) {
	switchDataChan := make(chan *switchdata.ClientInfo, constants.SwitchDataReceiveChanSize)

	// 启动 TCP 服务以接收交换数据
	go receiveSwitchDataThroughTCP(servPort, switchDataChan, errChan, sigCtx)
}
