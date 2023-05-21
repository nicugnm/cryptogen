package ro.cryptogen.cryptogenexpose.service

import org.springframework.context.event.EventListener
import org.springframework.stereotype.Service
import ro.cryptogen.cryptogenexpose.events.MessageReceivedEvent

@Service
class MessageProcessingService {

    final lateinit var lastMessage: String
        private set

    @EventListener
    fun onMessageReceived(event: MessageReceivedEvent) {
        lastMessage = event.message
        println("Received message: $lastMessage")
    }
}