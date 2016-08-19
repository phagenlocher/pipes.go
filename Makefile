CC = go
TARGET = pipes
bindir = /usr/local/bin

all: pipes.go
	$(CC) get -x github.com/rthornton128/goncurses
	$(CC) build -x -o $(TARGET) pipes.go

install: all
	mv $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

uninstall:
	rm -f $(DESTDIR)$(bindir)/$(TARGET)

clean:
	rm -f $(TARGET)
