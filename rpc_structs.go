package main

type ElectionArgs struct {
	AliveNodes []string
	OriginAddr string
}

type ElectionReply struct{}

type CoordinatorArgs struct {
	CoordinatorAddr string
	OriginAddr      string
}

type CoordinatorReply struct{}

type PingArgs struct {
}

type PingReply struct{}
