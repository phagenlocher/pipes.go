CC = go
TARGET = pipes
bindir = /usr/local/bin

$(TARGET): *.go
	$(CC) get -x github.com/rthornton128/goncurses
	$(CC) build -x -o $(TARGET) $<

.PHONY all: $(TARGET)

.PHONY install: all
	mv $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

.PHONY uninstall:
	rm -f $(DESTDIR)$(bindir)/$(TARGET)

.PHONY clean:
	rm -f $(TARGET)
