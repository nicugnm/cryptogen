package ro.cryptogen.cryptogenexpose.events

import org.springframework.context.ApplicationEvent

class MessageReceivedEvent(source: Any, val message: String) : ApplicationEvent(source)
