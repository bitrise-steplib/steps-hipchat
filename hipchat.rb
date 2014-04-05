require 'hipchat-api'
require 'optparse'

options = {
  token: ENV['HIPCHAT_TOKEN'],
  roomid: ENV['HIPCHAT_ROOMID'],
  fromname: ENV['HIPCHAT_FROMNAME'],
  message: ENV['HIPCHAT_MESSAGE']
}

p "Options: #{options}"

hipchat_api = HipChat::API.new(options[:token])
p hipchat_api.rooms_message(options[:roomid], options[:fromname], options[:message])