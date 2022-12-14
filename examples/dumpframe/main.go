package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"

	"github.com/shu1r0/srv6tracing_ebpfagent/pkg/ebpf"
)

func main() {
	dp, err := ebpf.NewTracingDataPlane(nil)
	if err != nil {
		panic(fmt.Errorf("Telemetry DP Error: %s", err))
	}
	defer dp.Close()

	if err := dp.AttachAll(); err != nil {
		panic(fmt.Errorf("Attach All Error: %s", err))
	}
	defer dp.DettachAll()

	if err := dp.SetMapConf(1); err != nil {
		panic(fmt.Errorf("Map Config Error: %s", err))
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	pktchan, err := dp.PacketInfoChan()
	if err != nil {
		panic(fmt.Errorf("Packet Info Error: %s", err))
	}
	go func() {
		for {
			pktinfo := <-pktchan
			fmt.Println("********** getPacket **********")
			fmt.Printf("Get Data: %s\n", hex.EncodeToString(pktinfo.Pkt))
			fmt.Printf("Packet : %s\n", hex.EncodeToString(pktinfo.Pkt[24:]))
			fmt.Printf("Packet ID : %b\n", pktinfo.PktId)
			fmt.Printf("Timestamp (mono): %b\n", pktinfo.MonotoricTimestamp)
			fmt.Printf("Hook: %d\n", pktinfo.Hookpoint)
			fmt.Println(hex.Dump(pktinfo.Pkt[24:]))
			//pkt := gopacket.NewPacket(pktinfo.Pkt[24:], layers.LayerTypeEthernet, gopacket.Default)
			//fmt.Println(pkt)
		}
	}()
	<-quit
}
