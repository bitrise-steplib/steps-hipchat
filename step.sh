#!/bin/bash

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${THIS_SCRIPTDIR}/_bash_utils/utils.sh"
source "${THIS_SCRIPTDIR}/_bash_utils/formatted_output.sh"

# init / cleanup the formatted output
echo "" > "${formatted_output_file_path}"

# Input validation
# - required
if [  -z "$HIPCHAT_TOKEN" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_TOKEN` is not provided!*'
  exit 1
fi

if [ -z "$HIPCHAT_ROOMID" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_ROOMID` is not provided!*'
  exit 1
fi

if [ -z "$HIPCHAT_FROM_NAME" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_FROM_NAME` is not provided!*'
  exit 1
fi

# - optional
if [ -z "$HIPCHAT_MESSAGE" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_MESSAGE` is not provided!*'
fi

if [ -z "$HIPCHAT_MESSAGE_COLOR" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_MESSAGE_COLOR` is not provided!*'
fi

if [ -z "$HIPCHAT_ERROR_FROM_NAME" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_ERROR_FROM_NAME` is not provided!*'
fi

if [ -z "$HIPCHAT_ERROR_MESSAGE_TEXT" ] ; then
  write_section_to_formatted_output '*Notice: `$HIPCHAT_ERROR_MESSAGE_TEXT` is not provided!*'
fi

# Build failed mode
isBuildFailedMode="0"
if [ -n "$STEPLIB_BUILD_STATUS" ] ; then
  isBuildFailedMode="${STEPLIB_BUILD_STATUS}"
fi

# Curl params
error_message="${HIPCHAT_MESSAGE}"
if [ -n "$HIPCHAT_ERROR_MESSAGE_TEXT" ] ; then
  error_message="${HIPCHAT_ERROR_FROM_NAME}"
fi

error_from_name="${HIPCHAT_FROM_NAME}"
if [ -n "$HIPCHAT_ERROR_FROM_NAME" ] ; then
  error_from_name="${HIPCHAT_ERROR_FROM_NAME}"
fi

message="${HIPCHAT_MESSAGE}"
if [[ "${isBuildFailedMode}" == "1" ]] ; then
  message="${error_message}"
fi

from_name="${HIPCHAT_FROM_NAME}"
if [[ "${isBuildFailedMode}" == "1" ]] ; then
  from_name="${error_from_name}"
fi

msg_color='yellow'
if [ -n "$HIPCHAT_MESSAGE_COLOR" ] ; then
  msg_color="${HIPCHAT_MESSAGE_COLOR}"
fi

echo "Configs:"
echo " * BUILD_FAILED_MODE: $isBuildFailedMode"
echo " * HIPCHAT_TOKEN: $HIPCHAT_TOKEN"
echo " * HIPCHAT_ROOMID: $HIPCHAT_ROOMID"
echo " * HIPCHAT_FROM_NAME: $from_name"
echo " * HIPCHAT_MESSAGE_COLOR: $msg_color"
echo " * HIPCHAT_MESSAGE: $message"
echo


urlencode() {
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

from_name=$(urlencode "$from_name")
msg_color=$(urlencode "$msg_color")

CONFIG="room_id=$HIPCHAT_ROOMID&from=$from_name&color=$msg_color"

curl_response=`curl -d $CONFIG --data-urlencode "message=$message" "https://api.hipchat.com/v1/rooms/message?auth_token=$HIPCHAT_TOKEN&format=json"`
echo "curl_response: $curl_response"
err_search=$(echo $curl_response | grep error)

if [[ "${isBuildFailedMode}" == "1" ]] ; then
  write_section_to_formatted_output "# Message send failed!"
  write_section_to_formatted_output "Error message:"
  write_section_to_formatted_output "    ${message}"
  exit 1
elif [ "$err_search" == "" ] ; then
  write_section_to_formatted_output "# Message successfully sent!"
  write_section_to_formatted_output "## From:"
  write_section_to_formatted_output "${HIPCHAT_FROM_NAME}"
  write_section_to_formatted_output "## To Room:"
  write_section_to_formatted_output "${HIPCHAT_ROOMID}"
  write_section_to_formatted_output "## Message:"
  write_section_to_formatted_output "${message}"
  exit 0
else
  echo "Failed"
  write_section_to_formatted_output "# Message send failed!"
  write_section_to_formatted_output "Error message:"
  write_section_to_formatted_output "    ${curl_response}"
fi

exit 1
