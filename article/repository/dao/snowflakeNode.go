package dao

import "github.com/bwmarrin/snowflake"

func NewSnowflakeNode0() *snowflake.Node {
	node, err := snowflake.NewNode(0)
	if err != nil {
		panic(err)
	}
	return node
}
