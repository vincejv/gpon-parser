#!/bin/sh

otop="/opt/lantiq/bin/otop"
otop_output="$($otop -b -g s 2>/dev/null)"

# === CPU Usage (1s sample) ===
cpu_usage() {
    read cpu user nice system idle rest < /proc/stat
    total1=$((user + nice + system + idle))
    idle1=$idle
    sleep 1
    read cpu user nice system idle rest < /proc/stat
    total2=$((user + nice + system + idle))
    idle2=$idle

    total_diff=$((total2 - total1))
    idle_diff=$((idle2 - idle1))
    usage=$((100 * (total_diff - idle_diff) / total_diff))
    echo "$usage"
}

# === Used RAM (%) ===
used_ram_percent() {
    total=$(awk '/MemTotal:/ {print $2}' /proc/meminfo)
    available=$(awk '/MemAvailable:/ {print $2}' /proc/meminfo)
    [ -z "$available" ] && available=$(awk '/MemFree:/ {print $2}' /proc/meminfo)
    used=$((total - available))
    echo $((100 * used / total))
}

# === Optical Stats from cached otop_output ===
extract_optical_info() {
    ddmi_voltage=$(echo "$otop_output" | awk '/^DDMI voltage/ { print $3 }')

    tx_power_dbm=$(echo "$otop_output" | grep -i '^tx power' | sed -n 's/.*\([-\+]\?[0-9]\+\.[0-9]\+dBm\).*/\1/p')
    rssi_power_dbm=$(echo "$otop_output" | grep -i '^RSSI 1490 power' | awk '{print $NF}')

    echo "DDMI Voltage          : $ddmi_voltage"
    echo "TX Power (dBm)        : $tx_power_dbm"
    echo "RSSI 1490 Power (dBm) : $rssi_power_dbm"
}

# === GPON Serial Number ===
gpon_serial() {
    serial=$(fw_printenv nSerial 2>/dev/null | cut -f2 -d=)
    echo "GPON Serial      : $serial"
}

# === Temperature (Die & Laser) from cached otop_output ===
temperature() {
    cpu_temp=$(echo "$otop_output" | grep 'temperature' | grep 'die' | cut -c 52-54)
    laser_temp=$(echo "$otop_output" | grep 'temperature' | grep 'laser' | cut -c 52-54)
    [ -n "$cpu_temp" ] && cpu_temp=$(expr $cpu_temp - 273)
    [ -n "$laser_temp" ] && laser_temp=$(expr $laser_temp - 273)
    echo "Temp (Die/Laser) : ${cpu_temp}℃ / ${laser_temp}℃"
}

# === System Uptime ===
uptime_str() {
    uptime_seconds=$(cut -d. -f1 /proc/uptime)
    echo "Uptime (secs)    : $uptime_seconds"
}

model() {
    vendorname=$(cat /tmp/vendorname 2>/dev/null)
    if [ "$vendorname" = "HUAWEI" ]; then
        modelname="SmartAX MA5671A"
    elif [ "$vendorname" = "Nokia" ]; then
        modelname="G-010S-A"
    else
        modelname="G-010S-P"
    fi
    echo "$modelname"
}

omcid_version() {
    omcid="/opt/lantiq/bin/omcid"
    if [ -x "$omcid" ]; then
        ver=$($omcid -v 2>/dev/null | tail -n 1 | cut -c 18-75)
        ver_o=$(echo "$ver" | grep -c '6BA1896SPE2C05')
        if [ "$ver_o" = "1" ]; then
            ver="6BA1896SPE2C05"
        fi
        echo "$ver"
    else
        echo "omcid not found"
    fi
}

bias_current() {
    if [ -n "$otop_output" ]; then
        echo "$otop_output" | grep -i '^actual bias / modulation current' | awk '{print $6}'
    else
        echo "otop output missing"
    fi
}

# === Output ===
echo "CPU Usage        : $(cpu_usage)%"
echo "Used RAM         : $(used_ram_percent)%"
echo "Model              : $(model)"
echo "omcid Version      : $(omcid_version)"
extract_optical_info
echo "Bias Current     : $(bias_current)"
gpon_serial
temperature
uptime_str

