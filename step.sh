#!/bin/bash

formatted_output_file_path="$BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH"

function echo_string_to_formatted_output {
  echo "$1" >> $formatted_output_file_path
}

function write_section_to_formatted_output {
  echo '' >> $formatted_output_file_path
  echo "$1" >> $formatted_output_file_path
  echo '' >> $formatted_output_file_path
}

echo "Configs:"
echo " * HIPCHAT_TOKEN: $HIPCHAT_TOKEN"
echo " * HIPCHAT_ROOMID: $HIPCHAT_ROOMID"
echo " * HIPCHAT_FROMNAME: $HIPCHAT_FROMNAME"
echo " * HIPCHAT_MESSAGE_COLOR: $HIPCHAT_MESSAGE_COLOR"
echo " * HIPCHAT_MESSAGE: $HIPCHAT_MESSAGE"
echo

# Input validation
if [ ! -n "$HIPCHAT_TOKEN" ]; then
  echo " [!] HIPCHAT_TOKEN is missing! Terminating..."
  echo
  write_section_to_formatted_output "# Error!"
  write_section_to_formatted_output "Reason: HipChat token (HIPCHAT_TOKEN) is missing!"
  exit 1
fi

if [ ! -n "$HIPCHAT_ROOMID" ]; then
  echo " [!] HIPCHAT_ROOMID is missing! Terminating..."
  echo
  write_section_to_formatted_output "# Error!"
  write_section_to_formatted_output "Reason: HipChat room id (HIPCHAT_ROOMID) is missing!"
  exit 1
fi

from_name='Bitrise'
if [ -n "$HIPCHAT_FROMNAME" ]; then
  from_name="$HIPCHAT_FROMNAME"
fi

msg_color='yellow'
if [ -n "$HIPCHAT_MESSAGE_COLOR" ]; then
  msg_color="$HIPCHAT_MESSAGE_COLOR"
fi

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

curl_response=`curl -d $CONFIG --data-urlencode "message=$HIPCHAT_MESSAGE" "https://api.hipchat.com/v1/rooms/message?auth_token=$HIPCHAT_TOKEN&format=json"`
echo "curl_response: $curl_response"
err_search=$(echo $curl_response | grep error)

if [ "$err_search" == "" ]; then
  write_section_to_formatted_output "# Message successfully sent!"
  write_section_to_formatted_output "### From: ${HIPCHAT_FROMNAME}"
  write_section_to_formatted_output "### Message:"
  write_section_to_formatted_output "${HIPCHAT_MESSAGE}"
  exit 0
else
  echo "Failed"
  write_section_to_formatted_output "# Message send failed!"
  write_section_to_formatted_output "Error message:"
  write_section_to_formatted_output "    ${curl_response}"
fi

exit 1
