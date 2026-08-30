# Benchmark

## Run

```shell
go test -bench='.' -benchtime=1s -benchmem
```

## Result

```
❯ go test -bench='.' -benchtime=1s -benchmem
goos: windows
goarch: amd64
pkg: gitlab.com/moneyoyo/go-webpush/v2
cpu: AMD Ryzen 9 7900 12-Core Processor
BenchmarkDefaultConfig/run_0-24                            29349             40833 ns/op            6714 B/op         63 allocs/op
BenchmarkDefaultConfig/run_1-24                            28999             41378 ns/op            6714 B/op         63 allocs/op
BenchmarkDefaultConfig/run_2-24                            29125             41415 ns/op            6713 B/op         63 allocs/op
BenchmarkDefaultConfig/run_3-24                            29073             40928 ns/op            6714 B/op         63 allocs/op
BenchmarkDefaultConfig/run_4-24                            29316             40909 ns/op            6713 B/op         63 allocs/op
BenchmarkDefaultConfig/run_5-24                            29437             40838 ns/op            6713 B/op         63 allocs/op
BenchmarkDefaultConfig/run_6-24                            29335             40929 ns/op            6713 B/op         63 allocs/op
BenchmarkOldImpl/run_0-24                                  14641             82076 ns/op           19431 B/op        209 allocs/op
BenchmarkOldImpl/run_1-24                                  14419             83261 ns/op           19561 B/op        209 allocs/op
BenchmarkOldImpl/run_2-24                                  14389             83540 ns/op           19561 B/op        209 allocs/op
BenchmarkOldImpl/run_3-24                                  14410             82976 ns/op           19431 B/op        209 allocs/op
BenchmarkOldImpl/run_4-24                                  14536             82653 ns/op           19431 B/op        209 allocs/op
BenchmarkOldImpl/run_5-24                                  14359             83535 ns/op           19432 B/op        209 allocs/op
BenchmarkOldImpl/run_6-24                                  14433             83348 ns/op           19559 B/op        209 allocs/op
BenchmarkNoCaching/run_0-24                                14647             82048 ns/op           16460 B/op        170 allocs/op
BenchmarkNoCaching/run_1-24                                14610             82364 ns/op           16573 B/op        170 allocs/op
BenchmarkNoCaching/run_2-24                                14493             82702 ns/op           16572 B/op        170 allocs/op
BenchmarkNoCaching/run_3-24                                14697             81788 ns/op           16460 B/op        170 allocs/op
BenchmarkNoCaching/run_4-24                                14572             82244 ns/op           16460 B/op        170 allocs/op
BenchmarkNoCaching/run_5-24                                14638             82150 ns/op           16460 B/op        170 allocs/op
BenchmarkNoCaching/run_6-24                                14577             82454 ns/op           16573 B/op        170 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_0-24                29174             41110 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_1-24                28911             41293 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_2-24                29006             41175 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_3-24                29329             40898 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_4-24                29025             40976 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_5-24                29020             41152 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingExpired/run_6-24                28898             41188 ns/op            7066 B/op         68 allocs/op
BenchmarkVAPIDCaching/run_0-24                                     29336             40726 ns/op            6714 B/op         63 allocs/op
BenchmarkVAPIDCaching/run_1-24                                     29185             40949 ns/op            6714 B/op         63 allocs/op
BenchmarkVAPIDCaching/run_2-24                                     29306             40989 ns/op            6714 B/op         63 allocs/op
BenchmarkVAPIDCaching/run_3-24                                     29253             40965 ns/op            6714 B/op         63 allocs/op
BenchmarkVAPIDCaching/run_4-24                                     29292             40903 ns/op            6714 B/op         63 allocs/op
BenchmarkVAPIDCaching/run_5-24                                     29241             40788 ns/op            6713 B/op         63 allocs/op
BenchmarkVAPIDCaching/run_6-24                                     29566             40972 ns/op            6714 B/op         63 allocs/op
BenchmarkLocalSecretCaching/run_0-24                               25490             47524 ns/op           14827 B/op        149 allocs/op
BenchmarkLocalSecretCaching/run_1-24                               25742             46649 ns/op           14940 B/op        149 allocs/op
BenchmarkLocalSecretCaching/run_2-24                               24446             47109 ns/op           14939 B/op        149 allocs/op
BenchmarkLocalSecretCaching/run_3-24                               25432             47633 ns/op           14827 B/op        149 allocs/op
BenchmarkLocalSecretCaching/run_4-24                               26053             47143 ns/op           14826 B/op        149 allocs/op
BenchmarkLocalSecretCaching/run_5-24                               26092             46688 ns/op           14827 B/op        149 allocs/op
BenchmarkLocalSecretCaching/run_6-24                               25494             46769 ns/op           14939 B/op        149 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_0-24                      200473              5801 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_1-24                      262711              6764 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_2-24                      192058              6480 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_3-24                      330228              6765 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_4-24                      186922              5845 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_5-24                      251011              5566 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCaching/run_6-24                      276121              6683 ns/op            5079 B/op         42 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_0-24              29359             41066 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_1-24              28866             41411 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_2-24              29073             41216 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_3-24              29533             41045 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_4-24              29154             41149 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_5-24              29592             40796 ns/op            7066 B/op         68 allocs/op
BenchmarkVapidAndLocalSecretCachingCacheInit/run_6-24              29121             41298 ns/op            7066 B/op         68 allocs/op
BenchmarkGetCachedKey/run_0-24                                     28508             41882 ns/op            9735 B/op        107 allocs/op
BenchmarkGetCachedKey/run_1-24                                     28410             41933 ns/op            9847 B/op        107 allocs/op
BenchmarkGetCachedKey/run_2-24                                     28726             41660 ns/op            9847 B/op        107 allocs/op
BenchmarkGetCachedKey/run_3-24                                     28755             41767 ns/op            9736 B/op        107 allocs/op
BenchmarkGetCachedKey/run_4-24                                     28627             41732 ns/op            9735 B/op        107 allocs/op
BenchmarkGetCachedKey/run_5-24                                     28774             41654 ns/op            9735 B/op        107 allocs/op
BenchmarkGetCachedKey/run_6-24                                     28933             41737 ns/op            9847 B/op        107 allocs/op
PASS
ok      gitlab.com/moneyoyo/go-webpush/v2       78.406s
```

# Conclusion

In the worst case scenario we achieve the same output compared to (sightly
optimized) [old implementation](https://github.com/SherClockHolmes/webpush-go) with lower allocations.

In the best case we achieve 14x performance; with the default config (only vapid cache enabled), we achieve 2x performance.