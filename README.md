# About
This is an excersie toy-load-balancer project.

# Running
- `go run main.go`

# Example program output
```
-----------------------------------------
Running balancer with strategy round-robin...
Register provider simple1 
Register provider simple2 
Register provider simple3 
Register provider faulty1_permanent 
Health check on faulty1_permanent result false
Register provider faulty2_temporary 
Health check on faulty2_temporary result false
Exclude provider simple1 
Get() = simple2 
Get() = simple3 
Get() = simple2 
Get() = simple3 
Get() = simple2 
Include provider simple1 
Get() = simple2 
Get() = simple3 
Get() = simple1 
Get() = simple2 
Get() = simple3 
Make sure some health checks have passed...
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Get() = faulty2_temporary 
Get() = simple1 
Get() = simple2 
Get() = simple3 
Get() = faulty2_temporary 
-----------------------------------------
Running balancer with strategy random...
Register provider simple1 
Register provider simple2 
Register provider simple3 
Register provider faulty1_permanent 
Health check on faulty1_permanent result false
Register provider faulty2_temporary 
Health check on faulty2_temporary result false
Exclude provider simple1 
Get() = simple3 
Get() = simple3 
Get() = simple3 
Get() = simple3 
Get() = simple3 
Include provider simple1 
Get() = simple1 
Get() = simple2 
Get() = simple3 
Get() = simple2 
Get() = simple1 
Make sure some health checks have passed...
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Health check on faulty1_permanent result false
Health check on faulty2_temporary result true
Get() = simple3 
Get() = faulty2_temporary 
Get() = simple3 
Get() = simple2 
Get() = simple1 
-----------------------------------------
Running empty balancer...
Get() error = too many requests in flight 
Get() error = too many requests in flight 
Get() error = too many requests in flight 
Get() error = too many requests in flight 
Get() error = too many requests in flight 
```