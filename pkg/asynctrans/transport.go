package asynctrans

import (
	"context"
)

type MsgHdr = string

const (
	MsgHdrMsgType MsgHdr = "MsgType"
	MsgHdrMsgName MsgHdr = "MsgName"
	MsgHdrGrpName MsgHdr = "GrpName"
	MsgHdrReplyTo MsgHdr = "ReplyTo"

	MsgHdrValUnk MsgHdr = "UNK"
)

type CtxKey int

const (
	CtxKeyMsgID CtxKey = iota
)

type MsgType = string

const (
	MsgTypeEvent MsgType = "EVENT"

	MsgTypeRq       MsgType = "RQ"
	MsgTypeRs       MsgType = "RS"
	MsgTypeErrRs    MsgType = "ERR_RS"
	MsgTypeErrEvent MsgType = "ERR_EVENT"

	MsgTypeUnk MsgType = "UNK"
)

type Handler interface {
	HandleRq(ctx context.Context, data []byte) ([]byte, error)
}

type MsgHandler func(ctx context.Context, data []byte) ([]byte, error)
