/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.ibm.com/mbg-agent/cmd/cluster/state"
	"github.ibm.com/mbg-agent/pkg/clusterProxy"
	pb "github.ibm.com/mbg-agent/pkg/protocol"
	"google.golang.org/grpc"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect flow connection to the closest MBG",
	Long:  `connect flow connection to the closest MBG`,
	Run: func(cmd *cobra.Command, args []string) {
		svcId, _ := cmd.Flags().GetString("serviceId")
		svcIdDest, _ := cmd.Flags().GetString("serviceIdDest")
		svcPolicy, _ := cmd.Flags().GetString("policy")
		SendReq, _ := cmd.Flags().GetString("SendConnectReq")

		state.UpdateState()

		if svcId == "" || svcIdDest == "" {
			fmt.Println("Error: please insert all flag arguments for connect command")
			os.Exit(1)
		}
		svc := state.GetService(svcId)
		destSvc := state.GetService(svcIdDest)
		mbgIP := state.GetMbgIP()
		if SendReq == "true" {
			SendConnectReq(svcId, svcIdDest, svcPolicy, mbgIP)
		}
		connectClient(svc.Service.Ip, destSvc.Service.Ip)

	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().String("serviceId", "", "Service Id that the cluster is listen")
	connectCmd.Flags().String("serviceIdDest", "", "Destination service id the cluster is connecting")
	connectCmd.Flags().String("policy", "", "Connection policy")
	connectCmd.Flags().String("SendConnectReq", "true", "Decide if to send connection request to MBG default:True")

}

func connectClient(source, dest string) {
	var c clusterProxy.ProxyClient
	//TBD add validity check for the source and dest  IP
	c.InitClient(source, dest)
	c.RunClient()
}

func SendConnectReq(svcId, svcIdDest, svcPolicy, mbgIP string) {
	log.Printf("Start connect Request to MBG %v for service %v", svcIdDest, mbgIP)

	conn, err := grpc.Dial(mbgIP, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewConnectClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ConnectCmd(ctx, &pb.ConnectRequest{Id: svcId, IdDest: svcIdDest, Policy: svcPolicy})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf(`Response Connect message:  %s`, r.GetMessage())
}
