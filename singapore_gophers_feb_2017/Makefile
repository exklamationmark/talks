OPEN ?= 0
IDLE ?= 0
WORKLOAD ?= one_time
INTERVAL ?= 100ms
DURATION ?= 5s
MULTIPLIER ?= 1.0
CEIL ?= 1000

all: build run

build:
	go build -race -v -o main

file := logs/workload.$(WORKLOAD).open.$(OPEN).idle.$(IDLE).interval.$(INTERVAL).duration.$(DURATION).mutiplier.$(MULTIPLIER).ceil.$(CEIL)
run:
	GOTRACEBACK=all ./main --max-open=$(OPEN) --max-idle=$(IDLE) --workload=$(WORKLOAD) --duration=$(DURATION) --interval=$(INTERVAL) --multiplier=$(MULTIPLIER) --ceil=$(CEIL) > $(file) 2>&1
	grep Stats $(file)
	exit 0
