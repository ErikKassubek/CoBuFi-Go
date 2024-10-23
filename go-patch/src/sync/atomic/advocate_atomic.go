package atomic

import "runtime"

// SwapInt32 atomically stores new into *addr and returns the previous *addr value.
// Consider using the more ergonomic and less error-prone [Int32.Swap] instead.
func SwapInt32(addr *int32, new int32) (old int32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.SwapOp)
	return SwapInt32Advocate(addr, new)
}

// SwapInt64 atomically stores new into *addr and returns the previous *addr value.
// Consider using the more ergonomic and less error-prone [Int64.Swap] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func SwapInt64(addr *int64, new int64) (old int64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.SwapOp)
	return SwapInt64Advocate(addr, new)
}

// SwapUint32 atomically stores new into *addr and returns the previous *addr value.
// Consider using the more ergonomic and less error-prone [Uint32.Swap] instead.
func SwapUint32(addr *uint32, new uint32) (old uint32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.SwapOp)
	return SwapUint32Advocate(addr, new)
}

// SwapUint64 atomically stores new into *addr and returns the previous *addr value.
// Consider using the more ergonomic and less error-prone [Uint64.Swap] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func SwapUint64(addr *uint64, new uint64) (old uint64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.SwapOp)
	return SwapUint64Advocate(addr, new)
}

// SwapUintptr atomically stores new into *addr and returns the previous *addr value.
// Consider using the more ergonomic and less error-prone [Uintptr.Swap] instead.
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.SwapOp)
	return SwapUintptrAdvocate(addr, new)
}

// SwapPointer atomically stores new into *addr and returns the previous *addr value.
// Consider using the more ergonomic and less error-prone [Pointer.Swap] instead.
// func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer) {
// 	return SwapPointerAdvocate(addr, new)
// }

// CompareAndSwapInt32 executes the compare-and-swap operation for an int32 value.
// Consider using the more ergonomic and less error-prone [Int32.CompareAndSwap] instead.
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicCompareAndSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.CompSwapOp)
	return CompareAndSwapInt32Advocate(addr, old, new)
}

// CompareAndSwapInt64 executes the compare-and-swap operation for an int64 value.
// Consider using the more ergonomic and less error-prone [Int64.CompareAndSwap] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicCompareAndSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.CompSwapOp)
	return CompareAndSwapInt64Advocate(addr, old, new)
}

// CompareAndSwapUint32 executes the compare-and-swap operation for a uint32 value.
// Consider using the more ergonomic and less error-prone [Uint32.CompareAndSwap] instead.
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicCompareAndSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.CompSwapOp)
	return CompareAndSwapUint32Advocate(addr, old, new)
}

// CompareAndSwapUint64 executes the compare-and-swap operation for a uint64 value.
// Consider using the more ergonomic and less error-prone [Uint64.CompareAndSwap] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicCompareAndSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.CompSwapOp)
	return CompareAndSwapUint64Advocate(addr, old, new)
}

// CompareAndSwapUintptr executes the compare-and-swap operation for a uintptr value.
// Consider using the more ergonomic and less error-prone [Uintptr.CompareAndSwap] instead.
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicCompareAndSwap, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.CompSwapOp)
	return CompareAndSwapUintptrAdvocate(addr, old, new)
}

// // CompareAndSwapPointer executes the compare-and-swap operation for a unsafe.Pointer value.
// // Consider using the more ergonomic and less error-prone [Pointer.CompareAndSwap] instead.
// func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool) {
// 	return CompareAndSwapPointerAdvocate(addr, old, new)
// }

// AddInt32 atomically adds delta to *addr and returns the new value.
// Consider using the more ergonomic and less error-prone [Int32.Add] instead.
func AddInt32(addr *int32, delta int32) (new int32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicAdd, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.AddOp)
	return AddInt32Advocate(addr, delta)
}

// AddUint32 atomically adds delta to *addr and returns the new value.
// To subtract a signed positive constant value c from x, do AddUint32(&x, ^uint32(c-1)).
// In particular, to decrement x, do AddUint32(&x, ^uint32(0)).
// Consider using the more ergonomic and less error-prone [Uint32.Add] instead.
func AddUint32(addr *uint32, delta uint32) (new uint32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicAdd, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.AddOp)
	return AddUint32Advocate(addr, delta)
}

// AddInt64 atomically adds delta to *addr and returns the new value.
// Consider using the more ergonomic and less error-prone [Int64.Add] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func AddInt64(addr *int64, delta int64) (new int64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicAdd, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.AddOp)
	return AddInt64Advocate(addr, delta)
}

// AddUint64 atomically adds delta to *addr and returns the new value.
// To subtract a signed positive constant value c from x, do AddUint64(&x, ^uint64(c-1)).
// In particular, to decrement x, do AddUint64(&x, ^uint64(0)).
// Consider using the more ergonomic and less error-prone [Uint64.Add] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func AddUint64(addr *uint64, delta uint64) (new uint64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicAdd, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.AddOp)
	return AddUint64Advocate(addr, delta)
}

// AddUintptr atomically adds delta to *addr and returns the new value.
// Consider using the more ergonomic and less error-prone [Uintptr.Add] instead.
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicAdd, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.AddOp)
	return AddUintptrAdvocate(addr, delta)
}

// LoadInt32 atomically loads *addr.
// Consider using the more ergonomic and less error-prone [Int32.Load] instead.
func LoadInt32(addr *int32) (val int32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicLoad, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.LoadOp)
	return LoadInt32Advocate(addr)
}

// LoadInt64 atomically loads *addr.
// Consider using the more ergonomic and less error-prone [Int64.Load] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func LoadInt64(addr *int64) (val int64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicLoad, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.LoadOp)
	return LoadInt64Advocate(addr)
}

// LoadUint32 atomically loads *addr.
// Consider using the more ergonomic and less error-prone [Uint32.Load] instead.
func LoadUint32(addr *uint32) (val uint32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicLoad, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.LoadOp)
	return LoadUint32Advocate(addr)
}

// LoadUint64 atomically loads *addr.
// Consider using the more ergonomic and less error-prone [Uint64.Load] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func LoadUint64(addr *uint64) (val uint64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicLoad, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.LoadOp)
	return LoadUint64Advocate(addr)
}

// LoadUintptr atomically loads *addr.
// Consider using the more ergonomic and less error-prone [Uintptr.Load] instead.
func LoadUintptr(addr *uintptr) (val uintptr) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicLoad, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.LoadOp)
	return LoadUintptrAdvocate(addr)
}

// // LoadPointer atomically loads *addr.
// // Consider using the more ergonomic and less error-prone [Pointer.Load] instead.
// func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer) {
// 	return LoadPointerAdvocate(addr)
// }

// StoreInt32 atomically stores val into *addr.
// Consider using the more ergonomic and less error-prone [Int32.Store] instead.
func StoreInt32(addr *int32, val int32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicStore, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.StoreOp)
	StoreInt32Advocate(addr, val)
}

// StoreInt64 atomically stores val into *addr.
// Consider using the more ergonomic and less error-prone [Int64.Store] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func StoreInt64(addr *int64, val int64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicStore, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.StoreOp)
	StoreInt64Advocate(addr, val)
}

// StoreUint32 atomically stores val into *addr.
// Consider using the more ergonomic and less error-prone [Uint32.Store] instead.
func StoreUint32(addr *uint32, val uint32) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicStore, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.StoreOp)
	StoreUint32Advocate(addr, val)
}

// StoreUint64 atomically stores val into *addr.
// Consider using the more ergonomic and less error-prone [Uint64.Store] instead
// (particularly if you target 32-bit platforms; see the bugs section).
func StoreUint64(addr *uint64, val uint64) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicStore, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.StoreOp)
	StoreUint64Advocate(addr, val)
}

// StoreUintptr atomically stores val into *addr.
// Consider using the more ergonomic and less error-prone [Uintptr.Store] instead.
func StoreUintptr(addr *uintptr, val uintptr) {
	wait, ch := runtime.WaitForReplay(runtime.OperationAtomicStore, 2)
	if wait {
		<- ch
	}
	runtime.AdvocateAtomic(addr, runtime.StoreOp)
	StoreUintptrAdvocate(addr, val)
}

// // StorePointer atomically stores val into *addr.
// // Consider using the more ergonomic and less error-prone [Pointer.Store] instead.
// func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer) {
// 	StorePointerAdvocate(addr, val)
// }
