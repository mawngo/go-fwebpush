# Benchmark

## Run

```shell
go test -bench=. -benchtime=1s -benchmem
```

## Result

```
❯ go test -bench="."  -benchtime=1s -benchmem
goos: windows
goarch: amd64
pkg: github.com/mawngo/go-fwebpush
cpu: AMD Ryzen 9 7900 12-Core Processor             
BenchmarkDefaultConfig/run_0-24                            30730             39019 ns/op            6963 B/op         64 allocs/op
BenchmarkDefaultConfig/run_1-24                            30111             39427 ns/op            6963 B/op         64 allocs/op
BenchmarkDefaultConfig/run_2-24                            30710             39206 ns/op            6963 B/op         64 allocs/op
BenchmarkDefaultConfig/run_3-24                            30648             39128 ns/op            6963 B/op         64 allocs/op
BenchmarkDefaultConfig/run_4-24                            30638             39122 ns/op            6963 B/op         64 allocs/op
BenchmarkDefaultConfig/run_5-24                            30720             39018 ns/op            6963 B/op         64 allocs/op
BenchmarkDefaultConfig/run_6-24                            30506             39268 ns/op            6963 B/op         64 allocs/op
BenchmarkOldImpl/run_0-24                                  10000            108521 ns/op           20542 B/op        230 allocs/op
BenchmarkOldImpl/run_1-24                                  10000            108801 ns/op           20702 B/op        230 allocs/op
BenchmarkOldImpl/run_2-24                                  10000            108669 ns/op           20704 B/op        230 allocs/op
BenchmarkOldImpl/run_3-24                                  10000            108582 ns/op           20541 B/op        230 allocs/op
BenchmarkOldImpl/run_4-24                                  10000            108126 ns/op           20541 B/op        230 allocs/op
BenchmarkOldImpl/run_5-24                                  10000            108535 ns/op           20541 B/op        230 allocs/op
BenchmarkOldImpl/run_6-24                                  10000            109035 ns/op           20703 B/op        230 allocs/op
BenchmarkNoCaching/run_0-24                                10000            105090 ns/op           16600 B/op        171 allocs/op
BenchmarkNoCaching/run_1-24                                10000            105518 ns/op           16711 B/op        171 allocs/op
BenchmarkNoCaching/run_2-24                                10000            105123 ns/op           16711 B/op        171 allocs/op
BenchmarkNoCaching/run_3-24                                10000            105673 ns/op           16599 B/op        171 allocs/op
BenchmarkNoCaching/run_4-24                                10000            105368 ns/op           16600 B/op        171 allocs/op
BenchmarkNoCaching/run_5-24                                10000            105895 ns/op           16601 B/op        171 allocs/op
BenchmarkNoCaching/run_6-24                                10000            105697 ns/op           16712 B/op        171 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_0-24                30404             39417 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_1-24                30195             39691 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_2-24                30157             39656 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_3-24                30130             39491 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_4-24                30573             39315 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_5-24                30518             39412 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_6-24                30093             39730 ns/op            7315 B/op         69 allocs/op
BenchmarkVAPIDCaching/run_0-24                                     30669             39098 ns/op            6963 B/op         64 allocs/op
BenchmarkVAPIDCaching/run_1-24                                     30346             39537 ns/op            6963 B/op         64 allocs/op
BenchmarkVAPIDCaching/run_2-24                                     30507             39349 ns/op            6963 B/op         64 allocs/op
BenchmarkVAPIDCaching/run_3-24                                     30318             39580 ns/op            6963 B/op         64 allocs/op
BenchmarkVAPIDCaching/run_4-24                                     30493             40068 ns/op            6963 B/op         64 allocs/op
BenchmarkVAPIDCaching/run_5-24                                     29517             40262 ns/op            6963 B/op         64 allocs/op
BenchmarkVAPIDCaching/run_6-24                                     29463             40136 ns/op            6963 B/op         64 allocs/op
BenchmarkLocalSecretCaching/run_0-24                               16866             71510 ns/op           14968 B/op        150 allocs/op
BenchmarkLocalSecretCaching/run_1-24                               17071             70307 ns/op           15079 B/op        150 allocs/op
BenchmarkLocalSecretCaching/run_2-24                               16840             71255 ns/op           15080 B/op        150 allocs/op
BenchmarkLocalSecretCaching/run_3-24                               17030             70465 ns/op           14967 B/op        150 allocs/op
BenchmarkLocalSecretCaching/run_4-24                               16965             70366 ns/op           14967 B/op        150 allocs/op
BenchmarkLocalSecretCaching/run_5-24                               17185             70005 ns/op           14968 B/op        150 allocs/op
BenchmarkLocalSecretCaching/run_6-24                               16914             70926 ns/op           15079 B/op        150 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_0-24                      347773              3671 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_1-24                      312048              3808 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_2-24                      330609              3651 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_3-24                      332665              3567 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_4-24                      329312              3748 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_5-24                      354231              3496 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_6-24                      319844              3771 ns/op            5330 B/op         43 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_0-24              29859             39815 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_1-24              30124             39783 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_2-24              30004             40083 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_3-24              30250             39569 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_4-24              30339             39595 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_5-24              30255             39718 ns/op            7315 B/op         69 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_6-24              30212             39778 ns/op            7315 B/op         69 allocs/op
BenchmarkGetCachedKey/run_0-24                                     18091             66345 ns/op            9637 B/op        107 allocs/op
BenchmarkGetCachedKey/run_1-24                                     18126             66244 ns/op            9749 B/op        107 allocs/op
BenchmarkGetCachedKey/run_2-24                                     18048             66557 ns/op            9749 B/op        107 allocs/op
BenchmarkGetCachedKey/run_3-24                                     18006             66378 ns/op            9639 B/op        107 allocs/op
BenchmarkGetCachedKey/run_4-24                                     18190             66357 ns/op            9638 B/op        107 allocs/op
BenchmarkGetCachedKey/run_5-24                                     18008             66895 ns/op            9636 B/op        107 allocs/op
BenchmarkGetCachedKey/run_6-24                                     18122             66057 ns/op            9749 B/op        107 allocs/op
PASS
ok      github.com/mawngo/go-fwebpush   74.227s
```

# Conclusion

In the worst case scenario we achieve the same output compared to (sightly
optimized) [old implementation](https://github.com/SherClockHolmes/webpush-go) with lower allocations.

In the best case we achieve 15x performance, with the default config (only vapid cache enabled), we achieve 2.5x performance.