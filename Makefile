all: glide_install clean test prometheus_ps

binary := "prometheus_ps"

prometheus_ps:
	go build $(binary)

glide_install:
	glide install

run: clean prometheus_ps
	./$(binary) --config conf.json

test:
	go test prometheus_ps

clean:
	go clean
	if [ -a prometheus_ps ] ; \
  	then \
    	rm prometheus_ps ; \
	fi;
