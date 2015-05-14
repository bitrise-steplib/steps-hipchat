#!/bin/bash

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${THIS_SCRIPTDIR}/_bash_utils/utils.sh"
source "${THIS_SCRIPTDIR}/_bash_utils/formatted_output.sh"

# init / cleanup the formatted output
echo "" > "${formatted_output_file_path}"

function write_error_section_to_foramtted_output {
  error_message="$1"

  echo "Failed"
  write_section_to_formatted_output "# Message send failed!"
  write_section_to_formatted_output "Error message:"
  write_section_to_formatted_output "    ${error_message}"
}

function write_success_section_to_foramtted_output {
  success_from_name="$1"
  success_room_id="$2"
  success_message="$3"

  echo "Success"
  write_section_to_formatted_output "# Message successfully sent!"
  write_section_to_formatted_output "## From:"
  write_section_to_formatted_output "     ${success_from_name}"
  write_section_to_formatted_output "## To Room:"
  write_section_to_formatted_output "     ${success_room_id}"
  write_section_to_formatted_output "## Message:"
  write_section_to_formatted_output "     ${success_message}"
}

# Input validation
# - required
if [  -z "$HIPCHAT_TOKEN" ] ; then
  write_error_section_to_foramtted_output '`$HIPCHAT_TOKEN` is not provided!'
  exit 1
fi

if [ -z "$HIPCHAT_ROOMID" ] ; then
  write_error_section_to_foramtted_output '`$HIPCHAT_ROOMID` is not provided!'
  exit 1
fi

if [ -z "$HIPCHAT_FROMNAME" ] ; then
  write_error_section_to_foramtted_output '`$HIPCHAT_FROMNAME` is not provided!'
  exit 1
fi

if [ -z "$HIPCHAT_MESSAGE" ] ; then
  write_error_section_to_foramtted_output '`$HIPCHAT_MESSAGE` is not provided!'
  exit 1
fi

# - optional
if [ -z "$HIPCHAT_MESSAGE_COLOR" ] ; then
  write_section_to_formatted_output '`$HIPCHAT_MESSAGE_COLOR` is not provided!'
fi

if [ -z "$HIPCHAT_ERROR_FROMNAME" ] ; then
  write_section_to_formatted_output '`$HIPCHAT_ERROR_FROMNAME` is not provided!'
fi

if [ -z "$HIPCHAT_ERROR_MESSAGE" ] ; then
  write_section_to_formatted_output '`$HIPCHAT_ERROR_MESSAGE` is not provided!'
fi

# Build failed mode
isBuildFailedMode="0"
if [ -n "$STEPLIB_BUILD_STATUS" ] ; then
  isBuildFailedMode="${STEPLIB_BUILD_STATUS}"
fi

# Curl params
message="${HIPCHAT_MESSAGE}"
if [[ "${isBuildFailedMode}" == "1" ]] ; then
  if [ -n "$HIPCHAT_ERROR_MESSAGE" ] ; then
    message="${HIPCHAT_ERROR_MESSAGE}"
  fi
fi

from_name="${HIPCHAT_FROMNAME}"
if [[ "${isBuildFailedMode}" == "1" ]] ; then
  if [ -n "$HIPCHAT_ERROR_FROMNAME" ] ; then
    from_name="${HIPCHAT_ERROR_FROMNAME}"
  fi
fi

msg_color='yellow'
if [ -n "$HIPCHAT_MESSAGE_COLOR" ] ; then
  msg_color="${HIPCHAT_MESSAGE_COLOR}"
fi

echo "Configs:"
echo " * BUILD_FAILED_MODE: $isBuildFailedMode"
echo " * HIPCHAT_TOKEN: $HIPCHAT_TOKEN"
echo " * HIPCHAT_ROOMID: $HIPCHAT_ROOMID"
echo " * HIPCHAT_FROMNAME: $from_name"
echo " * HIPCHAT_MESSAGE_COLOR: $msg_color"
echo " * HIPCHAT_MESSAGE: $message"
echo

function urlencode {
  # urlencode <string>
  #  source: https://gist.github.com/cdown/1163649

  local length="${#1}"
  for (( i = 0; i < length; i++ )); do
    local c="${1:i:1}"
    case $c in
      [a-zA-Z0-9.~_-]) printf "$c" ;;
      *) printf '%%%02X' "'$c"
    esac
  done
}

url_encoded_from_name=$(urlencode "$from_name")
msg_color=$(urlencode "$msg_color")

CONFIG="room_id=$HIPCHAT_ROOMID&from=$url_encoded_from_name&color=$msg_color"

curl_response=`curl -d $CONFIG --data-urlencode "message=$message" "https://api.hipchat.com/v1/rooms/message?auth_token=$HIPCHAT_TOKEN&format=json"`
echo "curl_response: $curl_response"
err_search=$(echo $curl_response | grep error)

if [[ "$err_search" != "" ]]; then
  write_error_section_to_foramtted_output "${curl_response}"
  exit 1
fi

write_success_section_to_foramtted_output "${from_name}" "${HIPCHAT_ROOMID}" "${message}"
exit 0
