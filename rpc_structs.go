package main

type ElectionArgs struct {
	AliveNodes []string
}

type ElectionReply struct{}

type CoordinatorArgs struct {
	CoordinatorAddr string
}

type CoordinatorReply struct{}

type PingArgs struct {
}

type PingReply struct{}
