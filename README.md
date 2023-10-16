# Audigent House Coding Challange 

The requirement for the coding challange is below

# Cache
​
Develop a cache that supports TTL with an active strategy instead of a lazy strategy (or passive).
You can read more about lazy vs active removing [here](https://www.pankajtanwar.in/blog/how-redis-expires-keys-a-deep-dive-into-how-ttl-works-internally-in-redis).
​
The cache should support a string-like type as key and a byte slice as value.
You can't use a `map` or an external library. You must create your own data structure. 
The cache functions can't produce more than 1 allocation per operation.
The code must come with a benchmark (in order to check for the allocations).
​
The cache should satisfy the following interface:
```go
type Cache interface {
	// Set will store the key value pair with a given TTL.
	Set(key, value []byte, ttl time.Duration)
​
	// Get returns the value stored using `key`.
	//
	// If the key is not present value will be set to nil.
	Get(key []byte) (value []byte, ttl time.Duration)
}
```
