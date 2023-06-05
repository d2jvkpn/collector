#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


# go test -bench BenchmarkXXXX -outputdir=result -v -memprofile=mem.out -cpuprofile=cpu.out -blockprofile=block.out

#### 01
echo -e "\n>>> BenchmarkUnmarshal_01"
go test -bench BenchmarkUnmarshal_01 -run none -count 5 -parallel 4

echo -e "\n>>> BenchmarkMarshal_01"
go test -bench BenchmarkMarshal_01 -run none -count 5 -parallel 4

#### 02
echo -e "\n>>> BenchmarkUnmarshal_02"
go test -bench BenchmarkUnmarshal_02 -run none -count 5 -parallel 4

echo -e "\n>>> BenchmarkMarshal_02"
go test -bench BenchmarkMarshal_02 -run none -count 5 -parallel 4
