package collector

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/prometheus/client_golang/prometheus"
	"log"
    "fmt"
)

var (
	interfaceRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ovs_interface_rx_bytes_total",
			Help: "Total received bytes on an OVS interface",
		},
		[]string{"bridge", "port"},
	)
	interfaceTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ovs_interface_tx_bytes_total",
			Help: "Total transmitted bytes on an OVS interface",
		},
		[]string{"bridge", "port"},
	)
	interfaceRxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ovs_interface_rx_packets_total",
			Help: "Total received packets on an OVS interface",
		},
		[]string{"bridge", "port"},
	)
	interfaceTxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ovs_interface_tx_packets_total",
			Help: "Total transmitted packets on an OVS interface",
		},
		[]string{"bridge", "port"},
	)
	interfaceRxErrors = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ovs_interface_rx_errors_total",
			Help: "Total receive errors on an OVS interface",
		},
		[]string{"bridge", "port"},
	)
	interfaceTxDropped = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ovs_interface_tx_dropped_total",
			Help: "Total dropped TX packets on an OVS interface",
		},
		[]string{"bridge", "port"},
	)
)

func init() {
	prometheus.MustRegister(
		interfaceRxBytes,
		interfaceTxBytes,
		interfaceRxPackets,
		interfaceTxPackets,
		interfaceRxErrors,
		interfaceTxDropped,
	)
}

func CollectOvsMetrics() {
	client := ovs.New()

	bridges, err := client.VSwitch.ListBridges()
	if err != nil {
		log.Printf("could not list bridges: %v\n", err)
		return
	}

	for _, bridge := range bridges {
		ports, err := client.OpenFlow.DumpPorts(bridge)
		if err != nil {
			log.Printf("could not list ports on %s: %v\n", bridge, err)
			continue
		}

		for _, port := range ports {
			labels := prometheus.Labels{
				"bridge": bridge,
                "port": fmt.Sprintf("%d", port.PortID),
			}

			interfaceRxBytes.With(labels).Set(float64(port.Received.Bytes))
			interfaceTxBytes.With(labels).Set(float64(port.Transmitted.Bytes))
			interfaceRxPackets.With(labels).Set(float64(port.Received.Packets))
			interfaceTxPackets.With(labels).Set(float64(port.Transmitted.Packets))
			interfaceRxErrors.With(labels).Set(float64(port.Received.Errors))
			interfaceTxDropped.With(labels).Set(float64(port.Transmitted.Dropped))
		}
	}
}

