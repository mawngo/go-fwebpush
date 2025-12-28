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
BenchmarkDefaultConfig/run_0-24                            30466             39365 ns/op            7495 B/op         67 allocs/op
BenchmarkDefaultConfig/run_1-24                            30643             39268 ns/op            7495 B/op         67 allocs/op
BenchmarkDefaultConfig/run_2-24                            30549             39351 ns/op            7494 B/op         67 allocs/op
BenchmarkDefaultConfig/run_3-24                            30668             38954 ns/op            7494 B/op         67 allocs/op
BenchmarkDefaultConfig/run_4-24                            30858             39052 ns/op            7494 B/op         67 allocs/op
BenchmarkDefaultConfig/run_5-24                            30526             39088 ns/op            7494 B/op         67 allocs/op
BenchmarkDefaultConfig/run_6-24                            30423             39328 ns/op            7494 B/op         67 allocs/op
BenchmarkOldImpl/run_0-24                                  10000            107726 ns/op           20543 B/op        230 allocs/op
BenchmarkOldImpl/run_1-24                                  10000            109092 ns/op           20703 B/op        230 allocs/op
BenchmarkOldImpl/run_2-24                                  10000            109221 ns/op           20704 B/op        230 allocs/op
BenchmarkOldImpl/run_3-24                                  10000            108164 ns/op           20542 B/op        230 allocs/op
BenchmarkOldImpl/run_4-24                                  10000            108128 ns/op           20542 B/op        230 allocs/op
BenchmarkOldImpl/run_5-24                                  10000            108264 ns/op           20541 B/op        230 allocs/op
BenchmarkOldImpl/run_6-24                                  10000            108724 ns/op           20704 B/op        230 allocs/op
BenchmarkNoCaching/run_0-24                                10000            105555 ns/op           17223 B/op        175 allocs/op
BenchmarkNoCaching/run_1-24                                10000            106378 ns/op           17336 B/op        175 allocs/op
BenchmarkNoCaching/run_2-24                                10000            105294 ns/op           17335 B/op        175 allocs/op
BenchmarkNoCaching/run_3-24                                10000            105797 ns/op           17222 B/op        175 allocs/op
BenchmarkNoCaching/run_4-24                                10000            104870 ns/op           17223 B/op        175 allocs/op
BenchmarkNoCaching/run_5-24                                10000            104960 ns/op           17223 B/op        175 allocs/op
BenchmarkNoCaching/run_6-24                                10000            115947 ns/op           17336 B/op        175 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_0-24                28590             42195 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_1-24                27214             43842 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_2-24                28329             42389 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_3-24                25562             44630 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_4-24                28231             42170 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_5-24                28560             42038 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_6-24                28592             42434 ns/op            7830 B/op         72 allocs/op
BenchmarkVAPIDCaching/run_0-24                                     28994             41514 ns/op            7494 B/op         67 allocs/op
BenchmarkVAPIDCaching/run_1-24                                     28790             41882 ns/op            7494 B/op         67 allocs/op
BenchmarkVAPIDCaching/run_2-24                                     28302             42108 ns/op            7494 B/op         67 allocs/op
BenchmarkVAPIDCaching/run_3-24                                     28977             41433 ns/op            7494 B/op         67 allocs/op
BenchmarkVAPIDCaching/run_4-24                                     28764             41698 ns/op            7494 B/op         67 allocs/op
BenchmarkVAPIDCaching/run_5-24                                     30006             39756 ns/op            7494 B/op         67 allocs/op
BenchmarkVAPIDCaching/run_6-24                                     29965             40003 ns/op            7494 B/op         67 allocs/op
BenchmarkLocalSecretCaching/run_0-24                               16873             70798 ns/op           16840 B/op        170 allocs/op
BenchmarkLocalSecretCaching/run_1-24                               16833             71440 ns/op           16951 B/op        170 allocs/op
BenchmarkLocalSecretCaching/run_2-24                               15850             75886 ns/op           16951 B/op        170 allocs/op
BenchmarkLocalSecretCaching/run_3-24                               15735             76013 ns/op           16840 B/op        170 allocs/op
BenchmarkLocalSecretCaching/run_4-24                               15811             75665 ns/op           16839 B/op        170 allocs/op
BenchmarkLocalSecretCaching/run_5-24                               15744             75997 ns/op           16840 B/op        170 allocs/op
BenchmarkLocalSecretCaching/run_6-24                               15681             76546 ns/op           16952 B/op        170 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_0-24                      267289              4641 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_1-24                      238840              4898 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_2-24                      257683              4747 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_3-24                      256612              4644 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_4-24                      271308              4880 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_5-24                      272463              4806 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_6-24                      203974              5124 ns/op            7110 B/op         62 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_0-24              28374             42279 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_1-24              26832             43420 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_2-24              27472             43420 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_3-24              28917             41633 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_4-24              28735             41800 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_5-24              28788             41803 ns/op            7830 B/op         72 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_6-24              28747             41962 ns/op            7830 B/op         72 allocs/op
PASS
ok      github.com/mawngo/go-fwebpush  65.818s
```

# Conclusion

In the worst case scenario we achieve the same output compared to (sightly
optimized) [old implementation](https://github.com/SherClockHolmes/webpush-go) with lower allocations.

In the best case we achieve 15x performance, with the default config (only vapid cache enabled), we achieve 2.5x performance.