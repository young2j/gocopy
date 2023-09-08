<h1 align="center">GoCopy</h1>

<p align="center">
  <img alt="Github top language" src="https://img.shields.io/github/languages/top/young2j/gocopy?color=56BEB8">
  <img alt="Github language count" src="https://img.shields.io/github/languages/count/young2j/gocopy?color=56BEB8">
  <img alt="Repository size" src="https://img.shields.io/github/repo-size/young2j/gocopy?color=56BEB8">
  <img alt="License" src="https://img.shields.io/github/license/young2j/gocopy?color=56BEB8">
<img alt="Github forks" src="https://img.shields.io/github/forks/young2j/gocopy?color=56BEB8" />
  <img alt="Github stars" src="https://img.shields.io/github/stars/young2j/gocopy?color=56BEB8" />
</p>


<p align="center">
  <a href="#dart-about">About</a> &#xa0; | &#xa0; 
  <a href="#sparkles-features">Features</a> &#xa0; | &#xa0;
  <a href="#white_check_mark-requirements">Requirements</a> &#xa0; | &#xa0;
  <a href="#hammer_and_wrench-installation">Installation</a> &#xa0; | &#xa0;
  <a href="#electric_plug-api">API</a> &#xa0; | &#xa0;
  <a href="#fire-benchmark">Benchmark</a> &#xa0;
</p>

<br>

>  :rage::rage::rage:Dumping nuclear wastewater into the ocean, damn it! :bomb::japan::boom::triumph::triumph::triumph:



## :dart: About ##

GoCopy is a simple golang repo for copying slice to slice, map to map, struct to struct, struct to map or bson.M.



## :sparkles: Features ##

​	:heavy_check_mark: copy slice to slice by type<br/>

​	:heavy_check_mark: copy map to map by type<br/>

​	:heavy_check_mark: copy struct to struct by field name<br/>

​	:heavy_check_mark: copy struct to map/bson.M by field name<br/>

​	:heavy_check_mark: support append values, change filed-case, rename or ignore any field when copying<br/>

​	:heavy_check_mark: support customized convert function to transform any filed to what you want<br/>



## :white_check_mark: Requirements ##

​	golang >=1.16



## :hammer_and_wrench: Installation ##

```shell
go get -u github.com/young2j/gocopy@latest
```



## :electric_plug: API

* `Copy(to, from interface{})`
* `CopyWithOption(to,from interface{},opt *Option)`

> Note: The arg `to` must be a reference value(usually `&to`), or copy maybe fail.

see more  at [`./example`](./example)

```shell
# path/to/example
$ go run .
```



## :beers: Related repository

* [`copier`](https://github.com/jinzhu/copier)



## :fire: Benchmark

```shell
go test -v . -bench=.  -benchmem -benchtime=1s -cpu=4
```

```shell
goos: darwin
goarch: amd64
pkg: github.com/young2j/gocopy
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkCopy
BenchmarkCopy-4     	  122139	      8884 ns/op	    5592 B/op	      81 allocs/op
BenchmarkCopier
BenchmarkCopier-4   	   62940	     18695 ns/op	   14640 B/op	     166 allocs/op
PASS
ok  	github.com/young2j/gocopy	4.999s
```



## :memo: License ##

This project is under license from MIT. For more details, see the [LICENSE](LICENSE.md) file.


Made with :heart: by <a href="https://github.com/young2j" target="_blank">young2j</a>
