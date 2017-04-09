PGDATA = $(shell "pwd")/data
INITDB = initdb
PG_CTL = pg_ctl
PG_LOG = $(PGDATA)/log.log
PG_CTL_D = $(PG_CTL) -D "$(PGDATA)"
PGPORT = 9998
PGDB = postgres

.PHONY: ci
ci: clean init start test stop

.PHONY: test
test:
	PGDATA="$(PGDATA)" \
	 PGPORT="$(PGPORT)" \
	 PGDB="$(PGDB)" \
	 go test -v .

.PHONY: init
init:
	$(INITDB) $(PGDATA)
	echo "port = $(PGPORT)" >> "$(PGDATA)/postgresql.conf"

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
