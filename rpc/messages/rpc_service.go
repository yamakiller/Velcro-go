// Code generated by Thrift Compiler (0.19.0). DO NOT EDIT.

package messages

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
	thrift "github.com/apache/thrift/lib/go/thrift"
	"strings"
	"regexp"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = errors.New
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal
// (needed by validator.)
var _ = strings.Contains
var _ = regexp.MatchString

type RpcService interface {}

type RpcServiceClient struct {
  c thrift.TClient
  meta thrift.ResponseMeta
}

func NewRpcServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *RpcServiceClient {
  return &RpcServiceClient{
    c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
  }
}

func NewRpcServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *RpcServiceClient {
  return &RpcServiceClient{
    c: thrift.NewTStandardClient(iprot, oprot),
  }
}

func NewRpcServiceClient(c thrift.TClient) *RpcServiceClient {
  return &RpcServiceClient{
    c: c,
  }
}

func (p *RpcServiceClient) Client_() thrift.TClient {
  return p.c
}

func (p *RpcServiceClient) LastResponseMeta_() thrift.ResponseMeta {
  return p.meta
}

func (p *RpcServiceClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
  p.meta = meta
}

type RpcServiceProcessor struct {
  processorMap map[string]thrift.TProcessorFunction
  handler RpcService
}

func (p *RpcServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
  p.processorMap[key] = processor
}

func (p *RpcServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
  processor, ok = p.processorMap[key]
  return processor, ok
}

func (p *RpcServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
  return p.processorMap
}

func NewRpcServiceProcessor(handler RpcService) *RpcServiceProcessor {

  self0 := &RpcServiceProcessor{handler:handler, processorMap:make(map[string]thrift.TProcessorFunction)}
return self0
}

func (p *RpcServiceProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
  name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
  if err2 != nil { return false, thrift.WrapTException(err2) }
  if processor, ok := p.GetProcessorFunction(name); ok {
    return processor.Process(ctx, seqId, iprot, oprot)
  }
  iprot.Skip(ctx, thrift.STRUCT)
  iprot.ReadMessageEnd(ctx)
  x1 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function " + name)
  oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
  x1.Write(ctx, oprot)
  oprot.WriteMessageEnd(ctx)
  oprot.Flush(ctx)
  return false, x1

}


// HELPER FUNCTIONS AND STRUCTURES


