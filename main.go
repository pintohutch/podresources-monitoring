package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	podresources "k8s.io/kubelet/pkg/apis/podresources/v1"
)

const (
	defaultPodResourcesSocket = "/var/lib/kubelet/pod-resources/kubelet.sock"
	defaultInterval           = 10 * time.Second
)

func prettyPrint(resp *podresources.ListPodResourcesResponse) {
	for _, p := range resp.PodResources {
		for _, c := range p.Containers {
			if len(c.Devices) == 0 && len(c.Memory) == 0 && len(c.DynamicResources) == 0 {
				continue
			}
			log.Printf("--------")
			log.Printf("pod: %s, namespace: %s, container: %s", p.Name, p.Namespace, c.Name)
			log.Printf("num of container devices: %d", len(c.Devices))
			for _, d := range c.Devices {
				log.Printf("..device resource name: %s", d.ResourceName)
				log.Printf("..device ids: %+v", d.DeviceIds)
				log.Printf("..device topology: %+v", d.Topology)
			}
			log.Printf("cpu ids: %+v", c.CpuIds)
			log.Printf("num of container memory assignments: %d", len(c.Memory))
			for _, m := range c.Memory {
				log.Printf("..memory type: %s", m.MemoryType)
				log.Printf("..memory size: %d", m.Size())
				log.Printf("..memory topology: %+v", m.Topology)
			}
			log.Printf("num of dynamic resources: %d", len(c.DynamicResources))
			for _, dr := range c.DynamicResources {
				log.Printf("..claim name: %s", dr.ClaimName)
				log.Printf("..claim namepace: %s", dr.ClaimNamespace)
				log.Printf("..claim resources: %+v", dr.ClaimResources)
			}
		}
	}
}

func main() {
	// Initialize flags.
	podresourcesSocket := flag.String("podresources-socket", defaultPodResourcesSocket, "Path to kubelet socket on disk.")
	interval := flag.Duration("interval", defaultInterval, "Polling interval for podresources api.")
	flag.Parse()

	// Create a channel to receive signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	// Ensure unix socket is available.
	if _, err := os.Stat(*podresourcesSocket); err != nil {
		log.Fatalf("podresources socket stat: %s", err)
	}

	// Initiailize grpc connection to socket.
	conn, err := grpc.NewClient(fmt.Sprintf("unix://%s", *podresourcesSocket),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("initializing grpc client against %q: %s", *podresourcesSocket, err)
	}
	defer conn.Close()

	// Initialize podresources API client.
	client := podresources.NewPodResourcesListerClient(conn)

	// Poll podresources API indefinitely.
	for {
		select {
		case sig := <-ch:
			fmt.Println("received signal:", sig)
			return // Exit the loop and the program
		default:
			// Continue with the loop
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.List(ctx, &podresources.ListPodResourcesRequest{})
			if err != nil {
				log.Printf("error listing pod resources: %s", err)
			} else {
				prettyPrint(resp)
			}
		}
		time.Sleep(*interval)
	}
}
