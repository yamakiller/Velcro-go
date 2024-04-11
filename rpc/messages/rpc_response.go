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

// Attributes:
//  - SequenceID
//  - Result_
type RpcResponseMessage struct {
  SequenceID int32 `thrift:"SequenceID,1" db:"SequenceID" json:"SequenceID"`
  Result_ []byte `thrift:"Result,2" db:"Result" json:"Result"`
}

func NewRpcResponseMessage() *RpcResponseMessage {
  return &RpcResponseMessage{}
}


func (p *RpcResponseMessage) GetSequenceID() int32 {
  return p.SequenceID
}

func (p *RpcResponseMessage) GetResult_() []byte {
  return p.Result_
}
func (p *RpcResponseMessage) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.I32 {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    case 2:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField2(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *RpcResponseMessage)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadI32(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.SequenceID = v
}
  return nil
}

func (p *RpcResponseMessage)  ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadBinary(ctx); err != nil {
  return thrift.PrependError("error reading field 2: ", err)
} else {
  p.Result_ = v
}
  return nil
}

func (p *RpcResponseMessage) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "RpcResponseMessage"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
    if err := p.writeField2(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *RpcResponseMessage) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "SequenceID", thrift.I32, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:SequenceID: ", p), err) }
  if err := oprot.WriteI32(ctx, int32(p.SequenceID)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.SequenceID (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:SequenceID: ", p), err) }
  return err
}

func (p *RpcResponseMessage) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "Result", thrift.STRING, 2); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:Result: ", p), err) }
  if err := oprot.WriteBinary(ctx, p.Result_); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.Result (2) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 2:Result: ", p), err) }
  return err
}

func (p *RpcResponseMessage) Equals(other *RpcResponseMessage) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.SequenceID != other.SequenceID { return false }
  if bytes.Compare(p.Result_, other.Result_) != 0 { return false }
  return true
}

func (p *RpcResponseMessage) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("RpcResponseMessage(%+v)", *p)
}

func (p *RpcResponseMessage) Validate() error {
  return nil
}
// Attributes:
//  - Err
type RpcError struct {
  Err string `thrift:"Err,1" db:"Err" json:"Err"`
}

func NewRpcError() *RpcError {
  return &RpcError{}
}


func (p *RpcError) GetErr() string {
  return p.Err
}
func (p *RpcError) Read(ctx context.Context, iprot thrift.TProtocol) error {
  if _, err := iprot.ReadStructBegin(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
  }


  for {
    _, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
    if err != nil {
      return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
    }
    if fieldTypeId == thrift.STOP { break; }
    switch fieldId {
    case 1:
      if fieldTypeId == thrift.STRING {
        if err := p.ReadField1(ctx, iprot); err != nil {
          return err
        }
      } else {
        if err := iprot.Skip(ctx, fieldTypeId); err != nil {
          return err
        }
      }
    default:
      if err := iprot.Skip(ctx, fieldTypeId); err != nil {
        return err
      }
    }
    if err := iprot.ReadFieldEnd(ctx); err != nil {
      return err
    }
  }
  if err := iprot.ReadStructEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
  }
  return nil
}

func (p *RpcError)  ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
  if v, err := iprot.ReadString(ctx); err != nil {
  return thrift.PrependError("error reading field 1: ", err)
} else {
  p.Err = v
}
  return nil
}

func (p *RpcError) Write(ctx context.Context, oprot thrift.TProtocol) error {
  if err := oprot.WriteStructBegin(ctx, "RpcError"); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err) }
  if p != nil {
    if err := p.writeField1(ctx, oprot); err != nil { return err }
  }
  if err := oprot.WriteFieldStop(ctx); err != nil {
    return thrift.PrependError("write field stop error: ", err) }
  if err := oprot.WriteStructEnd(ctx); err != nil {
    return thrift.PrependError("write struct stop error: ", err) }
  return nil
}

func (p *RpcError) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
  if err := oprot.WriteFieldBegin(ctx, "Err", thrift.STRING, 1); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:Err: ", p), err) }
  if err := oprot.WriteString(ctx, string(p.Err)); err != nil {
  return thrift.PrependError(fmt.Sprintf("%T.Err (1) field write error: ", p), err) }
  if err := oprot.WriteFieldEnd(ctx); err != nil {
    return thrift.PrependError(fmt.Sprintf("%T write field end error 1:Err: ", p), err) }
  return err
}

func (p *RpcError) Equals(other *RpcError) bool {
  if p == other {
    return true
  } else if p == nil || other == nil {
    return false
  }
  if p.Err != other.Err { return false }
  return true
}

func (p *RpcError) String() string {
  if p == nil {
    return "<nil>"
  }
  return fmt.Sprintf("RpcError(%+v)", *p)
}

func (p *RpcError) Validate() error {
  return nil
}
