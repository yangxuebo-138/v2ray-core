package ratelimit

import (
	"fmt"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
)

type RateLimitStatWriter struct {
	Bucket *Bucket
	Writer  buf.Writer
	SID    uint32
	Direct string
}



type RateLimitStatReader struct {
	Bucket *Bucket
	Reader buf.Reader
	SID    uint32
	Direct string
}

func (r *RateLimitStatReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	mb, err := r.Reader.ReadMultiBuffer()
	//fmt.Printf("Ratelimit reader len %d\n", mb.Len())
	fmt.Printf("Session id %d. Reader direct %s bytes %d\n", r.SID, r.Direct, int64(mb.Len()))
	mb2 := make(buf.MultiBuffer, 0, len(mb))
	for i := range mb {

		discardOK := r.Bucket.Wait(int64(mb[i].Len()), r.SID, r.Direct)
		//r.Bucket.WaitMaxDuration(int64(mb.Len()), 30e9)
		if discardOK {
			//fmt.Printf("Discard bytes len %d.\n", mb.Len())
			//buf.ReleaseMulti(mb)
			//mb[i].Release()
			//fmt.Printf("Discard after bytes len %d.\n", mb.Len())
			//mb[i] = nil
			break
		} else {
			buffer := buf.New()
			buffer.Write(mb[i].Bytes())
			//b := make([]int, elementCount)
			mb2 = append(mb2, buffer)
		}
	}
	//
	buf.ReleaseMulti(mb)

	//discardOK := r.Bucket.Wait(int64(mb.Len()), r.SID, r.Direct)
	//if discardOK {
	//	fmt.Printf("Discard bytes len %d.\n", mb.Len())
	//	buf.ReleaseMulti(mb)
	//	fmt.Printf("Discard after bytes len %d.\n", mb.Len())
	//}

	//discardOK := r.Bucket.Wait(int64(mb.Len()), r.SID, r.Direct)
	////r.Bucket.WaitMaxDuration(int64(mb.Len()), 30e9)
	//if discardOK {
	//	fmt.Printf("Discard bytes len %d.\n", mb.Len())
	//	buf.ReleaseMulti(mb)
	//	fmt.Printf("Discard after bytes len %d.\n", mb.Len())
	//}
	return mb2, err
}

func (r *RateLimitStatReader) Close() error {
	fmt.Printf("---------------- Ratelimit Close SID : %d. Direct %s. capacity %d --------------------\n",r.SID, r.Direct, r.Bucket.capacity)
	return common.Close(r.Reader)
}

func (r *RateLimitStatReader) Interrupt() {
	fmt.Printf("---------------- Ratelimit Interrupt SID : %d. Direct %s. capacity %d --------------------\n", r.SID, r.Direct, r.Bucket.capacity)
	//r.Bucket.Signal(r.SID)
	common.Interrupt(r.Reader)
}

func (w *RateLimitStatWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	fmt.Printf("Session id %d. Write direct %s bytes %d\n", w.SID, w.Direct, int64(mb.Len()))
	w.Bucket.Wait(int64(mb.Len()),w.SID, w.Direct)

	return w.Writer.WriteMultiBuffer(mb)
}

func (w *RateLimitStatWriter) Close() error {
	fmt.Printf("---------------- Ratelimit Close SID : %d. Direct %s. capacity %d --------------------\n",w.SID, w.Direct, w.Bucket.capacity)
	return common.Close(w.Writer)
}

func (w *RateLimitStatWriter) Interrupt() {
	//w.Bucket.Signal(w.SID)
	fmt.Printf("---------------- Ratelimit Interrupt SID : %d. Direct %s. capacity %d --------------------\n", w.SID, w.Direct, w.Bucket.capacity)
	common.Interrupt(w.Writer)
}