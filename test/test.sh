#!/bin/bash

docker run -v $(pwd):/d docker.io/airdock/lilypond:latest -dresolution=300 --png -dbackend=eps -dno-gs-load-fonts -dinclude-eps-fonts  -o /d/out /d/test.ly