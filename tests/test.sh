#!/usr/bin/env bash

cat <<'EOF'
  _                          
 | |__   ___  _   _ _ __ ___ 
 | '_ \ / _ \| | | | '__/ __|
 | | | | (_) | |_| | |  \__ \
 |_| |_|\___/ \__,_|_|  |___/

EOF
pass_count=0
fail_count=0

temp_dir=$(mktemp -d)
db_file_path="${temp_dir}/db.db"

echo "hours gen -y -d ${db_file_path}"
hours gen -y -d "${db_file_path}"

tests=(
    "log: today|hours -d ${db_file_path} log -p today|0"
    "log: yest|hours -d ${db_file_path} log -p yest|0"
    "log: 3d|hours -d ${db_file_path} log -p 3d|0"
    "log: week|hours -d ${db_file_path} log -p week|0"
    "log: date|hours -d ${db_file_path} log -p 2024/06/08|0"
    "log: date range|hours -d ${db_file_path} log -p 2024/06/08...2024/06/12|0"
    "log: incorrect argument|hours -d ${db_file_path} log -p blah|1"
    "log: incorrect date|hours -d ${db_file_path} log -p 2024/0608|1"
    "log: incorrect date range|hours -d ${db_file_path} log -p 2024/0608...2024/06/12|1"
    "report: today|hours -d ${db_file_path} report -p today|0"
    "report: yest|hours -d ${db_file_path} report -p yest|0"
    "report: 3d|hours -d ${db_file_path} report -p 3d|0"
    "report: week|hours -d ${db_file_path} report -p week|0"
    "report: date|hours -d ${db_file_path} report -p 2024/06/08|0"
    "report: date range|hours -d ${db_file_path} report -p 2024/06/08...2024/06/12|0"
    "report: incorrect argument|hours -d ${db_file_path} report -p blah|1"
    "report: incorrect date|hours -d ${db_file_path} report -p 2024/0608|1"
    "report: incorrect date range|hours -d ${db_file_path} report -p 2024/0608...2024/06/12|1"
    "report: date range too large|hours -d ${db_file_path} report -p 2024/06/08...2024/06/15|1"
    "stats: today|hours -d ${db_file_path} stats -p today|0"
    "stats: yest|hours -d ${db_file_path} stats -p yest|0"
    "stats: 3d|hours -d ${db_file_path} stats -p 3d|0"
    "stats: week|hours -d ${db_file_path} stats -p week|0"
    "stats: date|hours -d ${db_file_path} stats -p 2024/06/08|0"
    "stats: date range|hours -d ${db_file_path} stats -p 2024/06/08...2024/06/12|0"
    "stats: all|hours -d ${db_file_path} stats -p all|0"
    "stats: incorrect argument|hours -d ${db_file_path} stats -p blah|1"
    "stats: incorrect date|hours -d ${db_file_path} stats -p 2024/0608|1"
    "stats: incorrect date range|hours -d ${db_file_path} stats -p 2024/0608...2024/06/12|1"
)

for test in "${tests[@]}"; do
    IFS='|' read -r title cmd expected_exit_code <<<"$test"

    echo "> $title"
    echo "$cmd"
    echo
    eval "$cmd"
    exit_code=$?
    if [ $exit_code -eq $expected_exit_code ]; then
        echo "✅ command behaves as expected"
        ((pass_count++))
    else
        echo "❌ command returned $exit_code, expected $expected_exit_code"
        ((fail_count++))
    fi
    echo
    echo "==============================="
    echo
done

echo "Summary:"
echo "- Passed: $pass_count"
echo "- Failed: $fail_count"

if [ $fail_count -gt 0 ]; then
    exit 1
else
    exit 0
fi
