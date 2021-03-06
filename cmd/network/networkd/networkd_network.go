// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Address struct {
	Address string `json:"Address"`
	Peer    string `json:"Peer"`
	Label   string `json:"Label"`
	Scope   string `json:"Scope"`
}

type Route struct {
	Gateway         string `json:"Gateway"`
	GatewayOnlink   string `json:"GatewayOnlink"`
	Destination     string `json:"Destination"`
	Source          string `json:"Source"`
	PreferredSource string `json:"PreferredSource"`
	Table           string `json:"Table"`
}

type Network struct {
	ConfFile string `json:"ConfFile"`

	Match     interface{} `json:"Match"`
	Addresses interface{} `json:"Addresses"`
	Routes    interface{} `json:"Routes"`

	Gateway             string `json:"Gateway"`
	DHCP                string `json:"DHCP"`
	DNS                 string `json:"DNS"`
	Domains             string `json:"Domains"`
	NTP                 string `json:"NTP"`
	IPv6AcceptRA        string `json:"IPv6AcceptRA"`
	LinkLocalAddressing string `json:"LinkLocalAddressing"`
	LLDP                string `json:"LLDP"`
	EmitLLDP            string `json:"EmitLLDP"`

	Bridge  string `json:"Bridge"`
	Bond    string `json:"Bond"`
	VRF     string `json:"VRF"`
	VLAN    string `json:"VLAN"`
	MACVLAN string `json:"MACVLAN"`
	VXLAN   string `json:"VXLAN"`
	Tunnel  string `json:"Tunnel"`
}

func (network *Network) CreateNetworkMatchSectionConfig() string {
	conf := "[Match]\n"

	switch v := network.Match.(type) {
	case []interface{}:
		for _, b := range v {
			var mac string
			var driver string
			var name string

			if b.(map[string]interface{})["MAC"] != nil {
				mac = strings.TrimSpace(b.(map[string]interface{})["MAC"].(string))
			}

			if b.(map[string]interface{})["Driver"] != nil {
				driver = strings.TrimSpace(b.(map[string]interface{})["Driver"].(string))
			}

			if b.(map[string]interface{})["Name"] != nil {
				name = strings.TrimSpace(b.(map[string]interface{})["Name"].(string))
			}

			if mac != "" {
				mac := fmt.Sprintf("MACAddress=%s\n", mac)
				conf += mac
			}

			if driver != "" {
				driver := fmt.Sprintf("Driver=%s\n", driver)
				conf += driver
			}

			if name != "" {
				if network.ConfFile == "" {
					network.ConfFile = name
				}

				name := fmt.Sprintf("Name=%s\n", name)
				conf += name
			}
		}
		break
	}

	fmt.Println(conf)
	return conf
}

func (network *Network) CreateRouteSectionConfig() string {
	var routeConf string

	switch v := network.Routes.(type) {
	case []interface{}:
		for _, b := range v {
			var preferredSource string
			var gatewayOnLink string
			var destination string
			var gateway string
			var source string
			var table string

			if b.(map[string]interface{})["Gateway"] != nil {
				gateway = strings.TrimSpace(b.(map[string]interface{})["Gateway"].(string))
			}

			if b.(map[string]interface{})["GatewayOnlink"] != nil {
				gatewayOnLink = strings.TrimSpace(b.(map[string]interface{})["GatewayOnlink"].(string))
			}

			if b.(map[string]interface{})["Destination"] != nil {
				destination = strings.TrimSpace(b.(map[string]interface{})["Destination"].(string))
			}

			if b.(map[string]interface{})["Source"] != nil {
				source = strings.TrimSpace(b.(map[string]interface{})["Source"].(string))
			}

			if b.(map[string]interface{})["PreferredSource"] != nil {
				preferredSource = strings.TrimSpace(b.(map[string]interface{})["PreferredSource"].(string))
			}

			if b.(map[string]interface{})["Table"] != nil {
				table = strings.TrimSpace(b.(map[string]interface{})["Table"].(string))
			}

			routeConf += "\n[Route]\n"

			if len(gateway) != 0 {
				ip := net.ParseIP(gateway)
				if ip != nil {
					routeConf += "Gateway=" + gateway + "\n"
				} else {
					log.Error("Failed to parse Gateway: ", gateway)
				}
			}

			if len(gatewayOnLink) != 0 {
				onlink := strings.TrimSpace(gatewayOnLink)
				b, r := share.ParseBool(onlink)
				if r != nil {
					log.Error("Failed to parse GatewayOnlink: ", r, gatewayOnLink)
				} else {
					if b == true {
						routeConf += "GatewayOnlink=yes\n"
					} else {
						routeConf += "GatewayOnlink=no\n"
					}
				}
			}

			if len(destination) != 0 {
				ip := net.ParseIP(destination)
				if ip != nil {
					routeConf += "Destination=" + destination + "\n"
				} else {
					log.Error("Failed to parse Destination: ", destination)
				}
			}

			if len(source) != 0 {
				ip := net.ParseIP(source)
				if ip != nil {
					routeConf += "Source=" + source + "\n"
				} else {
					log.Error("Failed to parse Source: ", source)
				}
			}

			if len(preferredSource) != 0 {
				ip := net.ParseIP(preferredSource)
				if ip != nil {
					routeConf += "PreferredSource=" + preferredSource + "\n"
				} else {
					log.Error("Failed to parse PreferredSource: ", preferredSource)
				}
			}

			if len(table) != 0 {
				routeConf += "Table=" + table + "\n"
			}
		}
		break
	}

	return routeConf
}

func (network *Network) CreateAddressSectionConfig() string {
	var addressConf string

	switch v := network.Addresses.(type) {
	case []interface{}:
		for _, b := range v {
			var address string
			var peer string
			var scope string
			var label string

			if b.(map[string]interface{})["Address"] != nil {
				address = strings.TrimSpace(b.(map[string]interface{})["Address"].(string))
			}

			if b.(map[string]interface{})["Peer"] != nil {
				peer = strings.TrimSpace(b.(map[string]interface{})["Peer"].(string))
			}

			if b.(map[string]interface{})["Scope"] != nil {
				scope = strings.TrimSpace(b.(map[string]interface{})["Scope"].(string))
			}

			if b.(map[string]interface{})["Label"] != nil {
				label = strings.TrimSpace(b.(map[string]interface{})["Label"].(string))
			}

			if len(address) != 0 {
				addressConf += "\n[Address]\nAddress="

				ip := net.ParseIP(address)
				if ip != nil {
					addressConf += address + "\n"
				} else {
					log.Error("Failed to parse address: ", address)
				}

				if len(peer) != 0 {
					ip = net.ParseIP(peer)
					if ip != nil {
						addressConf += "Peer=" + peer + "\n"
					} else {
						log.Error("Failed to parse peer address: ", peer)
					}
				}

				if len(scope) != 0 {
					addressConf += "Scope=" + scope + "\n"
				}

				if len(label) != 0 {
					addressConf += "Label=" + label + "\n"
				}
			}
		}
		break
	}

	return addressConf
}

func (network *Network) CreateNetworkSectionConfig() string {
	conf := "[Network]\n"

	if network.DHCP != "" {
		dhcpConf := "DHCP="

		dhcp := strings.TrimSpace(network.DHCP)
		b, r := share.ParseBool(dhcp)
		if r != nil {
			switch dhcp {
			case "ipv4", "ipv6":
				dhcpConf += dhcp
				break
			default:
				log.Error("Failed to parse DHCP: ", r, network.DHCP)
			}
		} else {
			if b == true {
				dhcpConf += "yes"
			} else {
				dhcpConf += "no"
			}
		}
		conf += dhcpConf + "\n"
	}

	if network.Gateway != "" {
		gatewayConf := "Gateway="

		gw := strings.TrimSpace(network.Gateway)
		ip := net.ParseIP(gw)
		if ip != nil {
			gatewayConf += gw
			conf += gatewayConf + "\n"
		} else {
			log.Error("Failed to parse Gateway Address: ", network.Gateway)
		}
	}

	if network.DNS != "" {
		conf += "DNS=" + network.DNS
	}

	if network.Domains != "" {
		conf += "Domains=" + network.Domains + "\n"
	}

	if network.NTP != "" {
		conf += "NTP=" + network.NTP + "\n"
	}

	if network.IPv6AcceptRA != "" {
		IPv6AcceptRAConf := "IPv6AcceptRA="

		IPv6RA := strings.TrimSpace(network.IPv6AcceptRA)
		b, err := share.ParseBool(IPv6RA)
		if err != nil {
			log.Error("Failed to parse IPv6AcceptRA: ", err, network.IPv6AcceptRA)
		} else {
			if b == true {
				IPv6AcceptRAConf += "yes"
			} else {
				IPv6AcceptRAConf += "no"
			}
		}
		conf += IPv6AcceptRAConf + "\n"
	}

	if network.LinkLocalAddressing != "" {
		LinkLocalAddressingConf := "LinkLocalAddressing="

		IPv6RA := strings.TrimSpace(network.LinkLocalAddressing)
		b, err := share.ParseBool(IPv6RA)
		if err != nil {
			log.Error("Failed to parse LinkLocalAddressing: ", err, network.LinkLocalAddressing)
		} else {
			if b == true {
				LinkLocalAddressingConf += "yes"
			} else {
				LinkLocalAddressingConf += "no"
			}
		}
		conf += LinkLocalAddressingConf + "\n"
	}

	if network.LLDP != "" {
		LLDPConf := "LLDP="

		LLDP := strings.TrimSpace(network.LLDP)
		b, err := share.ParseBool(LLDP)
		if err != nil {
			log.Error("Failed to parse LLDP: ", err, network.LLDP)
		} else {
			if b == true {
				LLDPConf += "yes"
			} else {
				LLDPConf += "no"
			}
		}
		conf += LLDPConf + "\n"
	}

	if network.EmitLLDP != "" {
		EmitLLDPConf := "EmitLLDP="

		EmitLLDP := strings.TrimSpace(network.EmitLLDP)
		b, err := share.ParseBool(EmitLLDP)
		if err != nil {
			log.Error("Failed to parse EmitLLDP: ", err, network.EmitLLDP)
		} else {
			if b == true {
				EmitLLDPConf += "yes"
			} else {
				EmitLLDPConf += "no"
			}
		}
		conf += EmitLLDPConf + "\n"
	}

	if network.NTP != "" {
		conf += "NTP=" + network.NTP + "\n"
	}

	if network.Bridge != "" {
		conf += "Bridge=" + network.Bridge + "\n"
	}

	if network.Bond != "" {
		conf += "Bond=" + network.Bond + "\n"
	}

	if network.VLAN != "" {
		conf += "VLAN=" + network.VLAN + "\n"
	}

	return conf
}

func NetworkdParseJSONfromHTTPReq(req *http.Request) error {
	var configs map[string]interface{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %s", err)
		return err
	}

	json.Unmarshal([]byte(body), &configs)

	network := new(Network)
	json.Unmarshal([]byte(body), &network)

	matchConfig := network.CreateNetworkMatchSectionConfig()
	networkConfig := network.CreateNetworkSectionConfig()
	addressConfig := network.CreateAddressSectionConfig()
	routeConfig := network.CreateRouteSectionConfig()

	config := []string{matchConfig, networkConfig, addressConfig, routeConfig}

	fmt.Println(config)

	unitName := fmt.Sprintf("25-%s.network", network.ConfFile)
	unitPath := filepath.Join(NetworkdUnitPath, unitName)

	return share.WriteFullFile(unitPath, config)
}

func ConfigureNetworkFile(rw http.ResponseWriter, req *http.Request) {
	NetworkdParseJSONfromHTTPReq(req)
}
