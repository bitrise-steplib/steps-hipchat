Send HipChat message Step
=======

HipChat message sender Step.

This Step is part of the [Open StepLib](http://www.steplib.com/), you can find its StepLib page [here](http://www.steplib.com/step/hipchat).

# Input Environment Variables
- **HIPCHAT_TOKEN**

    at [https://bitfall-team.hipchat.com/admin/api](https://bitfall-team.hipchat.com/admin/api)
- **HIPCHAT_ROOMID**

	at [https://bitfall-team.hipchat.com/admin/rooms](https://bitfall-team.hipchat.com/admin/rooms)
- **HIPCHAT_FROMNAME**

	the name of the sender, maximum 15 characters; optional, "Bitrise" by default
- **HIPCHAT_MESSAGE_COLOR**

	the color of the message; optional, "yellow" by default
- **HIPCHAT_MESSAGE**

	the text message you would like to send

# How to test/run locally?

- clone this repository
- cd into the repository folder
- run: HIPCHAT_TOKEN=[your-hipchat-token] HIPCHAT_ROOMID=[your-hipchat-roomid] HIPCHAT_FROMNAME=[max-15-characters] HIPCHAT_MESSAGE_COLOR=[color] HIPCHAT_MESSAGE=[the-hipchat-message] bash step.sh
