require 'hipchat-api'

options = {
  token: ENV['HIPCHAT_TOKEN'],
  roomid: ENV['HIPCHAT_ROOMID'],
  fromname: ENV['HIPCHAT_FROMNAME'],
  message: ENV['HIPCHAT_MESSAGE']
}

p "Options: #{options}"

begin
  hipchat_api = HipChat::API.new(options[:token])
  resp = hipchat_api.rooms_message(options[:roomid], options[:fromname], options[:message])
  p resp
  if resp["error"]
    puts %{ [i] Error: #{resp["error"]}}
    exit 1
  end
rescue => ex
  puts "Exception happened: #{ex}"
  exit 1
end