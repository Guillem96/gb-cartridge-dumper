BUILDDIR=build
TARGET=gb-dumper

RPI_USER=pi
RPI_HOST=192.168.1.78
RPI_DESTINATION=/home/pi/$(TARGET)


directories:
	mkdir -p $(BUILDDIR)

build: directories
	go build -o $(BUILDDIR)/$(TARGET) main.go dumper.go

rpi-build:
	GOOS=linux GOARCH=arm GOARM=5 go build -o $(BUILDDIR)/rpi-$(TARGET) main.go dumper.go
	chmod +x $(BUILDDIR)/rpi-$(TARGET)

rpi-upload: rpi-build
	sshpass -p raspberry scp $(BUILDDIR)/rpi-$(TARGET) $(RPI_USER)@$(RPI_HOST):$(RPI_DESTINATION)

rpi-run: rpi-upload
	sshpass -p raspberry ssh $(RPI_USER)@$(RPI_HOST) "$(RPI_DESTINATION)"

rpi-get-rom:
	sshpass -p raspberry scp $(RPI_USER)@$(RPI_HOST):/home/pi/rom.gb rom.gb

run: build
	$(BUILDDIR)/$(TARGET)