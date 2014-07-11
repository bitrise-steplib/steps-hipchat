#!/bin/bash

#
# Run it from the directory which contains step.sh
#


# ------------------------
# --- Helper functions ---

function print_and_do_command {
  echo "$ $@"
  $@
}

function inspect_test_result {
  if [ $1 -eq 0 ]; then
    test_results_success_count=$[test_results_success_count + 1]
  else
    test_results_error_count=$[test_results_error_count + 1]
  fi
}

#
# First param is the expect message, other are the command which will be executed.
#
function expect_success {
  expect_msg=$1
  shift

  echo " -> $expect_msg"
  $@
  cmd_res=$?

  if [ $cmd_res -eq 0 ]; then
    echo " [OK] Expected zero return code, got: 0"
  else
    echo " [ERROR] Expected zero return code, got: $cmd_res"
    exit 1
  fi
}

#
# First param is the expect message, other are the command which will be executed.
#
function expect_error {
  expect_msg=$1
  shift

  echo " -> $expect_msg"
  $@
  cmd_res=$?

  if [ ! $cmd_res -eq 0 ]; then
    echo " [OK] Expected non-zero return code, got: $cmd_res"
  else
    echo " [ERROR] Expected non-zero return code, got: 0"
    exit 1
  fi
}

function is_dir_exist {
  if [ -d "$1" ]; then
    return 0
  else
    return 1
  fi
}

function is_file_exist {
  if [ -f "$1" ]; then
    return 0
  else
    return 1
  fi
}

function is_not_empty {
  if [[ $1 ]]; then
    return 0
  else
    return 1
  fi
}

function test_env_cleanup {
  unset HIPCHAT_TOKEN
  unset HIPCHAT_ROOMID
}

function print_new_test {
  echo
  echo "[TEST]"
}

function run_target_command {
  print_and_do_command ./step.sh
}

# -----------------
# --- Run tests ---

echo "Starting tests..."

test_ipa_path="tests/testfile.ipa"
test_results_success_count=0
test_results_error_count=0


# [TEST] Call the command with HIPCHAT_TOKEN not set, 
# it should raise an error message and exit
# 
(
  print_new_test
  test_env_cleanup

  # Set env vars
  export HIPCHAT_ROOMID="dsa4321"

  # HIPCHAT_TOKEN should NOT exist
  expect_error "HIPCHAT_TOKEN environment variable should NOT be set" is_not_empty "$HIPCHAT_TOKEN"
  expect_success "HIPCHAT_ROOMID environment variable should be set" is_not_empty "$HIPCHAT_ROOMID"

  # Deploy the file
  expect_error "The command should be called, but should not complete sucessfully" run_target_command  
)
test_result=$?
inspect_test_result $test_result


# [TEST] Call the command with HIPCHAT_ROOMID not set, 
# it should raise an error message and exit
# 
(
  print_new_test
  test_env_cleanup

  # Set env vars
  export HIPCHAT_TOKEN="asd1234"

  # HIPCHAT_ROOMID should NOT exist
  expect_error "HIPCHAT_ROOMID environment variable should NOT be set" is_not_empty "$HIPCHAT_ROOMID"
  expect_success "HIPCHAT_TOKEN environment variable should be set" is_not_empty "$HIPCHAT_TOKEN"

  # Deploy the file
  expect_error "The command should be called, but should not complete sucessfully" run_target_command 
)
test_result=$?
inspect_test_result $test_result


#final cleanup
test_env_cleanup

# --------------------
# --- Test Results ---

echo
echo "--- Results ---"
echo " * Errors: $test_results_error_count"
echo " * Success: $test_results_success_count"
echo "---------------"

if [ $test_results_error_count -eq 0 ]; then
  echo "-> SUCCESS"
else
  echo "-> FAILED"
fi
