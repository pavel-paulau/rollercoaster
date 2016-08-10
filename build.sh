#!/bin/bash -ex

make clean

for os in "windows" "linux" "darwin"; do
	ext=""
	if [ ${os} = "windows" ]; then
		ext=".exe"
	fi

	for arch in "386" "amd64"; do
		GOOS=${os} GOARCH=${arch} CGO_ENABLED=0 go build -a --ldflags "-s -w" -o build/rollercoaster_${os}_${arch}${ext}
	done
done
