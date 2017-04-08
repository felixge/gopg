PGDATA = $(shell "pwd")/data
INITDB = initdb
PG_CTL = pg_ctl
PG_LOG = $(PGDATA)/log.log
PG_CTL_D = $(PG_CTL) -D "$(PGDATA)"

.PHONY: ci
ci: clean init start stop

.PHONY: init
init:
	$(INITDB) $(PGDATA)
	echo "port = 9998" >> "$(PGDATA)/postgresql.conf"

.PHONY: start
start:
	$(PG_CTL_D) -w -l "$(PG_LOG)" start

.PHONY: start
stop:
	$(PG_CTL_D) -w stop

.PHONY: clean
clean:
	$(PG_CTL_D) stop || true
	rm -rf "$(PGDATA)"
