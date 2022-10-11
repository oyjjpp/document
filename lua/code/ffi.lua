#!/usr/bin/luajit

local ffi = require("ffi")

ffi.cdef[[
    int printf(const char *fmt, ...);
    void sleep(int ms);
    int poll(struct pollfd *fds, unsigned long nfds, int timeout);
]]

ffi.C.printf("Hello %s!", "world")

--access standard functions
local sleep

function sleep(s)
    ffi.C.poll(nil, 0 ,s*1000)
end

for i=1, 160 do
    io.write("."); io.flush()
    sleep(0.001)
end

io.write("\n")
