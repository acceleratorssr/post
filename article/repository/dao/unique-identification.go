package dao

import "github.com/bwmarrin/snowflake"

type UniqueID interface {
	Generate() uint64
}

type snowflakeNode struct {
	node *snowflake.Node
}

func (s *snowflakeNode) Generate() uint64 {
	return uint64(s.node.Generate().Int64())
}

func NewSnowflakeNode0() UniqueID {
	node, err := snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
	return &snowflakeNode{
		node: node,
	}
}
