// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
)

func main() {

	routes, ciliumLinks, err := findRoutesAndLinks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}

	fmt.Println("%v", routes)
	fmt.Println("%v", ciliumLinks)

}

const (
	ciliumLinkPrefix  = "cilium_"
	ciliumNetNSPrefix = "cilium-"
	hostLinkPrefix    = "lxc"
	hostLinkLen       = len(hostLinkPrefix + "XXXXX")
)

func linkMatch(linkName string) bool {
	return strings.HasPrefix(linkName, ciliumLinkPrefix) ||
		strings.HasPrefix(linkName, hostLinkPrefix) && len(linkName) == hostLinkLen
}

func findRoutesAndLinks() (map[int]netlink.Route, map[int]netlink.Link, error) {
	routesToRemove := map[int]netlink.Route{}
	linksToRemove := map[int]netlink.Link{}

	if routes, err := netlink.RouteList(nil, netlink.FAMILY_V4); err == nil {
		for _, r := range routes {
			link, err := netlink.LinkByIndex(r.LinkIndex)
			if err != nil {
				fmt.Printf("zxy0 LinkByIndex error")
				if strings.Contains(err.Error(), "Link not found") {
					continue
				}
				return routesToRemove, linksToRemove, err
			}

			linkName := link.Attrs().Name
			if !linkMatch(linkName) {
				fmt.Printf("zxy1 link not match %s\n", linkName)
				continue
			}
			routesToRemove[r.LinkIndex] = r
			linksToRemove[link.Attrs().Index] = link
		}
	}

	if links, err := netlink.LinkList(); err == nil {
		for _, link := range links {
			linkName := link.Attrs().Name
			if !linkMatch(linkName) {
				fmt.Printf("zxy2 link not match  %s\n", linkName)
				continue
			}
			linksToRemove[link.Attrs().Index] = link
		}
	}
	return routesToRemove, linksToRemove, nil
}

