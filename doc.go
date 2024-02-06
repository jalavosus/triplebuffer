// Package triplebuffer exposes a potentially viable triple-buffer system
// for managing data between a producer and consumer where avoiding read/write
// deadlocks might be problematic.
// Don't trust anything in this package for Production (prod) usage, oh my god.
package triplebuffer
