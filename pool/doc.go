/*
Package pool provides allocation pools for signal buffers.

To decrease a number of allocations at runtime, pools can be used.
Internally it relies on sync.Pool to manage objects in memory. The caller
is responsible for ensuring that buffer has same capacity before being
returned.
*/
package pool
