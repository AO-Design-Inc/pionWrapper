

pionHandler.so pionHandler.h: pion_handler.go
	go build -o $@ -buildmode=c-shared $<

.PHONY : clean
clean:
	rm *.so && rm *.h

