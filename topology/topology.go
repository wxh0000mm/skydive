/*
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package topology

import (
	"sync"
)

const (
	Root        = "root"
	NetNs       = "netns"
	LinuxBridge = "linuxbridge"
	OvsBridge   = "ovsbridge"
)

type Node struct {
	sync.RWMutex
	ID        string            `json:"-"`
	Metadatas map[string]string `json:",omitempty"`
	links     map[string]*Link  `json:"-"`
}

type Link struct {
	sync.RWMutex
	ID        string `json:"-"`
	Left      *Node
	Right     *Node
	Metadatas map[string]string `json:",omitempty"`
}

type Container struct {
	sync.RWMutex
	ID       string `json:"-"`
	Type     string
	Nodes    map[string]*Node
	Topology *Topology `json:"-"`
}

type Topology struct {
	sync.RWMutex
	Containers map[string]*Container
}

func (c *Container) GetNode(i string) *Node {
	c.Lock()
	defer c.Unlock()

	if node, ok := c.Nodes[i]; ok {
		return node
	}
	return nil
}

func (c *Container) DelNode(i string) {
	c.Lock()
	defer c.Unlock()

	n, ok := c.Nodes[i]
	if !ok {
		return
	}

	for id, link := range n.links {
		delete(link.Left.links, id)
		delete(link.Right.links, id)
	}

	delete(c.Nodes, i)
}

func (c *Container) NewNode(i string) *Node {
	c.Lock()
	defer c.Unlock()

	node := &Node{
		ID:        i,
		Metadatas: make(map[string]string),
		links:     make(map[string]*Link),
	}
	c.Nodes[i] = node

	return node
}

func (topo *Topology) DelContainer(i string) {
	topo.Lock()
	defer topo.Unlock()

	delete(topo.Containers, i)
}

func (topo *Topology) GetContainer(i string) *Container {
	c, ok := topo.Containers[i]
	if !ok {
		return nil
	}
	return c
}

func (topo *Topology) NewContainer(i string, t string) *Container {
	topo.Lock()
	defer topo.Unlock()

	container := &Container{
		ID:       i,
		Type:     t,
		Nodes:    make(map[string]*Node),
		Topology: topo,
	}
	topo.Containers[i] = container

	return container
}

func NewTopology() *Topology {
	return &Topology{
		Containers: make(map[string]*Container),
	}
}
