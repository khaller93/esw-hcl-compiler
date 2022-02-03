VERSION=`grep -oP 'const ApplicationVersion string = "\K[-\d.a-zA-Z]+' main.go`

build-docker-image:
	docker build . --build-arg ESW_COMP_VERSION=${version} -t khaller/esw-hcl-compiler:$(VERSION)