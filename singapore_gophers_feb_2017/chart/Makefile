WIDTH ?= 1400
HEIGHT ?= 600
RESOLUTION ?= 100ms

all: build run

build:
	go build -o chart

INPUT ?= workload.one_time.open.0.idle.0.interval.100ms.duration.5s.mutiplier.1.0.ceil.1000
input := ../logs/$(INPUT)
output := $(input).svg
run.once:
	./chart --width=$(WIDTH) --height=$(HEIGHT) --input=$(input) --output=$(output).res.$(RESOLUTION).svg --resolution=$(RESOLUTION)

run:
	@for file in `ls ../logs/ | grep -v '.svg'`; do \
		./chart --width=$(WIDTH) --height=$(HEIGHT) --input=../logs/$$file --output=../logs/$$file.res.$(RESOLUTION).svg --resolution=$(RESOLUTION); \
  done

sample: build
	@for file in `ls ../logs/ | grep -v '.svg'`; do \
		./chart --width=$(WIDTH) --height=$(HEIGHT) --input=../logs/$$file --output=../logs/$$file.res.100ms.svg --resolution=100ms; \
		./chart --width=$(WIDTH) --height=$(HEIGHT) --input=../logs/$$file --output=../logs/$$file.res.500ms.svg --resolution=500ms; \
  done
