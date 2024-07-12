package dao

import "github.com/bwmarrin/snowflake"

func NewSnowflakeNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}
