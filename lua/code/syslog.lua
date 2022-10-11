--[[
--@desc:syslog library
--@author oyj<ouyangjun@zhangyue.com>
--@date 2019-02-16
--]]

local ffi = require("ffi")
local json = require("json")

--syslog function
ffi.cdef[[
    void openlog(const char *ident, int option, int facility);
    void syslog(int priority, const char *format, ...);
    void closelog(void);
]]

--log type
local log_facility = {
    ["kern"] = 0,
    ["user"] = 8,
    ["mail"] = 16,
    ["daemon"] = 24,
    ["auth"] = 32,
    ["syslog"] = 40,
    ["lpr"] = 48,
    ["news"] = 56,
    ["uucp"] = 64,
    ["cron"] = 72,
    ["authoriv"] = 80,
    ["local0"] = 128,
    ["local1"] = 136,
    ["local2"] = 144,
    ["local3"] = 152,
    ["local4"] = 160,
    ["local5"] = 168,
    ["local6"] = 176,
    ["local7"] = 184,
}

--log level
local log_priority = {
    ["emerg"] = 0,
    ["alert"] = 1,
    ["crit"]  = 2,
    ["err"]   = 3,
    ["warning"] = 4,
    ["notice"] = 5,
    ["info"] = 6,
    ["debug"] = 7,
}

--rsyslog openlog function
local function openlog(ident, option, facility)
    ffi.C.openlog(ident, option, facility)
end

--rsyslog syslog function
local function syslog(priority,log)
    ffi.C.syslog(priority,log)
end

--rsyslog closelog function
local function closelog()
    ffi.C.closelog()
end

local function get_log_id()
    math.randomseed(os.time())
    num = math.random()
    return string.sub(num,-10)
end
--rsyslog write_log function
local function write_log(header, body, level, log_type)
    if header["rate"] == nil then
        header["rate"] = 0
    end

    if header["interval"] == nil then
        header["interval"] = 0
    end

    header["time"] = os.date("%Y-%m-%d %H:%M:%S")
    header["logId"] = get_log_id() 

    data = {
        ["header"] = header,
        ["body"] = body
    }

    content = json.encode(data)
    openlog("PHP_BIZ_ERRLOG", 1, log_type)
    syslog(level, " JSON"..content)
    closelog()
end

--get_log_id function
local function get_log_id_ver()
    math.randomseed(os.time())
    num = math.random()
    return string.sub(num,-10)
end

return {
    openlog = openlog,
    syslog = syslog,
    closelog = closelog,
    facility = log_facility,
    priority = log_priority,
    write_log = write_log
}
