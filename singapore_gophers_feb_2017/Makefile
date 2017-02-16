OPEN ?= 0
IDLE ?= 0 WORKLOAD ?= sustained
INTERVAL ?= 50ms
DURATION ?= 10s
MULTIPLIER ?= 1.0
CEIL ?= 1000

all: build run

build:
	go build -race -v -o main

file := logs/workload.$(WORKLOAD).open.$(OPEN).idle.$(IDLE).interval.$(INTERVAL).duration.$(DURATION).mutiplier.$(MULTIPLIER).ceil.$(CEIL)
run:
	GOTRACEBACK=all ./main --max-open=$(OPEN) --max-idle=$(IDLE) --workload=$(WORKLOAD) --duration=$(DURATION) --interval=$(INTERVAL) --multiplier=$(MULTIPLIER) --ceil=$(CEIL) > $(file) 2>&1

FILE ?= $(file)
stats:
	grep Stats $(FILE)

sample: build
	rm logs/*
	for ceil in 1000 2000 3000 4000 5000; \
		do GOTRACEBACK=all ./main --max-open=0 --max-idle=0 --workload=one_time --ceil=$$ceil > logs/workload.one_time.open.0.idle.0.interval.50ms.duration.10s.mutiplier.1.0.ceil.$$ceil 2>&1; \
	done
	for mul in 0.25 0.5 0.75 1.0 1.25 1.5 1.75 2.0; \
		do GOTRACEBACK=all ./main --max-open=20 --max-idle=20 --workload=sustained --duration=30s --interval=50ms --multiplier=$$mul --ceil=1000 > logs/workload.sustained.open.20.idle.20.interval.50ms.duration.30s.mutiplier.$$mul.ceil.1000 2>&1; \
	done
