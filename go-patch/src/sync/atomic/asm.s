// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

// ADVOCATE-CHANGE-START

#include "textflag.h"

TEXT ·SwapInt32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xchg(SB)

TEXT ·SwapUint32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xchg(SB)

TEXT ·SwapInt64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xchg64(SB)

TEXT ·SwapUint64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xchg64(SB)

TEXT ·SwapUintptrAdvocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xchguintptr(SB)

TEXT ·CompareAndSwapInt32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Cas(SB)

TEXT ·CompareAndSwapUint32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Cas(SB)

TEXT ·CompareAndSwapUintptrAdvocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Casuintptr(SB)

TEXT ·CompareAndSwapInt64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Cas64(SB)

TEXT ·CompareAndSwapUint64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Cas64(SB)

TEXT ·AddInt32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xadd(SB)

TEXT ·AddUint32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xadd(SB)

TEXT ·AddUintptrAdvocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xadduintptr(SB)

TEXT ·AddInt64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xadd64(SB)

TEXT ·AddUint64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Xadd64(SB)

TEXT ·LoadInt32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Load(SB)

TEXT ·LoadUint32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Load(SB)

TEXT ·LoadInt64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Load64(SB)

TEXT ·LoadUint64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Load64(SB)

TEXT ·LoadUintptrAdvocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Loaduintptr(SB)

TEXT ·LoadPointer(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Loadp(SB)

TEXT ·StoreInt32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Store(SB)

TEXT ·StoreUint32Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Store(SB)

TEXT ·StoreInt64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Store64(SB)

TEXT ·StoreUint64Advocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Store64(SB)

TEXT ·StoreUintptrAdvocate(SB),NOSPLIT,$0
	JMP	runtime∕internal∕atomic·Storeuintptr(SB)

// ADVOCATE-CHANGE-END
