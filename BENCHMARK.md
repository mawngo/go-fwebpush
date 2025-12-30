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
BenchmarkDefaultConfig/run_0-24                            29893             39896 ns/op            7475 B/op         65 allocs/op
BenchmarkDefaultConfig/run_1-24                            30127             39840 ns/op            7475 B/op         65 allocs/op
BenchmarkDefaultConfig/run_2-24                            30004             40024 ns/op            7475 B/op         65 allocs/op
BenchmarkDefaultConfig/run_3-24                            29847             40035 ns/op            7475 B/op         65 allocs/op
BenchmarkDefaultConfig/run_4-24                            30126             39794 ns/op            7475 B/op         65 allocs/op
BenchmarkDefaultConfig/run_5-24                            30296             39845 ns/op            7475 B/op         65 allocs/op
BenchmarkDefaultConfig/run_6-24                            30060             40307 ns/op            7475 B/op         65 allocs/op
BenchmarkOldImpl/run_0-24                                  10000            110664 ns/op           20545 B/op        230 allocs/op
BenchmarkOldImpl/run_1-24                                  10000            110548 ns/op           20706 B/op        230 allocs/op
BenchmarkOldImpl/run_2-24                                  10000            110220 ns/op           20707 B/op        230 allocs/op
BenchmarkOldImpl/run_3-24                                  10000            110394 ns/op           20544 B/op        230 allocs/op
BenchmarkOldImpl/run_4-24                                  10000            109622 ns/op           20542 B/op        230 allocs/op
BenchmarkOldImpl/run_5-24                                  10000            110317 ns/op           20542 B/op        230 allocs/op
BenchmarkOldImpl/run_6-24                                  10000            110630 ns/op           20704 B/op        230 allocs/op
BenchmarkNoCaching/run_0-24                                10000            107543 ns/op           17208 B/op        173 allocs/op
BenchmarkNoCaching/run_1-24                                10000            106585 ns/op           17319 B/op        173 allocs/op
BenchmarkNoCaching/run_2-24                                10000            106763 ns/op           17320 B/op        173 allocs/op
BenchmarkNoCaching/run_3-24                                10000            106672 ns/op           17208 B/op        173 allocs/op
BenchmarkNoCaching/run_4-24                                10000            107399 ns/op           17208 B/op        173 allocs/op
BenchmarkNoCaching/run_5-24                                10000            108211 ns/op           17208 B/op        173 allocs/op
BenchmarkNoCaching/run_6-24                                10000            107830 ns/op           17320 B/op        173 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_0-24                29269             40673 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_1-24                29371             40738 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_2-24                29815             40519 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_3-24                29827             40424 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_4-24                29472             40719 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_5-24                29142             40908 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_6-24                29161             41269 ns/op            7811 B/op         70 allocs/op
BenchmarkVAPIDCaching/run_0-24                                     29470             40636 ns/op            7475 B/op         65 allocs/op
BenchmarkVAPIDCaching/run_1-24                                     29241             40968 ns/op            7475 B/op         65 allocs/op
BenchmarkVAPIDCaching/run_2-24                                     29673             40556 ns/op            7475 B/op         65 allocs/op
BenchmarkVAPIDCaching/run_3-24                                     29913             40288 ns/op            7475 B/op         65 allocs/op
BenchmarkVAPIDCaching/run_4-24                                     29792             40212 ns/op            7475 B/op         65 allocs/op
BenchmarkVAPIDCaching/run_5-24                                     29660             40300 ns/op            7475 B/op         65 allocs/op
BenchmarkVAPIDCaching/run_6-24                                     29752             40603 ns/op            7475 B/op         65 allocs/op
BenchmarkLocalSecretCaching/run_0-24                               10000            107737 ns/op           17544 B/op        178 allocs/op
BenchmarkLocalSecretCaching/run_1-24                               10000            107244 ns/op           17656 B/op        178 allocs/op
BenchmarkLocalSecretCaching/run_2-24                               10000            108032 ns/op           17656 B/op        178 allocs/op
BenchmarkLocalSecretCaching/run_3-24                               10000            107651 ns/op           17544 B/op        178 allocs/op
BenchmarkLocalSecretCaching/run_4-24                               10000            107351 ns/op           17544 B/op        178 allocs/op
BenchmarkLocalSecretCaching/run_5-24                               10000            107368 ns/op           17543 B/op        178 allocs/op
BenchmarkLocalSecretCaching/run_6-24                               10000            108006 ns/op           17657 B/op        178 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_0-24                       29875             40241 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_1-24                       29414             40606 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_2-24                       29328             40756 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_3-24                       29502             40547 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_4-24                       29786             40514 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_5-24                       29173             40764 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_6-24                       29350             40584 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_0-24              29766             40365 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_1-24              29457             40627 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_2-24              29487             40712 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_3-24              29850             40245 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_4-24              29857             40469 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_5-24              29559             40792 ns/op            7811 B/op         70 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_6-24              29384             40930 ns/op            7811 B/op         70 allocs/op
PASS
ok      github.com/mawngo/go-fwebpush   65.106s
```

# Conclusion

In the worst case scenario we achieve the same output compared to (sightly
optimized) [old implementation](https://github.com/SherClockHolmes/webpush-go) with lower allocations.

In the best case we achieve 15x performance, with the default config (only vapid cache enabled), we achieve 2.5x performance.